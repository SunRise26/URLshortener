package main

import (
	"log"
	"net/http"
	"time"
)

type server struct {
	addr        	string
	db          	*database
	randLinkLen 	int
	dataLifetime	int64
}

var serv server

func main() {
	initConf()
	serv = server {
		conf.GetString("service.host") + ":" + conf.GetString("service.port"),
		OpenDatabase(),
		conf.GetInt("service.randLinkLen"),
		conf.GetInt64("service.dataLifetime"),
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/create", createLink)
	http.HandleFunc("/showLinks", showLinks)
	http.HandleFunc("/showRequests", showRequests)

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			serv.db.deleteOldRequests(serv.dataLifetime)
		}
	}()
	log.Fatal(http.ListenAndServe(serv.addr, nil))
}
