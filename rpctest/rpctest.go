package main

import (
	"github.com/mindtreatstudios/go_rpc/rpctest/api/HelloService"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/mindtreatstudios/go_rpc/json_multipart"

	"log"
	"net/http"
)

func main() {

	// Initialize the rpc server
	apiserver := rpc.NewServer()

	// Initialize the codecs
	jsoncodec := json.NewCodec()
	multipartcodec := json_multipart.NewCodec()

	// Register the codecs (support multipart as well as body format
	apiserver.RegisterCodec(multipartcodec, "multipart/form-data")
	apiserver.RegisterCodec(jsoncodec, "application/json")
	apiserver.RegisterCodec(jsoncodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8

	// Register the services that will handle the requests
	apiserver.RegisterService(new(HelloService.HelloService), "")

	// Create the request router
	muxer := mux.NewRouter()
	muxer.Handle("/api/", apiserver)

	// Start listening for requests
	err := http.ListenAndServe(":1111", muxer)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
