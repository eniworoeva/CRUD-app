package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"text/template"
	"time"
)

func Checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	router := chi.NewRouter()
	router.Get("/", HomePage)
	router.Post("/post", PostBlog)
	router.Get("/delete/{Id}", DeleteBlog)
	router.Get("/edit/{Id}", EditBlog)
	fmt.Println("Listening!")
	log.Fatal(http.ListenAndServe(":2020", router))
}

type Blog struct {
	Id      string
	Author  string
	Title   string
	Content string
	Time    string
	Date    string
}

var Blogposts []Blog
var data Blog

func HomePage(w http.ResponseWriter, request *http.Request) {
	temp := template.Must(template.ParseFiles("homepage.html"))
	err := temp.Execute(w, Blogposts)
	Checkerror(err)
}
func EditBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "Id")
	for i, _ := range Blogposts {
		if id == Blogposts[i].Id {
			data = Blogposts[i]
			temp := template.Must(template.ParseFiles("edit.html"))
			err := temp.Execute(w, data)
			Checkerror(err)
			Blogposts = append(Blogposts[:i], Blogposts[i+1:]...)
		}
	}

}
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "Id")
	log.Println("checking Id gotten first", id)
	fmt.Println(id)
	for i, item := range Blogposts {
		fmt.Println(item.Id)
		if id == item.Id {
			log.Println("checking Id gotten in the loop", id)
			Blogposts = append(Blogposts[:i], Blogposts[i+1:]...)
		}
		log.Println("checking resultant Blogposts", Blogposts)
	}
	http.Redirect(w, r, "/", 302)
}
func PostBlog(w http.ResponseWriter, r *http.Request) {
	InputAuthor := r.FormValue("author")
	InputTitle := r.FormValue("title")
	InputContent := r.FormValue("content")

	now := time.Now()
	m := now.Month()
	d := now.Day()
	hrs := now.Hour()
	min := now.Minute()
	time := fmt.Sprintf("%v:%v", hrs, min)
	date := fmt.Sprintf("%v %v", m, d)
	data = Blog{
		uuid.NewString(),
		InputAuthor,
		InputTitle,
		InputContent,
		time,
		date,
	}
	Blogposts = append(Blogposts, data)
	temp := template.Must(template.ParseFiles("homepage.html"))
	err := temp.Execute(w, Blogposts)
	Checkerror(err)
}
