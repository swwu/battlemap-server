package server

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/swwu/battlemap-server/logging"
	"github.com/swwu/battlemap-server/ruleset"
)

func Serve(gamespaces map[string]Gamespace,
	rulesets map[string]ruleset.Ruleset) {
	logging.Info.Println("Serving on :10010")

	checkGamespace := func(gamespace_id string, w http.ResponseWriter, cb func()) {
		if _, exists := gamespaces[gamespace_id]; exists {
			cb()
		} else {
			fmt.Fprintf(w, "Gamespace: %q doesn't exist", html.EscapeString(gamespace_id))
		}
	}

	router := mux.NewRouter()

	// Cross domain headers
	corsHandler := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
					w.Header().Set("Access-Control-Allow-Headers", acrh[0])
				}
				//w.Header().Set("Access-Control-Allow-Credentials", "True")
				if _, ok := r.Header["Access-Control-Allow-Origin"]; ok {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					if _, oko := r.Header["Origin"]; oko {
						w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
					} else {
						w.Header().Set("Access-Control-Allow-Origin", "*")
					}
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			} else {
				fn(w, r)
			}
		}
	}

	router.HandleFunc("/",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
		})).Methods("POST", "OPTIONS")

	router.HandleFunc("/gamespace/",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
		})).Methods("POST", "OPTIONS")

	router.HandleFunc("/gamespace/{gamespace_id}",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["gamespace_id"]
			checkGamespace(id, w, func() {
				fmt.Fprintf(w, "Gamespace: %q exists", html.EscapeString(id))
			})
		})).Methods("GET", "OPTIONS")

	router.HandleFunc("/gamespace/{gamespace_id}",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			// TODO: update baseValues
		})).Methods("POST", "PUT", "OPTIONS")

	router.HandleFunc("/gamespace/{gamespace_id}/entity/{entity_id}",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			gid := mux.Vars(r)["gamespace_id"]
			eid := mux.Vars(r)["entity_id"]
			checkGamespace(gid, w, func() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				bodyText, putErr := ioutil.ReadAll(r.Body)
				if putErr != nil {
					fmt.Fprintf(w, "some kind of put error")
				} else {
					gamespaces[gid].Entity(eid).JsonPut(bodyText)

					jsonDump, dumpErr := gamespaces[gid].Entity(eid).JsonDump()
					if dumpErr != nil {
						fmt.Fprintf(w, "some kind of dump error")
					} else {
						fmt.Fprintf(w, string(jsonDump))
					}
				}
			})
		})).Methods("POST", "PUT", "OPTIONS")

	router.HandleFunc("/gamespace/{gamespace_id}/entity/{entity_id}",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			gid := mux.Vars(r)["gamespace_id"]
			eid := mux.Vars(r)["entity_id"]
			checkGamespace(gid, w, func() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				jsonDump, err := gamespaces[gid].Entity(eid).JsonDump()
				if err != nil {
					fmt.Fprintf(w, "some kind of error")
				} else {
					fmt.Fprintf(w, string(jsonDump))
				}
				//w.Write(jsonDump)
			})
		})).Methods("GET", "OPTIONS")

	http.Handle("/", router)

	http.ListenAndServe(":10010", nil)
}
