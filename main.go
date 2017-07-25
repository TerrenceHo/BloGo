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
		title := "Hello World!"
		t := template.New("index.html")
		t, _ = t.ParseFiles("index.html")
		t.Execute(w, title)
	} else {
		f := "posts/" + r.URL.Path[1:] + ".md"
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[0])
		date := string(lines[1])
		summary := string(lines[2])
		body := strings.Join(lines[3:len(lines)], "\n")
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		post := Post{title, date, summary, body, r.URL.Path[1:]}
		t := template.New("post.html")
		t, _ = t.ParseFiles("post.html")
		t.Execute(w, post)
	}
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	log.Println(files)
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
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

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8000", nil)
}
