package api

import (
	"encoding/json"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"
	"time"
)

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
