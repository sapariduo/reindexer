# CMake generated Testfile for 
# Source directory: /home/vagrant/MDM/reindexer/connectors/py_reindexer
# Build directory: /home/vagrant/MDM/reindexer/build/connectors/py_reindexer
# 
# This file includes the relevant testing commands required for 
# testing this directory and lists subdirectories to be tested as well.
add_test(pyreindexer_test "python3" "-m" "unittest")
set_tests_properties(pyreindexer_test PROPERTIES  ENVIRONMENT "PYTHONPATH=/home/vagrant/MDM/reindexer/connectors/py_reindexer:/home/vagrant/MDM/reindexer/build/connectors/py_reindexer" WORKING_DIRECTORY "/home/vagrant/MDM/reindexer/connectors/py_reindexer/tests")
