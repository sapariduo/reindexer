if("b50b7b7a8c2835b45607cffabc18e4742072e9e6" STREQUAL "")
  message(FATAL_ERROR "Tag for git checkout should not be empty.")
endif()

set(run 0)

if("/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitinfo.txt" IS_NEWER_THAN "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitclone-lastrun.txt")
  set(run 1)
endif()

if(NOT run)
  message(STATUS "Avoiding repeated git clone, stamp file is up to date: '/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitclone-lastrun.txt'")
  return()
endif()

execute_process(
  COMMAND ${CMAKE_COMMAND} -E remove_directory "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to remove directory: '/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib'")
endif()

# try the clone 3 times incase there is an odd git clone issue
set(error_code 1)
set(number_of_tries 0)
while(error_code AND number_of_tries LESS 3)
  execute_process(
    COMMAND "/usr/bin/git" clone --origin "origin" "https://github.com/Restream/replxx" "replxx_lib"
    WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src"
    RESULT_VARIABLE error_code
    )
  math(EXPR number_of_tries "${number_of_tries} + 1")
endwhile()
if(number_of_tries GREATER 1)
  message(STATUS "Had to git clone more than once:
          ${number_of_tries} times.")
endif()
if(error_code)
  message(FATAL_ERROR "Failed to clone repository: 'https://github.com/Restream/replxx'")
endif()

execute_process(
  COMMAND "/usr/bin/git" checkout b50b7b7a8c2835b45607cffabc18e4742072e9e6
  WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to checkout tag: 'b50b7b7a8c2835b45607cffabc18e4742072e9e6'")
endif()

execute_process(
  COMMAND "/usr/bin/git" submodule init 
  WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to init submodules in: '/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib'")
endif()

execute_process(
  COMMAND "/usr/bin/git" submodule update --recursive 
  WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to update submodules in: '/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib'")
endif()

# Complete success, update the script-last-run stamp file:
#
execute_process(
  COMMAND ${CMAKE_COMMAND} -E copy
    "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitinfo.txt"
    "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitclone-lastrun.txt"
  WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib"
  RESULT_VARIABLE error_code
  )
if(error_code)
  message(FATAL_ERROR "Failed to copy script-last-run stamp file: '/home/vagrant/MDM/reindexer/build/cpp_src/cmd/reindexer_tool/replxx_lib-prefix/src/replxx_lib-stamp/replxx_lib-gitclone-lastrun.txt'")
endif()

