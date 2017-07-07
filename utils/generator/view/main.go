package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../../../docs/"))
	http.Handle("/", fs)
	log.Println("Listening on :3000...")
	http.ListenAndServe(":3000", nil)
}
