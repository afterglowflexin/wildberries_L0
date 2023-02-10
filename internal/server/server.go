package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/afterglowflexin/wildberries/level0/internal/cache"
)

type serverHTTP struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *serverHTTP {
	return &serverHTTP{cache}
}

func (s *serverHTTP) Start() {
	http.HandleFunc("/", s.Serve)
	log.Println("[DEBUG] starting http server")
	http.ListenAndServe(":8080", nil)
}

func (s *serverHTTP) Serve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := "../internal/static/index.html"

		http.ServeFile(w, r, path)
	case "POST":
		r.ParseMultipartForm(0)
		id := r.FormValue("message")

		log.Printf("[DEBUG] getting order with ID %s", id)

		order, err := s.cache.GetOrder(id)

		if err != nil {
			log.Printf("[ERROR] fetching product, %s", err)
			fmt.Fprintf(w, "No product with this ID")
			return
		}

		fmt.Fprintf(w, "Order data is %s", order)
	}
}
