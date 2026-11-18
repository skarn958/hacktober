package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	flag := os.Getenv("FLAG")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HTBR-FLAG", flag)
		fmt.Fprintln(w, "Welcome to Hacktober!")
	})

	http.ListenAndServe(":1337", nil)
}
