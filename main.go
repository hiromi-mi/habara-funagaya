package main

import (
	//"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// Event All Event Lists
type Event struct {
	Title   string
	Members map[string]string
}

var events = make(map[string]Event)

var templates = template.Must(template.ParseFiles("vote.html"))
var validPath = regexp.MustCompile("^/(events|new|register|unregister)/([a-zA-Z0-9]+)$")

func eventshandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	event, ok := events[m[2]]
	if !ok {
		http.Redirect(w, r, "/new/"+m[2], http.StatusFound)
		return
	}
	err := templates.ExecuteTemplate(w, "vote.html", event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func registershandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	event, ok := events[m[2]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	event.Members["hiromi_mi"] = "aaa"
	events[m[2]] = event
	http.Redirect(w, r, "/events/"+m[2], http.StatusFound)
}
func unregistershandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	event, ok := events[m[2]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	delete(event.Members, "hiromi_mi")
	events[m[2]] = event
	http.Redirect(w, r, "/events/"+m[2], http.StatusFound)
}

func main() {
	p := Event{Title: "TestEvent", Members: make(map[string]string)}
	events[p.Title] = p
	http.HandleFunc("/events/", eventshandler)
	http.HandleFunc("/register/", registershandler)
	http.HandleFunc("/unregister/", unregistershandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
