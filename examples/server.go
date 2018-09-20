package main

/**
 * @author: Alex Kozadaev
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"susanin/pkg/susanin/framework"
	"susanin/pkg/susanin/middleware"
	"susanin/pkg/susanin/router"
)

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fallback handler\n"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	values, ok := router.GetValues(r)
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

	values, ok := router.GetValues(r)
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

func postHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("response: %v\n", string(bytes))))
}

func main() {
	mux := http.NewServeMux()

	fw := framework.NewFramework()
	fw.Get("/", homeHandler)
	fw.Get("/home/*", homeHandler)
	fw.Get("/short", homeHandler)
	fw.Get("/hello/:fname/:lname/", helloHandler)
	fw.Get("/hello/:fname/*", helloSplatHandler)
	fw.Get("/*", fallbackHandler)
	fw.Post("/post/*", postHandler)

	fw.AttachMiddleware(middleware.TimerMiddleware)

	mux.HandleFunc("/", fw.Router())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}
