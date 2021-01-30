package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/josephbmanley/OpenSkins-Common/datastore"
	"github.com/josephbmanley/OpenSkins-Common/datatypes"
	"github.com/josephbmanley/OpenSkins/pluginmanager"
	"io"
	"log"
	"net/http"
)

var activeSkinstore datastore.Skinstore
var activeUserstore datastore.Userstore

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is healthy!")
}

func getCharacter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	return

	vars := mux.Vars(r)
	userID := vars["user"]
	charID := vars["char"]

	user, err := activeUserstore.GetUser(userID)

	// Check if there was an error getting the user
	if err != nil {
		log.Printf("An error occured when getting a character: %v", err.Error())
		w.WriteHeader(500)
		return
	}

	// Check if user exists
	if user == nil {
		w.WriteHeader(404)
		return
	}

	// Find character
	for _, character := range user.Characters {
		if character.UID == charID {
			json.NewEncoder(w).Encode(character) // Return character Json
			return                               // Return 200
		}
	}

	// Return 404
	w.WriteHeader(404)
	return
}

func createCharacter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	return

	vars := mux.Vars(r)
	userID := vars["user"]
	charID := vars["char"]

	user, err := activeUserstore.GetUser(userID)

	// Check if there was an error getting the user
	if err != nil {
		log.Printf("CREATE CHARACTER - Failed to get user: %v", err.Error())
		w.WriteHeader(500)
		return
	}

	// Check if user exists
	if user == nil {
		w.WriteHeader(400)
		w.Write([]byte("User does not exist!"))
		return
	}

	// Check if character already exists
	for _, character := range user.Characters {
		if character.UID == charID {
			w.WriteHeader(400)
			w.Write([]byte("Character already exists!"))
			return
		}
	}

	newCharacter := datatypes.Character{
		UID: charID,
		Skin: datatypes.Skin{
			UID: "none",
		},
	}

	// Add character to user object
	user.Characters = append(user.Characters, newCharacter)
	err = activeUserstore.SetUser(user)

	if err != nil {
		log.Printf("CREATE CHARACTER - Failed to set user %v", err.Error())
		w.WriteHeader(500)
	}

	return
}

func getSkin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skinID := vars["skin"]

	// Get skin
	skin, err := activeSkinstore.GetSkin(skinID)
	if err != nil {
		log.Printf("GET SKIN - Failed to get skin %v", err.Error())
		w.WriteHeader(500)
	}

	// Check if skin exists
	if skin == nil {
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(skin) // Return skin Json
	return                          // Return 200
}

func createSkin(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	skinID := vars["skin"]

	// Lookup skin object
	skin, err := activeSkinstore.GetSkin(skinID)
	if err != nil {
		log.Printf("CREATE SKIN - Failed to get skin: %v", err.Error())
		w.WriteHeader(500)
		return
	}

	// Check if skin exists
	if skin != nil {
		w.WriteHeader(400)
		w.Write([]byte("Skin already exists!"))
		return
	}

	// Load skin data from request
	skinBytes, err := readfile(r, "skin")
	if err != nil {
		log.Printf("CREATE SKIN - Failed to read upload: %v", err.Error())
	}

	// Verify data exists
	if skinBytes == nil {
		w.WriteHeader(400)
		w.Write([]byte("Missing skin data!"))
		return
	}

	// Create new skin object
	err = activeSkinstore.AddSkin(skinID, skinBytes)

	// Validate skin creation
	if err != nil {
		log.Printf("CREATE SKIN - Failed to add skin: %v", err.Error())
		w.WriteHeader(500)
		return
	}

	return // Return 200
}

// Helper function to read a file from a request form
func readfile(r *http.Request, fileName string) ([]byte, error) {

	err := r.ParseMultipartForm(32 << 20) // Limit upload size
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	file, _, err := r.FormFile(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Copy the file data to my buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	contents := buf.Bytes() // Get data
	buf.Reset()             // Clear buffer

	return contents, nil
}

// StartWebserver starts the webserver
func StartWebserver(port int) {

	activeSkinstore = pluginmanager.LoadedSkinstores[0] // Load first skinstore

	// Intialize router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Health check path
	myRouter.HandleFunc("/health", healthCheck).Methods("GET")

	// Character routes
	myRouter.HandleFunc("/get/character/{user}/{char}", getCharacter).Methods("GET")
	myRouter.HandleFunc("/create/character/{user}/{char}", createCharacter).Methods("POST")

	// Skin routes
	myRouter.HandleFunc("/get/skin/{skin}", getSkin).Methods("GET")
	myRouter.HandleFunc("/create/skin/{skin}", createSkin).Methods("POST")

	// Run webserver
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), myRouter))
}
