package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

var data = Flashcards{
	PageTitle: "LE Card",
	Cards: []Card{
		{"vcle", []string{"bridle", "cable", "able", "cradle", "rifle", "table", "bugle"}},
		{"vccle", []string{"jungle", "scramble", "handle", "tremble", "sample", "single", "muscle"}},
		{"vccle2", []string{"rattle", "bottle", "middle", "paddle", "battle", "scribble", "settle"}},
	},
}

var data1 = Flashcards{
	PageTitle: "L Card",
	Cards: []Card{
		{"le", []string{"angle", "struggle", "saddle", "couple", "needle", "bundle"}},
		{"el", []string{"cancel", "level", "vowel", "novel", "angel", "jewel"}},
		{"il", []string{"council", "stencil", "evil", "pupil", "fossil", "April"}},
		{"al", []string{"pedal", "total", "metal", "local", "journal", "signal"}},
		{"oddball", []string{"special", "fragile"}},
	},
}

func main() {

	pr := newPathResolver()
	pr.Add("GET /hello", hello)
	pr.Add("GET /goodbye", goodbye)
	pr.Add("GET /list", wip)
	pr.Add("GET /slideshow", slideshow)
	pr.Add("* *", generic)

	fmt.Fprint(os.Stdout, data)
	fmt.Fprint(os.Stdout, data1)

	http.ListenAndServe(":8080", pr)

}

func generic(res http.ResponseWriter, req *http.Request) {

	path := req.URL.Path
	parts := strings.Split(path, "/")

	file := parts[1]
	//view := parts[2]

	filecontents, err := ioutil.ReadFile(file + ".yml")
	if err != nil {
		res.WriteHeader(500)
	}

	type Data struct{}
	d := Data{}
	err = yaml.Unmarshal(filecontents, d)
	b, err := yaml.Marshal(d)
	res.Write(b)

}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}

type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func (p *pathResolver) Add(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, check); ok && err == nil {
			handlerFunc(res, req)
			return
		} else if err != nil {
			fmt.Fprint(res, err)
		}
	}
}

func wip(res http.ResponseWriter, req *http.Request) {

	t := FlashcardSet{[]Flashcards{data, data1}}

	d, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}
	res.Write(d)
}

type Card struct {
	Name string
	Word []string
}

type Flashcards struct {
	PageTitle string
	Cards     []Card
}

type FlashcardSet struct {
	Flashcards []Flashcards
}

func hello(res http.ResponseWriter, req *http.Request) {

	data, data1 := flashcard()

	tmpl := template.Must(template.ParseFiles("layout.html"))

	tmpl.Execute(res, data)
	tmpl.Execute(res, data1)

}

func flashcard() (Flashcards, Flashcards) {
	data := Flashcards{
		PageTitle: "LE Card",
		Cards: []Card{
			{"vcle", []string{"bridle", "cable", "able", "cradle", "rifle", "table", "bugle"}},
			{"vccle", []string{"jungle", "scramble", "handle", "tremble", "sample", "single", "muscle"}},
			{"vccle2", []string{"rattle", "bottle", "middle", "paddle", "battle", "scribble", "settle"}},
		},
	}

	data1 := Flashcards{
		PageTitle: "L Card",
		Cards: []Card{
			{"le", []string{"angle", "struggle", "saddle", "couple", "needle", "bundle"}},
			{"el", []string{"cancel", "level", "vowel", "novel", "angel", "jewel"}},
			{"il", []string{"council", "stencil", "evil", "pupil", "fossil", "April"}},
			{"al", []string{"pedal", "total", "metal", "local", "journal", "signal"}},
			{"oddball", []string{"special", "fragile"}},
		},
	}
	return data, data1
}

func slideshow(res http.ResponseWriter, req *http.Request) {

	flashcards := FlashcardSet{[]Flashcards{data, data1}}
	tmpl := template.Must(template.ParseFiles("slideshow.html"))

	tmpl.Execute(res, flashcards)

}

func goodbye(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")
	name := parts[2]
	if name == "" {
		name = "Monika"
	}

	fmt.Fprint(res, "Hej Då ", name)
	fmt.Fprint(res, ".")

}
