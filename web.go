package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
)

func index(c web.C, w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func clients(c web.C, w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT client_addr, count(*) AS count FROM pg_stat_activity
		GROUP BY client_addr ORDER BY count
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var client_addr string
		var count int
		err := rows.Scan(&client_addr, &count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s: %d\n", client_addr, count)
	}
}

func AppUp() {
	goji.Use(middleware.Recoverer)
	goji.Use(middleware.NoCache)

	goji.Get("/", index)
	goji.Get("/clients", clients)

	goji.DefaultMux.Compile()
	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", goji.DefaultMux)

	listener := bind.Default()
	log.Println("Starting Goji on", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()
	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() { log.Printf("Goji stopped") })

	err := graceful.Serve(listener, http.DefaultServeMux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
