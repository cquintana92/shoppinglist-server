package api

import (
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
)

func deleteAllChecked(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received DELETE ALL CHECKED request")
	err := storage.DeleteAllChecked()
	if err != nil {
		log.Logger.Errorf("Could not delete all checked: %+v", err)
		http.Error(w, "Could not delete all checked", 500)
	} else {
		respondAll(w, r)
	}
}
