package main

/**
 * @author: Alex Kozadaev
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/middleware/response"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("type=request method=%s uri=%s proto=%s headers=%v",
			r.Method, r.URL.Path, r.Proto, r.Header)

		wbuf := response.NewBuffer(w)
		defer wbuf.Flush()

		next.ServeHTTP(wbuf, r)

		body := wbuf.Body.Bytes()

		if body[len(body)-1] == '\n' {
			body = body[:len(body)-1]
		}

		log.Printf("type=response status=%d headers=%v body=%s time_diff=%s",
			wbuf.Status, wbuf.Header(), string(body), time.Since(start))
	})
}

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fallback handler\n"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if values, ok := framework.GetValues(r.Context()); ok {
		response.WithPayload(r.Context(), map[string]interface{}{
			"first_name": values["fname"],
			"last_name":  values["lname"],
		})
	}
}

func helloSplatHandler(w http.ResponseWriter, r *http.Request) {
	var message string

	uri := r.URL.Path

	if values, ok := framework.GetValues(r.Context()); ok {
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

	fw := framework.New()
	fw.WithPrefix("/api", func() {
		fw.Get("/test2", http.HandlerFunc(homeHandler))

		fw.WithPrefix("/v1", func() {
			fw.Get("/home/*", http.HandlerFunc(homeHandler))
			fw.Get("/test1", http.HandlerFunc(homeHandler))
			fw.Get("/hello/:fname/:lname/", http.HandlerFunc(helloHandler))
			fw.Get("/hello/:fname/*", http.HandlerFunc(helloSplatHandler))
			fw.Get("/*", http.HandlerFunc(fallbackHandler))
			fw.Post("/post/*", http.HandlerFunc(postHandler))
		})
	})

	fw.Get("/", http.HandlerFunc(homeHandler))
	fw.Get("/test3", http.HandlerFunc(homeHandler))

	fw.Attach(response.JSONEncoder, logMiddleware)

	mux.Handle("/", fw)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}
