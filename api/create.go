package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"
)

type newItem struct {
	Name string `json:"name"`
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received CREATE request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("Could not read the request body: %+v", err)
		http.Error(w, "Could not create a record", 500)
		return
	}

	req := newItem{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Logger.Errorf("Could not convert to JSON: %+v", err)
		http.Error(w, "Could not create a record", 500)
		return
	}
	if req.Name == "" {
		log.Logger.Errorf("Recieved a request to create with an empty name")
		http.Error(w, "Could not create a record", 400)
		return
	}
	req.Name = utils.SanitizeName(req.Name)
	log.Logger.Debugf("Received %+v", req)
	log.Logger.Infof("Creating item: %s", req.Name)
	err = storage.New(req.Name)
	if err != nil {
		if err == storage.ItemAlreadyExistsError {
			log.Logger.Errorf("Item [%s] already exists", req.Name)
			http.Error(w, "Item already exists", 409)
		} else {
			log.Logger.Errorf("Error creating a new record: %+v", err)
			http.Error(w, "Could not create a record", 500)
		}
	} else {
		respondAll(w, r)
	}
}
