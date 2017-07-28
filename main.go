package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("guestbook").Parse(guestbookTmpl))
		db, err := sql.Open("postgres", "")
		if err != nil {
			log.Printf("opening db: %v", err)
			http.Error(w, "Internal server error", 500)
			return
		}
		if r.Method == "POST" {
			query := "INSERT INTO guestbook (name, comment, created_at) VALUES ($1, $2, now())"
			if _, err := db.Exec(query, r.FormValue("name"), r.FormValue("comment")); err != nil {
				log.Printf("inserting into db: %v", err)
				http.Error(w, "Internal server error", 500)
				return
			}
		}
		query := "SELECT name, comment, created_at FROM guestbook ORDER BY created_at DESC"
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("getting from db: %v", err)
			http.Error(w, "Internal server error", 500)
			return
		}
		type comment struct {
			Name      string
			Comment   string
			CreatedAt time.Time
		}
		var comments []comment
		for rows.Next() {
			var c comment
			if err := rows.Scan(&c.Name, &c.Comment, &c.CreatedAt); err != nil {
				log.Printf("scanning row: %v", err)
				http.Error(w, "Internal server error", 500)
				return
			}
			comments = append(comments, c)
		}
		if err := tmpl.Execute(w, comments); err != nil {
			log.Printf("exec template: %v", err)
			http.Error(w, "Internal server error", 500)
			return
		}
	})
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const guestbookTmpl = `<!doctype html>
<title>Guestbook</title>
<h1>Guestbook</h1>
<form method="POST">
<div><label>Name <input name=name></label></div>
<div><label>Comment <textarea name=comment></textarea></div>
<button>Submit</button>
</form>
<hr>
<h2>Comments</h2>
<ul>
{{ range . }}
<li>{{.Name}} on {{.CreatedAt}}: {{.Comment}}</li>
{{ end }}
</ul>
`
