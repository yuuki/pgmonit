package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"

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

type Replication struct {
	Pid             int       `db:"pid"`
	UseSysid        string    `db:"usesysid"`
	UseName         string    `db:"usename"`
	ApplicationName string    `db:"application_name"`
	ClientAddr      string    `db:"client_addr"`
	ClientHostname  sql.NullString `db:"client_hostname"`
	ClientPort      string    `db:"client_port"`
	BackendStart    time.Time `db:"backend_start"`
	State           string    `db:"state"`
	SentLocation    string    `db:"sent_location"`
	WriteLocation   string    `db:"write_location"`
	FlushLocation   string    `db:"flush_location"`
	ReplayLocation  string    `db:"replay_location"`
	SyncPriority    int       `db:"sync_priority"`
	SyncState       string    `db:"sync_state"`
}

func index(c web.C, w http.ResponseWriter, r *http.Request) {
	clientCounts := []ClientCount{}
	err := DB.Select(&clientCounts, `
		SELECT client_addr, count(*) AS count FROM pg_stat_activity
		GROUP BY client_addr ORDER BY count
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, row := range clientCounts {
		fmt.Fprintf(w, "%s: %d\n", row.Addr, row.Count)
	}

	replications := []Replication{}
	err = DB.Select(&replications, "SELECT * FROM pg_stat_replication")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, row := range replications {
		var name string
		if row.ClientHostname.Valid {
			name = row.ClientHostname.String
		}

		fmt.Fprintf(w, "%s: %s %s %s %s %s %s %s %s %d %s\n",
			row.ApplicationName, row.ClientAddr, name,
			row.BackendStart, row.State, row.SentLocation,
			row.WriteLocation, row.FlushLocation, row.ReplayLocation,
			row.SyncPriority, row.SyncState,
		)
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
