package cproto

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sapariduo/reindexer/bindings"
	"github.com/sapariduo/reindexer/cjson"
)

type bufPtr struct {
	rseq uint32
	buf  *NetBuffer
}

type sig chan bufPtr

const bufsCap = 16 * 1024
const queueSize = 512
const maxSeqNum = queueSize * 1000000

const cprotoMagic = 0xEEDD1132
const cprotoVersion = 0x102
const cprotoMinCompatVersion = 0x101
const cprotoHdrLen = 16
const deadlineCheckPeriodSec = 1

const (
	cmdPing              = 0
	cmdLogin             = 1
	cmdOpenDatabase      = 2
	cmdCloseDatabase     = 3
	cmdDropDatabase      = 4
	cmdOpenNamespace     = 16
	cmdCloseNamespace    = 17
	cmdDropNamespace     = 18
	cmdTruncateNamespace = 19
	cmdAddIndex          = 21
	cmdEnumNamespaces    = 22
	cmdDropIndex         = 24
	cmdUpdateIndex       = 25
	cmdAddTxItem         = 26
	cmdCommitTx          = 27
	cmdRollbackTx        = 28
	cmdStartTransaction  = 29
	cmdDeleteQueryTx     = 30
	cmdUpdateQueryTx     = 31
	cmdCommit            = 32
	cmdModifyItem        = 33
	cmdDeleteQuery       = 34
	cmdUpdateQuery       = 35
	cmdSelect            = 48
	cmdSelectSQL         = 49
	cmdFetchResults      = 50
	cmdCloseResults      = 51
	cmdGetMeta           = 64
	cmdPutMeta           = 65
	cmdEnumMeta          = 66
	cmdCodeMax           = 128
)

type requestInfo struct {
	seqNum   uint32
	repl     sig
	deadline uint32
	isAsync  int32
	cmpl     bindings.RawCompletion
	cmplLock sync.Mutex
}

type connection struct {
	owner *NetCProto
	conn  net.Conn

	wrBuf, wrBuf2 *bytes.Buffer
	wrKick        chan struct{}

	rdBuf *bufio.Reader

	seqs chan uint32
	lock sync.RWMutex

	err   error
	errCh chan struct{}

	lastReadStamp int64

	now    uint32
	termCh chan struct{}

	requests [queueSize]requestInfo
}

func newConnection(ctx context.Context, owner *NetCProto) (c *connection, err error) {
	c = &connection{
		owner:  owner,
		wrBuf:  bytes.NewBuffer(make([]byte, 0, bufsCap)),
		wrBuf2: bytes.NewBuffer(make([]byte, 0, bufsCap)),
		wrKick: make(chan struct{}, 1),
		seqs:   make(chan uint32, queueSize),
		errCh:  make(chan struct{}),
		termCh: make(chan struct{}),
	}
	for i := 0; i < queueSize; i++ {
		c.seqs <- uint32(i)
		c.requests[i].repl = make(sig)
	}

	go c.deadlineTicker()

	intCtx, cancel := applyTimeout(ctx, uint32(owner.timeouts.LoginTimeout/time.Second))
	if cancel != nil {
		defer cancel()
	}

	if err = c.connect(intCtx); err != nil {
		c.onError(err)
		return
	}

	if err = c.login(intCtx, owner); err != nil {
		c.onError(err)
		return
	}
	return
}

func seqNumIsValid(seqNum uint32) bool {
	if seqNum < maxSeqNum {
		return true
	}
	return false
}

func (c *connection) deadlineTicker() {
	timeout := time.Second * time.Duration(deadlineCheckPeriodSec)
	ticker := time.NewTicker(timeout)
	atomic.StoreUint32(&c.now, 0)
	for range ticker.C {
		select {
		case <-c.errCh:
			return
		case <-c.termCh:
			return
		default:
		}
		now := atomic.AddUint32(&c.now, deadlineCheckPeriodSec)
		for i := range c.requests {
			seqNum := atomic.LoadUint32(&c.requests[i].seqNum)
			if !seqNumIsValid(seqNum) {
				continue
			}
			deadline := atomic.LoadUint32(&c.requests[i].deadline)
			if deadline != 0 && now >= deadline && atomic.CompareAndSwapUint32(&c.requests[i].deadline, deadline, 0) {
				if atomic.LoadInt32(&c.requests[i].isAsync) != 0 {
					c.requests[i].cmplLock.Lock()
					if c.requests[i].cmpl != nil && deadline == atomic.LoadUint32(&c.requests[i].deadline) {
						cmpl := c.requests[i].cmpl
						c.requests[i].cmpl = nil
						seqNum = atomic.LoadUint32(&c.requests[i].seqNum)
						atomic.StoreUint32(&c.requests[i].seqNum, maxSeqNum)
						atomic.StoreInt32(&c.requests[i].isAsync, 0)
						c.requests[i].cmplLock.Unlock()
						c.seqs <- nextSeqNum(seqNum)
						cmpl(nil, context.DeadlineExceeded)
					} else {
						c.requests[i].cmplLock.Unlock()
					}
				}
			}
		}
	}
}

