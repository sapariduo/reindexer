file(READ /home/vagrant/MDM/reindexer/build/cpp_src/server/swagger/index.html indexhtml)
    string(REPLACE "http://petstore.swagger.io/v2/swagger.json" "swagger.yml" indexhtml "${indexhtml}")
    file(WRITE /home/vagrant/MDM/reindexer/build/cpp_src/server/swagger/index.html "${indexhtml}")