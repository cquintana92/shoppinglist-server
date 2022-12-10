package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"strconv"
)

func toggleChecked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Logger.Errorf("Error parsing id to int: %+v", err)
		http.Error(w, "Bad id", 400)
		return
	}
	log.Logger.Infof("Received TOGGLE CHECKED [id=%d] request", id)
	err = storage.ToggleChecked(int(id))
	if err != nil {
		log.Logger.Errorf("Could not toggle checked id=%d: %+v", id, err)
		http.Error(w, "Could not toggle checked", 400)
	} else {
		respondAll(w, r)
	}
}
