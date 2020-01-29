# This file will be configured to contain variables for CPack. These variables
# should be set in the CMake list file of the project before CPack module is
# included. The list of available CPACK_xxx variables and their associated
# documentation may be obtained using
#  cpack --help-variable-list
#
# Some variables are common to all generators (e.g. CPACK_PACKAGE_NAME)
# and some are specific to a generator
# (e.g. CPACK_NSIS_EXTRA_INSTALL_COMMANDS). The generator specific variables
# usually begin with CPACK_<GENNAME>_xxxx.


SET(CPACK_ARCHIVE_COMPONENT_INSTALL "ON")
SET(CPACK_BINARY_7Z "")
SET(CPACK_BINARY_BUNDLE "")
SET(CPACK_BINARY_CYGWIN "")
SET(CPACK_BINARY_DEB "")
SET(CPACK_BINARY_DRAGNDROP "")
SET(CPACK_BINARY_IFW "")
SET(CPACK_BINARY_NSIS "")
SET(CPACK_BINARY_OSXX11 "")
SET(CPACK_BINARY_PACKAGEMAKER "")
SET(CPACK_BINARY_RPM "")
SET(CPACK_BINARY_STGZ "")
SET(CPACK_BINARY_TBZ2 "")
SET(CPACK_BINARY_TGZ "")
SET(CPACK_BINARY_TXZ "")
SET(CPACK_BINARY_TZ "")
SET(CPACK_BINARY_WIX "")
SET(CPACK_BINARY_ZIP "")
SET(CPACK_CMAKE_GENERATOR "Unix Makefiles")
SET(CPACK_COMPONENTS_ALL "dev;server")
SET(CPACK_COMPONENT_UNSPECIFIED_HIDDEN "TRUE")
SET(CPACK_COMPONENT_UNSPECIFIED_REQUIRED "TRUE")
SET(CPACK_DEBIAN_DEV_FILE_NAME "DEB-DEFAULT")
SET(CPACK_DEBIAN_DEV_PACKAGE_DEPENDS "libleveldb-dev,libleveldb-dev,libsnappy-dev,libgoogle-perftools4")
SET(CPACK_DEBIAN_PACKAGE_CONTROL_EXTRA "/home/vagrant/MDM/reindexer/cpp_src/cmd/reindexer_server/contrib/deb/postinst")
SET(CPACK_DEBIAN_PACKAGE_DEPENDS "libleveldb-dev,libsnappy-dev,libgoogle-perftools4")
SET(CPACK_DEBIAN_SERVER_FILE_NAME "DEB-DEFAULT")
SET(CPACK_DEB_COMPONENT_INSTALL "ON")
SET(CPACK_GENERATOR "DEB")
SET(CPACK_INSTALL_CMAKE_PROJECTS "/home/vagrant/MDM/reindexer/build;reindexer;ALL;/")
SET(CPACK_INSTALL_PREFIX "/usr/local")
SET(CPACK_MODULE_PATH "/home/vagrant/MDM/reindexer/cpp_src/cmake/modules")
SET(CPACK_NSIS_DISPLAY_NAME "reindexer 2.4.5.34.g1eeae5fa")
SET(CPACK_NSIS_INSTALLER_ICON_CODE "")
SET(CPACK_NSIS_INSTALLER_MUI_ICON_CODE "")
SET(CPACK_NSIS_INSTALL_ROOT "$PROGRAMFILES")
SET(CPACK_NSIS_PACKAGE_NAME "reindexer 2.4.5.34.g1eeae5fa")
SET(CPACK_OUTPUT_CONFIG_FILE "/home/vagrant/MDM/reindexer/build/CPackConfig.cmake")
SET(CPACK_PACKAGE_CONTACT "Oleg Gerasimov <ogerasimov@gmail.com>")
SET(CPACK_PACKAGE_DEFAULT_LOCATION "/")
SET(CPACK_PACKAGE_DESCRIPTION_FILE "/usr/share/cmake-3.5/Templates/CPack.GenericDescription.txt")
SET(CPACK_PACKAGE_DESCRIPTION_SUMMARY "ReindexerDB server package")
SET(CPACK_PACKAGE_FILE_NAME "reindexer-2.4.5.34.g1eeae5fa-Linux")
SET(CPACK_PACKAGE_INSTALL_DIRECTORY "reindexer 2.4.5.34.g1eeae5fa")
SET(CPACK_PACKAGE_INSTALL_REGISTRY_KEY "reindexer 2.4.5.34.g1eeae5fa")
SET(CPACK_PACKAGE_NAME "reindexer")
SET(CPACK_PACKAGE_RELOCATABLE "true")
SET(CPACK_PACKAGE_VENDOR "Reindexer")
SET(CPACK_PACKAGE_VERSION "2.4.5.34.g1eeae5fa")
SET(CPACK_PACKAGE_VERSION_MAJOR "0")
SET(CPACK_PACKAGE_VERSION_MINOR "1")
SET(CPACK_PACKAGE_VERSION_PATCH "1")
SET(CPACK_RESOURCE_FILE_LICENSE "/home/vagrant/MDM/reindexer/cpp_src/../LICENSE")
SET(CPACK_RESOURCE_FILE_README "/usr/share/cmake-3.5/Templates/CPack.GenericDescription.txt")
SET(CPACK_RESOURCE_FILE_WELCOME "/usr/share/cmake-3.5/Templates/CPack.GenericWelcome.txt")
SET(CPACK_RPM_COMPONENT_INSTALL "ON")
SET(CPACK_RPM_DEV_FILE_NAME "RPM-DEFAULT")
SET(CPACK_RPM_DEV_PACKAGE_REQUIRES "leveldb,leveldb,snappy,gperftools-libs")
SET(CPACK_RPM_PACKAGE_REQUIRES "leveldb,snappy,gperftools-libs")
SET(CPACK_RPM_POST_INSTALL_SCRIPT_FILE "/home/vagrant/MDM/reindexer/cpp_src/cmd/reindexer_server/contrib/rpm/postinst")
SET(CPACK_RPM_RELOCATION_PATHS "/etc")
SET(CPACK_RPM_SERVER_FILE_NAME "RPM-DEFAULT")
SET(CPACK_SET_DESTDIR "ON")
SET(CPACK_SOURCE_7Z "")
SET(CPACK_SOURCE_CYGWIN "")
SET(CPACK_SOURCE_GENERATOR "TBZ2;TGZ;TXZ;TZ")
SET(CPACK_SOURCE_OUTPUT_CONFIG_FILE "/home/vagrant/MDM/reindexer/build/CPackSourceConfig.cmake")
SET(CPACK_SOURCE_TBZ2 "ON")
SET(CPACK_SOURCE_TGZ "ON")
SET(CPACK_SOURCE_TXZ "ON")
SET(CPACK_SOURCE_TZ "ON")
SET(CPACK_SOURCE_ZIP "OFF")
SET(CPACK_SYSTEM_NAME "Linux")
SET(CPACK_TOPLEVEL_TAG "Linux")
SET(CPACK_WIX_SIZEOF_VOID_P "8")

if(NOT CPACK_PROPERTIES_FILE)
  set(CPACK_PROPERTIES_FILE "/home/vagrant/MDM/reindexer/build/CPackProperties.cmake")
endif()

if(EXISTS ${CPACK_PROPERTIES_FILE})
  include(${CPACK_PROPERTIES_FILE})
endif()
