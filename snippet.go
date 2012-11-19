package main

import (
	"html/template"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	_ "log"
	"os/exec"
	"io"
	"net/http"
	"net/url"
)

const runUrl = "http://golang.org/compile?output=json"

//Stores a single instance of a code snippet
type Snippet struct {
	Id              bson.ObjectId `bson:"_id"`
	Name            string
	Description     string
	Code            string
	HighlightedCode template.HTML
}

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

func (s *Snippet) run(w io.Writer) {
	v := url.Values{}
	v.Set("body", s.Code)
	resp, _ := http.PostForm(runUrl, v)
	io.Copy(w, resp.Body)
}

func NewSnippet(name, description, code string) *Snippet {
	return &Snippet{Id: bson.NewObjectId(), Name: name, Description: description, Code: code}
}

//Returns a list of all snippets stored in the database
//TODO: add pagination
func AllSnippets() (snippets []Snippet) {
	session.DB("gosnip").C("snippets").Find(nil).Iter().All(&snippets)
	return
}

func FindSnippet(id interface{}) (snippet *Snippet) {
	session.DB("gosnip").C("snippets").FindId(id).One(&snippet)
	return snippet
}

func CreateSnippet(snippet *Snippet) {
	snippet.Pygmentize()
	session.DB("gosnip").C("snippets").Insert(snippet)
}
