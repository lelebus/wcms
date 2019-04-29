package purchase

type Purchase struct {
	Wine     int    `json:"id"`
	Date     string `json:"date"`
	Supplier string `json:"supplier"`
	Quantity int    `json:"quantity"`
	Cost     string `json:"cost"`
}
