package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
	"os"
	"os/signal"
	"context"
)

var DEST string
var PATH string
var PORT string

// Serve a directory to a port.
func serve() {
	router := mux.NewRouter()
	router.PathPrefix(PATH).Handler(http.FileServer(http.Dir(DEST)))

	srv := &http.Server{
		Addr: "0.0.0.0: " + PORT,
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

	log.Println("Web server is running on port " + PORT)
	log.Println("To stop, press Control + C.")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30))
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down web server on port " + PORT)
	os.Exit(0)
}

func parseCLI() {
	pflag.StringP("destination", "d", "/", "The new root directory.")
	pflag.StringP("path", "p", ".", "The webpage directory to expose.")
	pflag.StringP("port", "o", "8000", "The port to expose the webpage on.")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	DEST = viper.GetString("destination")
	PATH = viper.GetString("path")
	PORT = viper.GetString("port")
}

func main() {
	parseCLI()
	serve()
}