func (c *connection) connect(ctx context.Context) (err error) {
	var d net.Dialer
	c.conn, err = d.DialContext(ctx, "tcp", c.owner.url.Host)
	if err != nil {
		return err
	}
	c.conn.(*net.TCPConn).SetNoDelay(true)
	c.rdBuf = bufio.NewReaderSize(c.conn, bufsCap)

	go c.writeLoop()
	go c.readLoop()
	return
}

func (c *connection) login(ctx context.Context, owner *NetCProto) (err error) {
	password, username, path := "", "", owner.url.Path
	if owner.url.User != nil {
		username = owner.url.User.Username()
		password, _ = owner.url.User.Password()
	}
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	buf, err := c.rpcCall(ctx, cmdLogin, 0, username, password, path, c.owner.connectOpts.CreateDBIfMissing)
	if err != nil {
		c.err = err
		return
	}
	defer buf.Free()

	if len(buf.args) > 1 {
		owner.checkServerStartTime(buf.args[1].(int64))
	}
	return
}

func (c *connection) readLoop() {
	var err error
	var hdr = make([]byte, cprotoHdrLen)
	for {
		if err = c.readReply(hdr); err != nil {
			c.onError(err)
			return
		}
		atomic.StoreInt64(&c.lastReadStamp, time.Now().Unix())
	}
}

func (c *connection) readReply(hdr []byte) (err error) {
	if _, err = io.ReadFull(c.rdBuf, hdr); err != nil {
		return
	}
	ser := cjson.NewSerializer(hdr)
	magic := ser.GetUInt32()

	if magic != cprotoMagic {
		return fmt.Errorf("Invalid cproto magic '%08X'", magic)
	}

	version := ser.GetUInt16()
	_ = int(ser.GetUInt16())
	size := int(ser.GetUInt32())
	rseq := uint32(ser.GetUInt32())

	if version < cprotoMinCompatVersion {
		return fmt.Errorf("Unsupported cproto version '%04X'. This client expects reindexer server v1.9.8+", version)
	}

	if !seqNumIsValid(rseq) {
		return fmt.Errorf("invalid seq num: %d", rseq)
	}
	reqID := rseq % queueSize
	if atomic.LoadUint32(&c.requests[reqID].seqNum) != rseq {
		io.CopyN(ioutil.Discard, c.rdBuf, int64(size))
		return
	}
	repCh := c.requests[reqID].repl
	answ := newNetBuffer()
	answ.reset(size, c)

	if _, err = io.ReadFull(c.rdBuf, answ.buf); err != nil {
		return
	}

	if atomic.LoadInt32(&c.requests[reqID].isAsync) != 0 {
		c.requests[reqID].cmplLock.Lock()
		if c.requests[reqID].cmpl != nil && atomic.LoadUint32(&c.requests[reqID].seqNum) == rseq {
			cmpl := c.requests[reqID].cmpl
			c.requests[reqID].cmpl = nil
			atomic.StoreUint32(&c.requests[reqID].seqNum, maxSeqNum)
			atomic.StoreInt32(&c.requests[reqID].isAsync, 0)
			c.requests[reqID].cmplLock.Unlock()
			c.seqs <- nextSeqNum(rseq)
			cmpl(answ, answ.parseArgs())
		} else {
			c.requests[reqID].cmplLock.Unlock()
		}
	} else if repCh != nil {
		repCh <- bufPtr{rseq, answ}
	} else {
		return fmt.Errorf("unexpected answer: %v", answ)
	}
	return
}

func (c *connection) write(buf []byte) {
	c.lock.Lock()
	c.wrBuf.Write(buf)
	c.lock.Unlock()
	select {
	case c.wrKick <- struct{}{}:
	default:
	}
}

func (c *connection) writeLoop() {
	for {
		select {
		case <-c.errCh:
			return
		case <-c.wrKick:
		}
		c.lock.Lock()
		if c.wrBuf.Len() == 0 {
			err := c.err
			c.lock.Unlock()
			if err == nil {
				continue
			} else {
				return
			}
		}
		c.wrBuf, c.wrBuf2 = c.wrBuf2, c.wrBuf
		c.lock.Unlock()

		if _, err := c.wrBuf2.WriteTo(c.conn); err != nil {
			c.onError(err)
			return
		}
	}
}

func nextSeqNum(seqNum uint32) uint32 {
	seqNum += queueSize
	if seqNum < maxSeqNum {
		return seqNum
	}
	return seqNum - maxSeqNum
}

func (c *connection) packRPC(cmd int, seq uint32, execTimeout int, args ...interface{}) {
	in := newRPCEncoder(cmd, seq)
	for _, a := range args {
		switch t := a.(type) {
		case bool:
			in.boolArg(t)
		case int:
			in.intArg(t)
		case int32:
			in.intArg(int(t))
		case int64:
			in.int64Arg(t)
		case string:
			in.stringArg(t)
		case []byte:
			in.bytesArg(t)
		case []int32:
			in.int32ArrArg(t)
		}
	}

	in.startArgsChunck()
	in.int64Arg(int64(execTimeout))

	c.write(in.ser.Bytes())
	in.ser.Close()
}

