package server

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/ruleset"
)

func Serve(gamespaces map[string]Gamespace,
	rulesets map[string]ruleset.Ruleset) {
	logging.Info.Println("Serving on :10010")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	checkGamespace := func(gamespace_id string, w http.ResponseWriter, cb func()) {
		if _, exists := gamespaces[gamespace_id]; exists {
			cb()
		} else {
			fmt.Fprintf(w, "Gamespace: %q doesn't exist", html.EscapeString(gamespace_id))
		}
	}

	fmt.Println(upgrader)

	router := mux.NewRouter()

	router.HandleFunc("/gamespace/",
		func(w http.ResponseWriter, r *http.Request) {
		}).Methods("POST")

	router.HandleFunc("/gamespace/{gamespace_id}",
		func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["gamespace_id"]
			checkGamespace(id, w, func() {
				fmt.Fprintf(w, "Gamespace: %q exists", html.EscapeString(id))
			})
		}).Methods("GET")

	router.HandleFunc("/gamespace/{gamespace_id}",
		func(w http.ResponseWriter, r *http.Request) {
			// TODO: update baseValues
		}).Methods("POST", "PUT")

	router.HandleFunc("/gamespace/{gamespace_id}/entity",
		func(w http.ResponseWriter, r *http.Request) {
			gid := mux.Vars(r)["gamespace_id"]
			//eid := mux.Vars(r)["entity_id"]
			checkGamespace(gid, w, func() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			})
		}).Methods("POST", "PUT")

	router.HandleFunc("/gamespace/{gamespace_id}/entity/{entity_id}",
		func(w http.ResponseWriter, r *http.Request) {
			gid := mux.Vars(r)["gamespace_id"]
			eid := mux.Vars(r)["entity_id"]
			checkGamespace(gid, w, func() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				jsonDump, err := gamespaces[gid].Entity(eid).JsonDump()
				if err != nil {
					fmt.Fprintf(w, "some kind of error")
				} else {
					fmt.Fprintf(w, string(jsonDump))
				}
				//w.Write(jsonDump)
			})
		}).Methods("GET")

	router.HandleFunc("/{ruleName}",
		func(w http.ResponseWriter, r *http.Request) {
			ruleName := mux.Vars(r)["ruleName"]
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(ruleName))
		})
	http.Handle("/", router)

	http.ListenAndServe(":10010", nil)
}
