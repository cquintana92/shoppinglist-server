package storage

type ItemDB struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Checked   int    `json:"checked"`
	ListOrder int    `json:"listOrder"`
	CreatedAt string `json:"createdAt"`
}
