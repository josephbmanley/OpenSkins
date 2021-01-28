package runtime

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/josephbmanley/OpenSkins-Common/datastore"
	"log"
	"net/http"
)

var activeSkinstore datastore.Skinstore
var activeUserstore datastore.Userstore

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is healthy!")
}

func getCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user"]
	charID := vars["char"]

	user, err := activeUserstore.GetUser(userID)

	// Check if there was an error getting the user
	if err != nil {
		log.Panicf("An error occured when getting a character: %v", err.Error())
		w.WriteHeader(500)
		return
	}

	// Check if user exists
	if user == nil {
		w.WriteHeader(404)
		return
	}

	for _, character := range user.Characters {
		if character.UID == charID {
			json.NewEncoder(w).Encode(character)
			return
		}
	}

	w.WriteHeader(404)
	return
}

// StartWebserver starts the webserver
func StartWebserver(port int) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/health", healthCheck)
	myRouter.HandleFunc("/get/character/{user}/{char}", getCharacter)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), myRouter))
}
