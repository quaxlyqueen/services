package services

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Serve a directory to a port.
func servePage(router *mux.Router, dest string, path string, port int) {
	router.PathPrefix(path).Handler(http.FileServer(http.Dir(dest)))

	srv := &http.Server{
		Addr: "0.0.0.0: " + string(port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	log.Println("Web server is running on port " + string(port))
}
