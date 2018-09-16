package main

import (
	"fmt"
	"log"
	"net/http"

	"httprouter/pkg/susanin"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	values, ok := susanin.GetValues(r)
	if !ok {
		log.Println("Empty arguments")
	}

	message := fmt.Sprintf("hello %s\n", values["name"])
	w.Write([]byte(message))
}

func helloSplatHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	w.WriteHeader(200)

	values, ok := susanin.GetValues(r)
	if !ok {
		log.Println("Empty arguments")
	}

	message := fmt.Sprintf("hello %s [uri: %s]\n", values["name"], uri)
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
	router.Handle("/hello/:name", helloHandler)
	router.Handle("/hello/:name/*", helloSplatHandler)

	dh := susanin.DispatchHandler{}
	dh.Attach(susanin.TimerMiddleware)

	mux.HandleFunc("/", dh.Handler(router.Router))
	http.ListenAndServe(":8080", mux)
}
