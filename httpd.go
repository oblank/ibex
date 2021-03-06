package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"code.google.com/p/go-imap/go1/imap"
	"time"
)

var c *imap.Client

func requestLogger(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	handler.ServeHTTP(w, r)
    })
}

func inboxHandler(w http.ResponseWriter, r *http.Request) {
	selectMailbox(c, "INBOX", true)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprint(w, string(listRecent(c, 20)))
}

func allMailHandler(w http.ResponseWriter, r *http.Request) {
	selectMailbox(c, "[Gmail]/All Mail", true)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprint(w, string(listRecent(c, 20)))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := mux.Vars(r)["messageID"]
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprint(w, string(fetchMessage(c, messageID)))
}

func httpMain(debug bool) {
	c = initClient(debug)
	if (c == nil) {	return }
	selectMailbox(c, "INBOX", true)
	defer c.Logout(30 * time.Minute)

	r := mux.NewRouter()
	r.HandleFunc("/Inbox.json", inboxHandler)
	r.HandleFunc("/AllMail.json", allMailHandler)
	r.HandleFunc("/Messages/{messageID}", messageHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))
	http.Handle("/", r)
	http.ListenAndServe(":8080", requestLogger(http.DefaultServeMux))
}
