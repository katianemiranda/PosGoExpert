package main

import (
	"log"
	"net/http"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				//debug.PrintStack()
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	log.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", recoverMiddleware(mux)); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", ":3000", err)
	}

}
