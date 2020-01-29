# Install script for directory: /home/vagrant/MDM/reindexer/cpp_src

# Set the install prefix
if(NOT DEFINED CMAKE_INSTALL_PREFIX)
  set(CMAKE_INSTALL_PREFIX "/usr/local")
endif()
string(REGEX REPLACE "/$" "" CMAKE_INSTALL_PREFIX "${CMAKE_INSTALL_PREFIX}")

# Set the install configuration name.
if(NOT DEFINED CMAKE_INSTALL_CONFIG_NAME)
  if(BUILD_TYPE)
    string(REGEX REPLACE "^[^A-Za-z0-9_]+" ""
           CMAKE_INSTALL_CONFIG_NAME "${BUILD_TYPE}")
  else()
    set(CMAKE_INSTALL_CONFIG_NAME "RelWithDebInfo")
  endif()
  message(STATUS "Install configuration: \"${CMAKE_INSTALL_CONFIG_NAME}\"")
endif()

# Set the component getting installed.
if(NOT CMAKE_INSTALL_COMPONENT)
  if(COMPONENT)
    message(STATUS "Install component: \"${COMPONENT}\"")
    set(CMAKE_INSTALL_COMPONENT "${COMPONENT}")
  else()
    set(CMAKE_INSTALL_COMPONENT)
  endif()
endif()

# Install shared libraries without execute permission?
if(NOT DEFINED CMAKE_INSTALL_SO_NO_EXE)
  set(CMAKE_INSTALL_SO_NO_EXE "1")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/tools" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/tools/errors.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/tools" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/tools/serializer.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/tools" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/tools/varint.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/tools" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/tools/stringstools.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/tools" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/tools/customhash.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/reindexer.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/type_consts.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/item.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/payload" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/payload/payloadvalue.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/payload" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/payload/payloadiface.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/indexopts.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/namespacedef.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/keyvalue" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/keyvalue/variant.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/sortingprioritiestable.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/rdxcontext.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/activity_context.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/cbinding" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/cbinding/reindexer_c.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/cbinding" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/cbinding/reindexer_ctypes.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/transaction.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/query" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/query/query.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/query" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/query/queryentry.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/queryresults" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/queryresults/queryresults.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/indexdef.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/queryresults" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/queryresults/aggregationresult.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core/queryresults" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/queryresults/itemref.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/core" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/core/expressiontree.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/h_vector.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/string_view.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/mutex.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/intrusive_ptr.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/trivial_reverse_iterator.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/estl" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/estl/span.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/reindexer.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/item.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/reindexerconfig.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/queryresults.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/resultserializer.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/client" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/client/internalrdxcontext.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/debug" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/debug/backtrace.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/include/reindexer/debug" TYPE FILE FILES "/home/vagrant/MDM/reindexer/cpp_src/debug/allocdebug.h")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/lib" TYPE STATIC_LIBRARY FILES "/home/vagrant/MDM/reindexer/build/cpp_src/libreindexer.a")
endif()

if(NOT CMAKE_INSTALL_COMPONENT OR "${CMAKE_INSTALL_COMPONENT}" STREQUAL "dev")
  file(INSTALL DESTINATION "${CMAKE_INSTALL_PREFIX}/lib" TYPE DIRECTORY FILES "/home/vagrant/MDM/reindexer/build/cpp_src/pkgconfig")
endif()

if(NOT CMAKE_INSTALL_LOCAL_ONLY)
  # Include the install script for each subdirectory.
  include("/home/vagrant/MDM/reindexer/build/cpp_src/server/cmake_install.cmake")
  include("/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/cmake_install.cmake")
  include("/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_server/cmake_install.cmake")
  include("/home/vagrant/MDM/reindexer/build/cpp_src/doc/cmake_install.cmake")

endif()

