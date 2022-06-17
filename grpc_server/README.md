# grpc_server

## Build
From `golang` folder
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc_server/protos/kv.proto
```