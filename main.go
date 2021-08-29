package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Example object used to demonstrate backend data.
type object struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Value string `json:"value"`
}

var listOfObjects []object

func homePage(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		log.Fatal("Error with Homepage")
	}
	log.Info("Homepage Hit")
}

// Handle GET POST PUT PATCH DELETE methods
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/objects", getAllObjects).Methods(http.MethodGet)
	myRouter.HandleFunc("/object", postObject).Methods(http.MethodPost)
	myRouter.HandleFunc("/object/{id}", getObject).Methods(http.MethodGet)
	myRouter.HandleFunc("/object/{id}", deleteObject).Methods(http.MethodDelete)
	myRouter.HandleFunc("/object/{id}", putObject).Methods(http.MethodPut)
	myRouter.HandleFunc("/object/{id}", patchObject).Methods(http.MethodPatch)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	// Create some test objects
	listOfObjects = []object{
		{
			Name:  "alpha",
			ID:    3343,
			Value: "PTR452s",
		},
		{
			Name:  "beta",
			ID:    8374,
			Value: "LSD532j",
		},
		{
			Name:  "gamma",
			ID:    1201,
			Value: "WLD293i",
		},
	}

	handleRequests()
}

// Return all objects
func getAllObjects(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(listOfObjects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error encoding JSON")
	}
}

// Return a specific object with a given ID in the URI
func getObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])
	found := false

	for index := range listOfObjects {
		if listOfObjects[index].ID == key {
			found = true
			err := json.NewEncoder(w).Encode(listOfObjects[index])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal("Error encoding JSON")
			}
			break
		}
	}

	if found == false {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Could not find object")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}
	}
}

// POST an object
func postObject(w http.ResponseWriter, r *http.Request) {

	requestBody, _ := ioutil.ReadAll(r.Body)
	var newObject object

	err := json.Unmarshal(requestBody, &newObject)
	if err != nil {
		log.Error("Unable to decode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode("Could not unmarshal JSON, invalid request")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}

	} else {

		listOfObjects = append(listOfObjects, newObject)

		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(newObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}
	}
}

// DELETE a specific object with a given ID
func deleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])
	found := false
	location := 0

	for index := range listOfObjects {
		if listOfObjects[index].ID == key {
			found = true
			location = index
			break
		}
	}

	if found == true {
		listOfObjects = append(listOfObjects[:location], listOfObjects[location+1:]...)
		err := json.NewEncoder(w).Encode("Removed")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Could not find object")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}
	}
}

// PUT an object
func putObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])
	found := false

	requestBody, _ := ioutil.ReadAll(r.Body)
	var newObject object

	err := json.Unmarshal(requestBody, &newObject)
	if err != nil {
		log.Error("Unable to decode JSON")

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode("Could not unmarshal JSON, invalid request")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}

	} else {

		for index := range listOfObjects {
			if listOfObjects[index].ID == key {
				found = true
				listOfObjects[index] = newObject
				err := json.NewEncoder(w).Encode(newObject)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Fatal("Error encoding JSON")
				}
				break
			}
		}

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode("Could not find object")
			if err != nil {
				log.Fatal("Error encoding JSON")
			}
		}
	}
}

// PATCH an object
func patchObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])
	found := false

	requestBody, _ := ioutil.ReadAll(r.Body)
	var newObject object

	err := json.Unmarshal(requestBody, &newObject)
	if err != nil {
		log.Error("Unable to decode JSON")

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode("Could not unmarshal JSON, invalid request")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error encoding JSON")
		}

	} else {

		for index := range listOfObjects {
			if listOfObjects[index].ID == key {
				found = true
				listOfObjects[index].merge(newObject)
				err := json.NewEncoder(w).Encode(listOfObjects[index])
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Fatal("Error encoding JSON")
				}
				break
			}
		}

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode("Could not find object")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal("Error encoding JSON")
			}
		}
	}
}

// When patching an object, assess which values are non-nil to move over
func (oldStruct *object) merge(newStruct object) {
	if newStruct.ID != 0 {
		oldStruct.ID = newStruct.ID
	}

	if newStruct.Name != "" {
		oldStruct.Name = newStruct.Name
	}

	if newStruct.Value != "" {
		oldStruct.Value = newStruct.Value
	}
}
