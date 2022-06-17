package main

import (
	"example.com/http_kv/cache"
	"example.com/http_kv/grpc_server"
	"example.com/http_kv/http_server"
	"example.com/http_kv/metrics"
	"flag"
	"sync"
)

var (
	defaultTtlSeconds = flag.Uint("ttl", 20, "Default item TTL in seconds")
	grpcPort          = flag.Int("grpc-port", 4400, "The gRPC server port")
	httpPort          = flag.Int("http-port", 4000, "The HTTP server port")
	httpPortSyncMap   = flag.Int("http-port-sync", 4001, "The HTTP sync-map-cache server port")
)

func main() {
	flag.Parse()
	metrics.InitMetrics(":9000") // metrics port

	var wg sync.WaitGroup
	var c1 cache.Cache = cache.NewVacuumedStdMapCache(*defaultTtlSeconds)
	server := http_server.NewServer(c1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start(httpPort) // start service at port 4000 with StdMapCache
	}()

	var c2 = cache.NewVacuumedSyncMapCache(*defaultTtlSeconds)
	server2 := http_server.NewServer(c2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		server2.Start(httpPortSyncMap) // start service at port 4000 with SyncMapCache
	}()

	// TODO add gRPC server
	grpcServer := grpc_server.NewServer(c1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Start(grpcPort)
	}()
	wg.Wait()
}
