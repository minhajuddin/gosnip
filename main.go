package main

import (
	"strconv"
	"flag"
	"log"
	"html/template"
	"net/http"
	"code.google.com/p/gorilla/mux"
	"labix.org/v2/mgo"
)

var session *mgo.Session

var indexTemplate, err = template.ParseFiles("views/index.html")

func homeHandler(rw http.ResponseWriter, req * http.Request){
	indexTemplate.Execute(rw, AllSnippets())
}

func createHandler(rw http.ResponseWriter, req * http.Request){
	s := NewSnippet(req.FormValue("name"), req.FormValue("description"), req.FormValue("code"))
	go CreateSnippet(s)
	//TODO: should redirect to the snippet details view
	http.Redirect(rw, req, "/", 302)
}


func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/create", createHandler).Methods("POST")
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	return r
}

func main(){
	session, err = mgo.Dial("localhost")
	exitIfError(err)
	defer session.Close()
	port := flag.Int("port", 3000, "port to run snippet server")
	flag.Parse()
	log.Println("started server on", *port)
	http.Handle("/", router())
	err := http.ListenAndServe(":" + strconv.Itoa(*port), nil)
	exitIfError(err)
}
//helpers
func exitIfError(err error){
	if err != nil {
		log.Panic(err)
	}
}