func (c *connection) awaitSeqNum(ctx context.Context) (seq uint32, remainingTimeout int, err error) {
	select {
	case seq = <-c.seqs:
		if err = ctx.Err(); err != nil {
			c.seqs <- seq
			return
		}
		if execDeadline, ok := ctx.Deadline(); ok {
			remainingTimeout = int(execDeadline.Sub(time.Now()) / time.Millisecond)
			if remainingTimeout <= 0 {
				c.seqs <- seq
				err = context.DeadlineExceeded
			}
		}
	case <-ctx.Done():
		err = ctx.Err()
	}
	return
}

func applyTimeout(ctx context.Context, timeout uint32) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, time.Second*time.Duration(timeout))
}

func (c *connection) rpcCallAsync(ctx context.Context, cmd int, netTimeout uint32, cmpl bindings.RawCompletion, args ...interface{}) {
	if err := c.curError(); err != nil {
		cmpl(nil, err)
		return
	}

	intCtx, cancel := applyTimeout(ctx, netTimeout)
	if cancel != nil {
		defer cancel()
	}

	seq, timeout, err := c.awaitSeqNum(intCtx)
	if err != nil {
		cmpl(nil, err)
		return
	}
	reqID := seq % queueSize
	c.requests[reqID].cmplLock.Lock()
	c.requests[reqID].cmpl = cmpl
	atomic.StoreInt32(&c.requests[reqID].isAsync, 1)
	atomic.StoreUint32(&c.requests[reqID].seqNum, seq)
	c.requests[reqID].cmplLock.Unlock()

	if timeout != 0 {
		atomic.StoreUint32(&c.requests[reqID].deadline, atomic.LoadUint32(&c.now)+uint32(timeout))
	}

	c.packRPC(cmd, seq, timeout, args...)

	return
}

func (c *connection) rpcCall(ctx context.Context, cmd int, netTimeout uint32, args ...interface{}) (buf *NetBuffer, err error) {
	intCtx, cancel := applyTimeout(ctx, netTimeout)
	if cancel != nil {
		defer cancel()
	}
	seq, timeout, err := c.awaitSeqNum(intCtx)
	if err != nil {
		return nil, err
	}

	reqID := seq % queueSize
	reply := c.requests[reqID].repl

	atomic.StoreUint32(&c.requests[reqID].seqNum, seq)
	c.packRPC(cmd, seq, timeout, args...)

for_loop:
	for {
		select {
		case bufPtr := <-reply:
			if bufPtr.rseq == seq {
				buf = bufPtr.buf
				break for_loop
			} else {
				bufPtr.buf.Free()
			}
		case <-c.errCh:
			c.lock.RLock()
			err = c.err
			c.lock.RUnlock()
			break for_loop
		case <-intCtx.Done():
			err = intCtx.Err()
			break for_loop
		}
	}

	atomic.StoreUint32(&c.requests[reqID].seqNum, maxSeqNum)

	select {
	case bufPtr := <-reply:
		bufPtr.buf.Free()
	default:
	}

	c.seqs <- nextSeqNum(seq)
	if err != nil {
		return
	}
	if err = buf.parseArgs(); err != nil {
		return
	}
	return
}

func (c *connection) onError(err error) {
	c.lock.Lock()
	if c.err == nil {
		c.err = err
		if c.conn != nil {
			c.conn.Close()
		}
		select {
		case <-c.errCh:
		default:
			close(c.errCh)
		}

		for i := range c.requests {
			if atomic.LoadInt32(&c.requests[i].isAsync) != 0 {
				c.requests[i].cmplLock.Lock()
				if c.requests[i].cmpl != nil {
					cmpl := c.requests[i].cmpl
					c.requests[i].cmpl = nil
					seqNum := atomic.LoadUint32(&c.requests[i].seqNum)
					atomic.StoreUint32(&c.requests[i].seqNum, maxSeqNum)
					atomic.StoreInt32(&c.requests[i].isAsync, 0)
					c.requests[i].cmplLock.Unlock()
					c.seqs <- nextSeqNum(seqNum)
					cmpl(nil, err)
				} else {
					c.requests[i].cmplLock.Unlock()
				}
			}
		}
	}
	c.lock.Unlock()
}

func (c *connection) hasError() (has bool) {
	c.lock.RLock()
	has = c.err != nil
	c.lock.RUnlock()
	return
}

func (c *connection) curError() error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.err
}

func (c *connection) lastReadTime() time.Time {
	stamp := atomic.LoadInt64(&c.lastReadStamp)
	return time.Unix(stamp, 0)
}

func (c *connection) Finalize() error {
	close(c.termCh)
	return nil
}
