package main

import (
	//"errors"
	//"database/sql"
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

var events = make(map[string]*Event)

var templates = template.Must(template.ParseFiles("vote.html", "new.html", "index.html"))
var validPath = regexp.MustCompile("^/(events|new|register|unregister)/([a-zA-Z0-9]+)$")
var validTitle = regexp.MustCompile("^[a-zA-Z0-9]+$")

func metahandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		title := m[2]
		fn(w, r, title)
	}
}

func eventshandler(w http.ResponseWriter, r *http.Request, title string) {
	event, ok := events[title]
	if !ok {
		http.NotFound(w, r)
		//http.Redirect(w, r, "/new/", http.StatusFound)
		return
	}
	err := templates.ExecuteTemplate(w, "vote.html", event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func registershandler(w http.ResponseWriter, r *http.Request, title string) {
	id_raw := r.FormValue("id")
	id := template.HTMLEscapeString(id_raw)
	hitokoto_raw := r.FormValue("hitokoto")
	hitokoto := template.HTMLEscapeString(hitokoto_raw)
	event, ok := events[title]
	if !ok {
		http.NotFound(w, r)
		return
	}
	event.Members[id] = hitokoto
	events[title] = event
	http.Redirect(w, r, "/events/"+title, http.StatusFound)
}
func unregistershandler(w http.ResponseWriter, r *http.Request, title string) {
	event, ok := events[title]
	if !ok {
		http.NotFound(w, r)
		return
	}
	delete(event.Members, "hiromi_mi")
	events[title] = event
	http.Redirect(w, r, "/events/"+title, http.StatusFound)
}

func neweventhandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "new.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func indexhandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createeventhandler(w http.ResponseWriter, r *http.Request) {
	title_raw := r.FormValue("eventname") // get POST
	title := validTitle.FindString(title_raw)
	if title == "" {
		http.Error(w, "Event Name Not Found", http.StatusInternalServerError)
	}
	events[title] = &Event{Title: title, Members: make(map[string]string)}
	http.Redirect(w, r, "/events/"+title, http.StatusFound)
}

func main() {
	p := &Event{Title: "TestEvent", Members: make(map[string]string)}
	events[p.Title] = p
	http.HandleFunc("/index.html", indexhandler)
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/events/", metahandler(eventshandler))
	http.HandleFunc("/create/", createeventhandler)
	http.HandleFunc("/new/", neweventhandler)
	http.HandleFunc("/register/", metahandler(registershandler))
	http.HandleFunc("/unregister/", metahandler(unregistershandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
