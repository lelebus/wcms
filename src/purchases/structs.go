package purchase

type Purchase struct {
	ID       int    `json:"id"`
	Wine     int    `json:"wine"`
	Date     string `json:"date"`
	Supplier string `json:"supplier"`
	Quantity int    `json:"quantity"`
	Cost     string `json:"cost"`
}
