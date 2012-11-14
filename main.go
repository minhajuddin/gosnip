package main

import (
	"strconv"

	"flag"
	"log"

	"net/http"
	"html/template"

	"code.google.com/p/gorilla/mux"
)

type Snippet struct{
	Name string
	Description string
	Code string
}

var data []Snippet
var indexTemplate, err = template.ParseFiles("views/index.html")

func homeHandler(rw http.ResponseWriter, req * http.Request){
	indexTemplate.Execute(rw, data)
}

func main(){
	data = make([]Snippet, 100)
	data[0] = Snippet{"helloworld", "Hello world snip", `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello World")
		}`}
	port := flag.Int("port", 3000, "port to run snippet server")
	flag.Parse()

	log.Println("started server on", *port)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
