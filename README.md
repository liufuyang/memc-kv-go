# Http Key Value Storage
Kind of a cache

## Http server

## Planning
- [x] Step 1: a simple http server
- [x] Step 2: support simple get/set
- [x] Step 3: modularize the repo
- [x] Step 4: benchmark things somehow
- [x] Step 5: add some metrics
- [ ] Step 6: http endpoint unit tests
- [ ] Step 7: add gRPC endpoint
- [ ] Step 8: add memcached protocol
- [ ] A rust version?

## Run test or start server locally
In the project folder:
```
go test -v ./...

go run -race main.go
```

## Manual test and benchmark
```
curl -X POST http://localhost:4000/set?somekey=somevalue
curl -X GET http://localhost:4000/get?key=somekey
```

```
ab -p empty.txt -n 10 -c 10 127.0.0.1:4000/set?somekey=somevalue
ab -n 10 -c 10 127.0.0.1:4000/get?key=somekey 
```

Or using a benchmark repo at https://github.com/liufuyang/autocannon-go
```
go build && ./autocannon-go --connections=20 --pipelining=10 --duration=300 --uri=http://localhost:4000
```

## GRPC server

```
grpcurl -plaintext -d '{"key":"k1", "value":"value1", "ttl_seconds": 4}'  localhost:4400 kv.Cache/Set 
grpcurl -plaintext -d '{"key":"k1", "value":"value1"}'  localhost:4400 kv.Cache/Set

grpcurl -plaintext -d '{"key":"k1"}'  localhost:4400 kv.Cache/Get
```

## memcached interface

For comparing with memcached
```
docker run --name mc -d --rm -p 11211:11211 memcached memcached -m 64
```

A bench test
```
docker run --rm redislabs/memtier_benchmark --protocol=memcache_text --server 192.168.0.25 --port=6001 --generate-keys -n 1000 --key-maximum=10000 --ratio=2:8
```
