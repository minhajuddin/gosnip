package main

import (
	"code.google.com/p/gorilla/pat"
	"flag"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"io/ioutil"
	"strconv"
)

const databaseName = "gosnip"

type context struct {
	Request        *http.Request
	ResponseWriter *http.ResponseWriter
	DbSession      *mgo.Session
	Database       *mgo.Database
  Params map[string]string
}

func createContext(w http.ResponseWriter, r * http.Request) *context {
  dbSession := session.Clone()
  ctx := context{
    Request: r,
    ResponseWriter: &w,
    DbSession: dbSession,
    Database: dbSession.DB(databaseName),
    Params: make(map[string]string,1)}
  return &ctx
}

func httpHandler(w http.ResponseWriter, r * http.Request) {
  //ctx := createContext(w, r)
}


//TODO: should be cloning a session for each request handler, probably put it in a title
var session *mgo.Session

var funcs = template.FuncMap{"appVersion": func() string { return appVersion }}
var templates = template.Must(template.New("base").Funcs(funcs).ParseGlob("views/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.Lookup(tmpl+".html").Funcs(funcs).Execute(w, data)
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

func compileHandler(w http.ResponseWriter, r *http.Request) {
	id := getObjectId(getParam("id", r))
	s := FindSnippet(id)
	if s == nil {
		http.NotFound(w, r)
		return
	}
	s.run(w)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", AllSnippets())
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about", nil)
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
	r.Get("/about", aboutHandler)
	r.Post("/create", createHandler)
	r.Post("/compile/{id}", compileHandler)
	r.Get("/", indexHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	return r
}

var err error
var appVersion = ""

func main() {
	appVersion = findAppVersion()
	session, err = mgo.Dial("localhost")
	exitIfError(err)
	defer session.Close()
	port := flag.Int("port", 3000, "port to run snippet server")
	flag.Parse()
	log.Println("running V", appVersion)
	log.Println("started server on", *port)
	http.Handle("/", router())
	exitIfError(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func findAppVersion() string {
	gitVersion, _ := ioutil.ReadFile("REVISION")
	return string(gitVersion)
}

//helpers
func exitIfError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
