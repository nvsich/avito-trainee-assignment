package model

type EmployeeInfo struct {
	Coins       int
	Inventory   []InventoryItem
	CoinHistory CoinHistory
}

type InventoryItem struct {
	Type     string
	Quantity int
}

type CoinHistory struct {
	Received []CoinTransaction
	Sent     []CoinTransaction
}

type CoinTransaction struct {
	User   string
	Amount int
}
