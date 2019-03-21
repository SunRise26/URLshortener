package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Available characters for random link creation
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// unavailable path's
var occupied = []string{
	"create",
	"showLinks",
	"showRequests",
	"testNotFound",
	"testFound",
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	if path == "" {
		showIndexPage(w)
	} else {
		origURI := serv.db.getOrigURL(path)

		if origURI != "" {
			http.Redirect(w, r, origURI, http.StatusFound)
			serv.db.addRequest(origURI, path)
		} else {
			serv.db.addRequest("Error 404", path)
			http.Error(w, "Unknown URL", http.StatusNotFound)
		}
	}
}

func createLink(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	origURL := strings.TrimSpace(r.FormValue("origURL"))
	shortURL := strings.TrimSpace(r.FormValue("shortURL"))

	if origURL == "" {
		http.Error(w, "Error: Empty request", http.StatusBadRequest)
		return
	} else if shortURL != "" && checkNewLink(shortURL) == false {
		http.Error(w, "Error: Shortened URL '" + shortURL + "' already exist", http.StatusBadRequest)
		return
	}
	if shortURL == "" {
		shortURL = randomLink()
	}
	serv.db.addLink(origURL, shortURL)
	fmt.Fprintf(w, "Shortened URL 'http://%s/%s' was created for\n'%s'",
		serv.addr, shortURL, origURL)
}

func randomLink() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, serv.randLinkLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	newLink := string(b)
	if checkNewLink(newLink) == false {
		return randomLink()
	}
	return newLink
}

func checkNewLink(link string) bool {
	for _, s := range occupied {
		if link == s {
			return false
		}
	}
	if serv.db.getOrigURL(link) != "" {
		return false
	}
	return true
}

func showLinks(w http.ResponseWriter, r *http.Request) {
	table := serv.db.getLinks()

	fmt.Fprintln(w, table)
}

func showRequests(w http.ResponseWriter, r *http.Request) {
	table := serv.db.getRequests()

	fmt.Fprintln(w, table)
}

func showIndexPage(w io.Writer) {
	type par struct {
		Addr string
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, par{serv.addr})
}
