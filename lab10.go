package main

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

type Book struct {
    ID     int
    Title  string
    Author string
}

var db *sql.DB
var tpl *template.Template

func init() {
    tpl = template.Must(template.ParseGlob("templates/*.html"))

    var err error
    db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
}

func createHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        author := r.FormValue("author")

        _, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", title, author)
        if err != nil {
            log.Fatal(err)
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        tpl.ExecuteTemplate(w, "create.html", nil)
    }
}

func readHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM books")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    books := []Book{}
    for rows.Next() {
        book := Book{}
        err := rows.Scan(&book.ID, &book.Title, &book.Author)
        if err != nil {
            log.Fatal(err)
        }
        books = append(books, book)
    }

    tpl.ExecuteTemplate(w, "read.html", books)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        id := r.FormValue("id")
        title := r.FormValue("title")
        author := r.FormValue("author")

        _, err := db.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", title, author, id)
        if err != nil {
            log.Fatal(err)
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        id := r.FormValue("id")

        var book Book
        err := db.QueryRow("SELECT * FROM books WHERE id = ?", id).Scan
