#ifndef CMRC_CMRC_HPP_INCLUDED
#define CMRC_CMRC_HPP_INCLUDED

#include <string>
#include <map>
#include <mutex>

#define CMRC_INIT(libname) \
    do { \
	extern void cmrc_init_resources_##libname(); \
	cmrc_init_resources_##libname(); \
    } while (0)

namespace cmrc {

class resource {
    const char* _begin = nullptr;
    const char* _end = nullptr;
public:
    const char* begin() const { return _begin; }
    const char* end() const { return _end; }

    resource() = default;
    resource(const char* beg, const char* end) : _begin(beg), _end(end) {}
};

using resource_table = std::map<std::string, resource>;

namespace detail {

inline resource_table& table_instance() {
    static resource_table table;
    return table;
}

inline std::mutex& table_instance_mutex() {
    static std::mutex mut;
    return mut;
}

// We restrict access to the resource table through a mutex so that multiple
// threads can access it safely.
template <typename Func>
inline auto with_table(Func fn) -> decltype(fn(std::declval<resource_table&>())) {
    std::lock_guard<std::mutex> lk{ table_instance_mutex() };
    return fn(table_instance());
}

}

inline resource open(const char* fname) {
    return detail::with_table([fname](const resource_table& table) {
	auto iter = table.find(fname);
	if (iter == table.end()) {
	    return resource {};
	}
	return iter->second;
    });
}

inline resource open(const std::string& fname) {
    return open(fname.data());
}

}

#endif // CMRC_CMRC_HPP_INCLUDED
