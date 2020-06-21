package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"
	"strconv"
	"time"
)

type newItem struct {
	Name string `json:"name"`
}

type ResponseItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Checked   bool      `json:"checked"`
	ListOrder int       `json:"listOrder"`
	CreatedAt time.Time `json:"createdAt"`
}

func dbItemToResponseItem(item *storage.ItemDB) ResponseItem {
	createdAt, err := utils.DateFrom(item.CreatedAt)
	if err != nil {
		log.Logger.Errorf("Error parsing date: %+v", err)
		createdAt = time.Now()
	}
	return ResponseItem{
		Id:        item.Id,
		Name:      item.Name,
		Checked:   item.Checked == 1,
		ListOrder: item.ListOrder,
		CreatedAt: createdAt,
	}
}

func respondAll(w http.ResponseWriter, r *http.Request) {
	items, err := storage.GetAll()
	converted := make([]ResponseItem, len(items))
	for i, e := range items {
		converted[i] = dbItemToResponseItem(e)
	}
	if err != nil {
		log.Logger.Errorf("Error retrieving items: %+v", err)
		http.Error(w, "Could not retrieve items", 500)
	} else {
		log.Logger.Debug("Items retrieved")
		bytes, err := json.Marshal(converted)
		if err != nil {
			log.Logger.Errorf("Error marshalling to JSON: %+v", err)
			http.Error(w, "Could not retrieve items", 500)
		} else {
			w.WriteHeader(200)
			w.Write(bytes)
		}
	}
}

func getAll(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received GET ALL request")
	respondAll(w, r)
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
	req.Name = utils.CapitalizeName(req.Name)
	log.Logger.Debugf("Received %+v", req)
	log.Logger.Infof("Creating item: %s", req.Name)
	err = storage.New(req.Name)
	if err != nil {
		log.Logger.Errorf("Error creating a new record: %+v", err)
		http.Error(w, "Could not create a record", 500)
	} else {
		respondAll(w, r)
	}
}

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
	req.Name = utils.CapitalizeName(req.Name)
	log.Logger.Debugf("Received id=%d updates=%+v", id, req)
	err = storage.Update(req.Name, int(id))
	if err != nil {
		log.Logger.Errorf("Error updating record: %+v", err)
		http.Error(w, "Could not update the record", 400)
	} else {
		respondAll(w, r)
	}
}

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

func Run(port int) error {
	r := mux.NewRouter()
	r.HandleFunc("/", getAll).Methods("GET")
	r.HandleFunc("/", create).Methods("POST")
	r.HandleFunc("/", deleteAllChecked).Methods("DELETE")
	r.HandleFunc("/{id}", toggleChecked).Methods("PATCH")
	r.HandleFunc("/{id}", update).Methods("PUT")
	r.HandleFunc("/{id}", deleteOne).Methods("DELETE")
	r.HandleFunc("/{id}/position", updatePosition).Methods("PUT")

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Logger.Infof("[API] Started server at: %s", addr)
	return srv.ListenAndServe()
}
