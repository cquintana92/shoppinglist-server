package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"
	"strconv"
)

type setName struct {
	Name string `json:"name"`
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received UPDATE request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("Could not read the request body: %+v", err)
		http.Error(w, "Could not update the record", 500)
		return
	}

	req := setName{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Logger.Errorf("Could not convert to JSON: %+v", err)
		http.Error(w, "Could not update the record", 500)
		return
	}
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Logger.Errorf("Could not read the id as int: %+v", err)
		http.Error(w, "Could not update the record", 400)
		return
	}
	if req.Name == "" {
		log.Logger.Errorf("Recieved a request to update with an empty name")
		http.Error(w, "Could not create a record", 400)
		return
	}
	req.Name = utils.SanitizeName(req.Name)
	log.Logger.Debugf("Received id=%d updates=%+v", id, req)
	err = storage.Update(req.Name, int(id))
	if err != nil {
		log.Logger.Errorf("Error updating record: %+v", err)
		http.Error(w, "Could not update the record", 400)
	} else {
		respondAll(w, r)
	}
}
