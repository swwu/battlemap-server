package server

import (
	"fmt"
	"html"
	"io/ioutil"
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

	// Cross domain headers
	corsHandler := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hallo")
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

	router.HandleFunc("/gamespace/{gamespace_id}/entity",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hallo postput")
			gid := mux.Vars(r)["gamespace_id"]
			eid := mux.Vars(r)["entity_id"]
			checkGamespace(gid, w, func() {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				bodyText, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Fprintf(w, "some kind of error")
				} else {
					gamespaces[gid].Entity(eid).JsonPut(bodyText)
				}
			})
		})).Methods("POST", "PUT", "OPTIONS")

	router.HandleFunc("/gamespace/{gamespace_id}/entity/{entity_id}",
		corsHandler(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hallo get")
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

	router.HandleFunc("/{ruleName}",
		func(w http.ResponseWriter, r *http.Request) {
			ruleName := mux.Vars(r)["ruleName"]
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(ruleName))
		})
	http.Handle("/", router)

	http.ListenAndServe(":10010", nil)
}
