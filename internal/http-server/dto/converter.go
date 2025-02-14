package dto

import (
	resp "avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/model"
)

func ToInfoResponse(employeeInfo model.EmployeeInfo) resp.InfoResponse {
	return resp.InfoResponse{
		Coins: employeeInfo.Coins,
		CoinHistory: resp.CoinHistory{
			Received: convertTransactions(employeeInfo.CoinHistory.Received),
			Sent:     convertTransactions(employeeInfo.CoinHistory.Sent),
		},
		Inventory: convertInventory(employeeInfo.Inventory),
	}
}

func convertTransactions(transactions []model.CoinTransaction) []resp.CoinTransaction {
	converted := make([]resp.CoinTransaction, len(transactions))
	for i := range transactions {
		converted[i] = resp.CoinTransaction{
			User:   transactions[i].User,
			Amount: transactions[i].Amount,
		}
	}
	return converted
}

func convertInventory(inventory []model.InventoryItem) []resp.InventoryItem {
	converted := make([]resp.InventoryItem, len(inventory))
	for i := range inventory {
		converted[i] = resp.InventoryItem{
			Type:     inventory[i].Type,
			Quantity: inventory[i].Quantity,
		}
	}
	return converted
}
