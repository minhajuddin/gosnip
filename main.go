package main

import (
	"strconv"

	"flag"
	"log"

	"html/template"
	"net/http"

	"os/exec"
	"io/ioutil"

	"code.google.com/p/gorilla/mux"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"

)

func (s *Snippet) Pygmentize() {
	pygmentCmd := exec.Command("bash", "-c", "pygmentize -l go -f html")
	pygmentIn, _ := pygmentCmd.StdinPipe()
	pygmentOut, _ := pygmentCmd.StdoutPipe()
	pygmentCmd.Start()
	pygmentIn.Write([]byte(s.Code))
	pygmentIn.Close()
	highlightedCodeBytes, _ := ioutil.ReadAll(pygmentOut)
	s.HighlightedCode = template.HTML(highlightedCodeBytes)
	pygmentCmd.Wait()
}


//Stores a single instance of a code snippet
type Snippet struct{
	Name string
	Description string
	Code string
	HighlightedCode template.HTML
}

//Returns a list of all snippets stored in the database
//TODO: add pagination
func AllSnippets()(snippets []Snippet) {
	session.DB("gosnip").C("snippets").Find(nil).Iter().All(&snippets)
	return
}

func CreateSnippet(snippet *Snippet) {
	snippet.Pygmentize()
	session.DB("gosnip").C("snippets").Insert(snippet)
}

var session *mgo.Session

var indexTemplate, err = template.ParseFiles("views/index.html")

func homeHandler(rw http.ResponseWriter, req * http.Request){
	indexTemplate.Execute(rw, AllSnippets())
}

func createHandler(rw http.ResponseWriter, req * http.Request){
	s := &Snippet{Name: req.FormValue("name"), Description: req.FormValue("description"), Code: req.FormValue("code") }
	go CreateSnippet(s)
	http.Redirect(rw, req, "/", 302)
}

func main(){
	session, err = mgo.Dial("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	port := flag.Int("port", 3000, "port to run snippet server")
	flag.Parse()

	log.Println("started server on", *port)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/create", createHandler).Methods("POST")
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.Handle("/", r)

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
