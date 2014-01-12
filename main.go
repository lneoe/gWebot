package main

import (
	"gobot/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/home2", handler.ReciveMsgHandler)
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
