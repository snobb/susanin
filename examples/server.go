package main

/**
 * @author: Alex Kozadaev
 */

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"susanin/pkg/susanin"
	"susanin/pkg/susanin/middleware"
)

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fallback handler\n"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	values, ok := susanin.GetValues(r)
	if !ok {
		log.Println("Empty arguments")
	}

	message := fmt.Sprintf("Hello %s %s\n",
		strings.Title(values["fname"]), strings.Title(values["lname"]))
	w.Write([]byte(message))
}

func helloSplatHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	w.WriteHeader(200)

	values, ok := susanin.GetValues(r)
	if !ok {
		log.Println("Empty arguments")
	}

	message := fmt.Sprintf("Hello %s [uri: %s]\n", values["fname"], uri)
	w.Write([]byte(message))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("home!!!\n"))
}

func main() {
	mux := http.NewServeMux()

	router := susanin.NewSusanin()
	router.Handle("/home/*", homeHandler)
	router.Handle("/short", homeHandler)
	router.Handle("/hello/:fname/:lname/", helloHandler)
	router.Handle("/hello/:fname/*", helloSplatHandler)
	router.Handle("/*", fallbackHandler)

	dh := susanin.DispatchHandler{}
	dh.Attach(middleware.TimerMiddleware)

	mux.HandleFunc("/", dh.Handler(router.Router))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}
