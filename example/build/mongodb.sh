docker run -d --name example-mongo \
    -e MONGO_INITDB_ROOT_USERNAME=example \
    -e MONGO_INITDB_ROOT_PASSWORD=example \
    -p 27017:27017 \
    mongo