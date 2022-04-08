package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"net/http"
	"text/template"
	"time"
)

var DB *sql.DB

func Checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	openDB()
	router := chi.NewRouter()
	router.Get("/", HomePage)
	//router.Get("/", HomePage)
	router.Post("/post", PostBlog)
	router.Get("/delete/{Id}", DeleteBlog)
	router.Get("/edit/{Id}", EditBlog)
	router.Post("/edit/{Id}", PostEdit)
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

func openDB() {
	db, err := sql.Open("mysql", "root:houseno6@tcp(127.0.0.1:3306)/orevaDB")
	if err != nil {
		fmt.Println(err)
	}
	log.Println("connected to database")
	DB = db

}

func HomePage(w http.ResponseWriter, request *http.Request) {

	//field := "SELECT * FROM blog"
	ans, _ := DB.Query("SELECT * FROM blog")
	//fmt.Println(ans)
	defer ans.Close()

	for ans.Next() {
		var j Blog
		ans.Scan(&j.Id, &j.Author, &j.Title, &j.Content, &j.Time, &j.Date)
		Blogposts = append(Blogposts, j)
	}

	temp := template.Must(template.ParseFiles("index.html"))
	err := temp.Execute(w, Blogposts)
	Checkerror(err)

	Blogposts = nil
}

func PostBlog(w http.ResponseWriter, r *http.Request) {
	InputAuthor := r.FormValue("author")
	InputTitle := r.FormValue("title")
	InputContent := r.FormValue("content")

	now := time.Now()
	m := now.Month()
	d := now.Day()
	y := now.Year()
	hrs := now.Hour()
	min := now.Minute()
	time := fmt.Sprintf("%v:%v	", hrs, min)
	date := fmt.Sprintf("%v %v, %v", m, d, y)
	data = Blog{
		uuid.NewString(),
		InputAuthor,
		InputTitle,
		InputContent,
		time,
		date,
	}
	DB.Query("INSERT INTO blog(id, Author, Title, Content, Time, Date) VALUES (?,?,?,?,?,?)", data.Id, data.Author, data.Title, data.Content, data.Time, data.Date)

	//Blogposts = append(Blogposts, data)
	http.Redirect(w, r, "/", 302)
}

func EditBlog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := chi.URLParam(r, "Id")
	editValue, err := DB.Query("SELECT * FROM blog WHERE id=?", id)

	blog := Blog{}
	for editValue.Next() {
		err := editValue.Scan(&blog.Id, &blog.Author, &blog.Title, &blog.Content, &blog.Time, &blog.Date)
		Checkerror(err)
	}
	log.Println(blog)

	temp := template.Must(template.ParseFiles("edit.html"))
	err = temp.Execute(w, blog)
	Checkerror(err)

}

func PostEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := chi.URLParam(r, "Id")

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
		id,
		InputAuthor,
		InputTitle,
		InputContent,
		time,
		date,
	}
	_, err := DB.Query("UPDATE blog SET Author=?, Title=?, Content=?, Time=?, Date=? WHERE id=?", data.Author, data.Title, data.Content, data.Time, data.Date, data.Id)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/", 302)

}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "Id")
	log.Println("checking Id gotten first", id)
	fmt.Println(id)
	_, err := DB.Query("DELETE FROM blog WHERE id=?", id)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/", 302)

}
