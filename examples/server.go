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
	if _, err := w.Write([]byte("fallback handler\n")); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if values, ok := framework.GetValues(r.Context()); ok {
		response := response.New(w)
		err := response.Payload(r.Context(), map[string]interface{}{
			"first_name": values["fname"],
			"last_name":  values["lname"],
		})
		if err != nil {
			w.WriteHeader(500)
		}
	}
}

func helloSplatHandler(w http.ResponseWriter, r *http.Request) {
	var message string

	uri := r.URL.Path

	if values, ok := framework.GetValues(r.Context()); ok {
		message = fmt.Sprintf("Hello %s [uri: %s]\n", values["fname"], uri)
	}

	if _, err := w.Write([]byte(message)); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("home!!!\n")); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	if _, err := w.Write([]byte(fmt.Sprintf("response: %v\n", string(bytes)))); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func main() {
	fw := framework.New()
	fw = fw.WithNotFoundHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"status":"Endpoint not found"}\n`))
		}),
	)

	fw.Get("/", http.HandlerFunc(homeHandler))
	fw.Get("/test3", http.HandlerFunc(homeHandler))

	fw.WithDefaultPrefix("/api")
	fw.Get("/test2", http.HandlerFunc(homeHandler))

	fw.WithPrefix("/v1", func() {
		fw.Get("/home/*", http.HandlerFunc(homeHandler))
		fw.WithPrefix("/test", func() {
			fw.Get("/1", http.HandlerFunc(homeHandler))
			fw.Get("/2", http.HandlerFunc(homeHandler))
			fw.Get("/3", http.HandlerFunc(homeHandler))
		})
		fw.Get("/hello/:fname/:lname/", http.HandlerFunc(helloHandler))
		fw.Get("/hello/:fname/*", http.HandlerFunc(helloSplatHandler))
		fw.Get("/*", http.HandlerFunc(fallbackHandler))
		fw.Post("/post/*", http.HandlerFunc(postHandler))
	})

	fw.Attach(logMiddleware)

	err := http.ListenAndServe(":8080", fw)
	if err != nil {
		log.Println(err.Error())
	}
}
