package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sqladmin/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			panic(err)
		}

	})
	http.ListenAndServe(":8700", mux)
}
