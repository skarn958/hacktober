package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Challenge running")
	flag := os.Getenv("FLAG")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		w.Header().Add("HTBR-FLAG", flag)
		fmt.Fprintln(w, "Welcome to Hacktober!")
	})

	http.ListenAndServe(":1337", nil)
}
