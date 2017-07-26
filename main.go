package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {
		posts := getPosts()
		t := template.New("index.html")
		t, err := t.ParseFiles("index.html")
		must(err)
		t.Execute(w, posts)
	} else {
		f := "posts/" + r.URL.Path[1:] + ".md"
		log.Println(f)
		fileread, err := ioutil.ReadFile(f)
		must(err)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		post := Post{title, date, summary, body, r.URL.Path[1:]}
		t := template.New("post.html")
		t, err = t.ParseFiles("post.html")
		must(err)
		t.Execute(w, post)
	}
}

func getPosts() []Post {
	a := []Post{}
	files, err := filepath.Glob("posts/*")
	must(err)
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, err := ioutil.ReadFile(f)
		must(err)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		a = append(a, Post{title, date, summary, body, file})
	}
	return a
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Serving on port :8000")
	http.ListenAndServe(":8000", nil)
}
