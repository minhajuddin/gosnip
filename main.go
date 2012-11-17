package main

import (
	"code.google.com/p/gorilla/pat"
	"flag"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
)

//TODO: should be cloing a session for each request handler, probably put it in a title
var session *mgo.Session

var templates = template.Must(template.ParseGlob("views/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

func getObjectId(str *string) *bson.ObjectId {
	if str != nil && len(*str) == 24 {
		oid := bson.ObjectIdHex(*str)
		if oid.Valid() {
			return &oid
		}
	}
	return nil
}

func getParam(name string, r *http.Request) *string {
	val := (r.URL.Query().Get(":" + name))
	return &val
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	id := getObjectId(getParam("id", r))
	s := FindSnippet(id)
	if s == nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "show", s)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", AllSnippets())
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	s := NewSnippet(r.FormValue("name"), r.FormValue("description"), r.FormValue("code"))
	go CreateSnippet(s)
	//TODO: should redirect to the snippet details view
	http.Redirect(w, r, "/", 302)
}

func router() *pat.Router {
	r := pat.New()
	r.Get("/new", newHandler)
	r.Get("/show/{id}", showHandler)
	r.Post("/create", createHandler)
	r.Get("/", indexHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	return r
}

var err error

func main() {
	session, err = mgo.Dial("localhost")
	exitIfError(err)
	defer session.Close()
	port := flag.Int("port", 3000, "port to run snippet server")
	flag.Parse()
	log.Println("started server on", *port)
	http.Handle("/", router())
	exitIfError(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

//helpers
func exitIfError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
