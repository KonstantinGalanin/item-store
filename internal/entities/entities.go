package entities

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string
	Coins       int
	Inventory   []*Item
	CoinHistory CoinHistory
}

type Item struct {
	ItemType string `json:"type"`
	Quantity int    `json:"quantity"`
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []*Item     `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type CoinHistory struct {
	Received []*ReceiveOperation `json:"received"`
	Sent     []*SentOperation    `json:"sent"`
}

type SentOperation struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type ReceiveOperation struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}
