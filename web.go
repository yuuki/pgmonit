package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
)

type ClientCount struct {
	Addr  string `db:"client_addr"`
	Count int    `db:"count"`
}

func index(c web.C, w http.ResponseWriter, r *http.Request) {
	var clientCounts []ClientCount
	err := DB.Select(&clientCounts, `
		SELECT client_addr, count(*) AS count FROM pg_stat_activity
		GROUP BY client_addr ORDER BY count
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, count := range clientCounts {
		fmt.Fprintf(w, "%s: %d\n", count.Addr, count.Count)
	}
}

func AppUp() {
	goji.Use(middleware.Recoverer)
	goji.Use(middleware.NoCache)

	goji.Get("/", index)

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
