package main

import (
	"strconv"
	//"fmt"
	"code.google.com/p/gorilla/mux"
	"net/http"
	"flag"
	"log"
	"html/template"
	"time"
)

var indexTemplate, err = template.ParseFiles("views/index.html")

func homeHandler(rw http.ResponseWriter, req * http.Request){
	indexTemplate.Execute(rw, struct{Time time.Time}{time.Now()})
}

func main(){
	port := flag.Int("port", 3050, "port to run snippet server")
	flag.Parse()

	log.Println("started server on", *port)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
