package lib

import (
	"fmt"
	"html/template"
	"net/http"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["email"])
	fmt.Println(r.Method)

	data := authData{
		IsAuthed: false,
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Println("Not Authenticated")
		data.IsAuthed = false
	} else {
		fmt.Println("Authenticated")
		data.IsAuthed = true
	}

	t, _ := template.ParseFiles("./static/contact.html")
	t.Execute(w, data)
}
