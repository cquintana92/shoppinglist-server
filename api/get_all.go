package api

import (
	"net/http"
	"shoppinglistserver/log"
)

func getAll(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received GET ALL request")
	respondAll(w, r)
}
