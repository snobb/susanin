package main

/**
 * @author: Alex Kozadaev
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware/request"
	"github.com/snobb/susanin/pkg/middleware/response"
)

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fallback handler\n"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if values, ok := framework.GetValues(r); ok {
		response.WithPayload(r.Context(), map[string]interface{}{
			"first_name": values["fname"],
			"last_name":  values["lname"],
		})
	}
}

func helloSplatHandler(w http.ResponseWriter, r *http.Request) {
	var message string

	uri := r.URL.Path

	if values, ok := framework.GetValues(r); ok {
		message = fmt.Sprintf("Hello %s [uri: %s]\n", values["fname"], uri)
	}

	w.WriteHeader(200)
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

	logger := logging.New("example", os.Stderr)

	fw := framework.NewFrameworkWithPrefix("/api/v1")
	fw.Get("/", homeHandler)
	fw.Get("/home/*", homeHandler)
	fw.Get("/short", homeHandler)
	fw.Get("/test1", homeHandler)
	fw.Get("/test2", homeHandler)
	fw.Get("/test3", homeHandler)
	fw.Get("/test4", homeHandler)
	fw.Get("/test5", homeHandler)
	fw.Get("/hello/:fname/:lname/", helloHandler)
	fw.Get("/hello/:fname/*", helloSplatHandler)
	fw.Get("/*", fallbackHandler)
	fw.Post("/post/*", postHandler)

	// fw.Attach(middleware.Debug)
	fw.Attach(request.NewLogger(logger))
	fw.Attach(response.NewJSONEncoder(logger))
	fw.Attach(response.NewLogger(logger))
	fw.Attach(response.NewTimer(logger))

	mux.Handle("/", fw.Router())
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}
