package http_server

import (
	"example.com/http_kv/cache"
	"example.com/http_kv/metrics"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	cache cache.Cache
}

func NewServer(c cache.Cache) *Server {
	metrics.RegisterCacheSizeGauge(c)
	return &Server{cache: c}
}

func (s *Server) Start(port *int) {
	address := fmt.Sprintf(":%d", *port)
	log.Printf("HTTP server with %v is listening on address: %v\n", s.cache.Name(), address)

	setChain := metrics.CreateHttpHandleChain(s.set, "set", s.cache.Name(), validationMiddleware)
	getChain := metrics.CreateHttpHandleChain(s.get, "get", s.cache.Name(), validationMiddleware)
	sizeChain := metrics.CreateHttpHandleChain(s.size, "size", s.cache.Name(), validationMiddleware)

	serveMux := http.NewServeMux()
	serveMux.Handle("/set", setChain)
	serveMux.Handle("/get", getChain)
	serveMux.Handle("/size", sizeChain)
	err := http.ListenAndServe(address, serveMux)
	if err != nil {
		panic(err)
	}
}

func (s *Server) set(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Only POST method is supported for set", http.StatusMethodNotAllowed)
		return
	}

	for k, v := range req.URL.Query() {
		if len(v) > 1 {
			http.Error(w, "Only one value is allowed for each key", http.StatusBadRequest)
			return
		}
		if len(k) > 100 {
			http.Error(w, "Key length must be less than 100", http.StatusMethodNotAllowed)
			return
		}
		s.cache.Set(k, v[0])
		// fmt.Printf("k: %s, v: %s\n", k, v[0])
		fmt.Fprintf(w, "set key: %s, value: %s \n", k, v[0])
	}
}

func (s *Server) get(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Only GET method is supported for get", http.StatusMethodNotAllowed)
		return
	}

	key := req.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Only GET method is supported for get", http.StatusMethodNotAllowed)
		return
	}
	if len(key) > 100 {
		http.Error(w, "Key length must be less than 100", http.StatusMethodNotAllowed)
		return
	}
	v := s.cache.Get(key)
	fmt.Fprintf(w, "%s\n", v)
}

func (s *Server) size(w http.ResponseWriter, req *http.Request) {
	size := s.cache.Size()
	fmt.Fprintf(w, "%d\n", size)
}

// validationMiddleware validates request url length and body length.
func validationMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 504 {
			http.Error(w, "Key and value length must be less than 500 char", http.StatusBadRequest)
			return
		}
		if len(r.URL.Path) < 2 {
			http.Error(w, "Must provide a key in the path", http.StatusBadRequest)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 32)
		f.ServeHTTP(w, r)
	}
}
