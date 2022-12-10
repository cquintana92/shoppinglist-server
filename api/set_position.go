package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"strconv"
)

type setPosition struct {
	Position int `json:"position"`
}

func updatePosition(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received UPDATE POSITION request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("Could not read the request body: %+v", err)
		http.Error(w, "Could not update the record", 500)
		return
	}

	req := setPosition{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Logger.Errorf("Could not convert to JSON: %+v", err)
		http.Error(w, "Could not update the record position", 400)
		return
	}
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Logger.Errorf("Could not read the id as int: %+v", err)
		http.Error(w, "Could not update the record position", 400)
		return
	}
	log.Logger.Debugf("Received id=%d updates=%+v", id, req)
	err = storage.MoveToPosition(int(id), req.Position)
	if err != nil {
		log.Logger.Errorf("Error updating record position: %+v", err)
		http.Error(w, "Could not update the record position", 400)
	} else {
		respondAll(w, r)
	}
}
