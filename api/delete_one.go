package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"strconv"
)

func deleteOne(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received DELETE ONE request")
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Logger.Errorf("Could not read the id as int: %+v", err)
		http.Error(w, "Could not delete the record", 400)
		return
	}
	log.Logger.Debugf("Received id=%d", id)
	err = storage.DeleteOne(int(id))
	if err != nil {
		log.Logger.Errorf("Error deleting record: %+v", err)
		http.Error(w, "Could not delete the record", 400)
	} else {
		respondAll(w, r)
	}
}
