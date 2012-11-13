package main

import (
	"strconv"
	"fmt"
	"code.google.com/p/gorilla/mux"
	"net/http"
	"flag"
)

func homeHandler(rw http.ResponseWriter, req * http.Request){
	rw.Write([]byte("This is an awesome response"))
}

func main(){
	port := flag.Int("port", 3050, "port to run snippet server")
	flag.Parse()

	fmt.Println("started server on", *port)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
