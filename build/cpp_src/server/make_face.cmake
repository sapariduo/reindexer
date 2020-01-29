file(
	DOWNLOAD "https://codeload.github.com/restream/reindexer-face-dist/tar.gz/master"
	"/home/vagrant/MDM/reindexer/build/cpp_src/server/face.tar.gz"
     )
     execute_process(
	 COMMAND "/usr/bin/cmake" -E remove_directory "/home/vagrant/MDM/reindexer/build/cpp_src/server/face"
	 RESULT_VARIABLE ret
     )
     execute_process(
	 COMMAND "/usr/bin/cmake" -E tar xzf "face.tar.gz" WORKING_DIRECTORY /home/vagrant/MDM/reindexer/build/cpp_src/server
	 RESULT_VARIABLE ret
     )
     if (NOT "${ret}" STREQUAL "0")
	 message(FATAL_ERROR "Could not untar 'face.tar.gz'")
     endif()
     file(RENAME "/home/vagrant/MDM/reindexer/build/cpp_src/server/reindexer-face-dist-master" "/home/vagrant/MDM/reindexer/build/cpp_src/server/face")
     file(REMOVE "/home/vagrant/MDM/reindexer/build/cpp_src/server/face.tar.gz")