file(
	DOWNLOAD "https://codeload.github.com/swagger-api/swagger-ui/tar.gz/2.x"
	"/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger.tar.gz"
     )
     execute_process(
       COMMAND "/usr/bin/cmake" -E tar xzf "swagger.tar.gz" WORKING_DIRECTORY /home/vagrant/MDM/reindexer/build/cpp_src/server
       RESULT_VARIABLE ret
     )
     if (NOT "${ret}" STREQUAL "0")
	 message(FATAL_ERROR "Could not untar 'swagger.tar.gz'")
     endif()
     execute_process(
       COMMAND "/usr/bin/cmake" -E copy_directory "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger-ui-2.x/dist" "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger"
       RESULT_VARIABLE ret
     )
     if (NOT "${ret}" STREQUAL "0")
	 message(FATAL_ERROR "Could not copy directory '/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger-ui-2.x/dist'")
     endif()
     execute_process(
       COMMAND "/usr/bin/cmake" -E copy "/home/vagrant/MDM/reindexer/cpp_src/server/contrib/server.yml" "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger/swagger.yml"
       RESULT_VARIABLE ret
     )
     if (NOT "${ret}" STREQUAL "0")
	 message(FATAL_ERROR "Could not copy '/home/vagrant/MDM/reindexer/cpp_src/server/contrib/server.yml'")
     endif()
     file(RENAME "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger/swagger-ui.min.js" "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger/swagger-ui.js")
     execute_process(COMMAND "/usr/bin/cmake" -P "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger_replace.cmake")
     execute_process(
       COMMAND /usr/bin/cmake -E remove_directory "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger-ui-2.x"
       RESULT_VARIABLE ret
     )
     file(REMOVE "/home/vagrant/MDM/reindexer/build/cpp_src/server/swagger.tar.gz")