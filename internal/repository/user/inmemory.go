package repository

import (
	"fmt"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
)

type InmemoryRepo struct {
	items map[string]int
	users []*entities.User
}

func NewInmemoryRepo() *InmemoryRepo {
	return &InmemoryRepo{
		items: map[string]int{
			"t-shirt": 80,
			"cup": 20,
			"book": 50,
			"pen": 10,
		},
		users: []*entities.User{
			{
				ID: 1,
				Username: "Konstantin",
				Password: "12345",
				Coins: 100,
				Inventory: make([]*entities.Item, 0),
				CoinHistory: entities.CoinHistory{
					Received: make([]*entities.ReceiveOperation, 0),
					Sent: make([]*entities.SentOperation, 0),
				},
			},
			{
				ID: 2,
				Username: "Masha",
				Password: "12345",
				Coins: 100,
				Inventory: make([]*entities.Item, 0),
				CoinHistory: entities.CoinHistory{
					Received: make([]*entities.ReceiveOperation, 0),
					Sent: make([]*entities.SentOperation, 0),
				},
			},
			{
				ID: 3,
				Username: "Jack",
				Password: "12345",
				Coins: 100,
				Inventory: make([]*entities.Item, 0),
				CoinHistory: entities.CoinHistory{
					Received: make([]*entities.ReceiveOperation, 0),
					Sent: make([]*entities.SentOperation, 0),
				},
			},
		},
	}
}

func (i *InmemoryRepo) GetUserByName(userName string) (*entities.User, error) {
	for _, user := range i.users {
		if user.Username == userName {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (i *InmemoryRepo) AddToInventory(user *entities.User, itemName string) {
	for _, item := range user.Inventory {
		if item.ItemType == itemName {
			item.Quantity++
			return
		}
	}
	user.Inventory = append(user.Inventory, &entities.Item{
		ItemType: itemName,
		Quantity: 1,
	})
} 

func (i *InmemoryRepo) BuyItem(userName string, itemName string) error {
	price := i.items[itemName]
	var user *entities.User
	for ind := range i.users {
		if i.users[ind].Username == userName {
			user = i.users[ind]
			break
		}
	}
	balance := user.Coins
	if balance < price {
		return fmt.Errorf("not enough balance")
	}

	user.Coins -= price

	i.AddToInventory(user, itemName)

	return nil
}


func (i *InmemoryRepo) AddReceiveRecord(user *entities.User, fromUser string, amount int) {
	user.CoinHistory.Received = append(user.CoinHistory.Received, &entities.ReceiveOperation{
		FromUser: fromUser,
		Amount: amount,
	})
}

func (i *InmemoryRepo) AddSentRecord(user *entities.User, toUser string, amount int) {
	user.CoinHistory.Sent = append(user.CoinHistory.Sent, &entities.SentOperation{
		ToUser: toUser,
		Amount: amount,
	})
}

func (i *InmemoryRepo) SendCoin(fromUserName, toUserName string, amount int) error {
	fromUser, err := i.GetUserByName(fromUserName)
	if err != nil {
		return err
	}

	toUser, err := i.GetUserByName(toUserName)
	if err != nil {
		return err
	}

	if fromUser.Coins < amount {
		return fmt.Errorf("not enough balance")
	}

	fromUser.Coins -= amount
	toUser.Coins += amount


	i.AddReceiveRecord(toUser, fromUser.Username, amount)
	i.AddSentRecord(fromUser, toUser.Username, amount)


	return nil
}

func (i *InmemoryRepo) GetInfo(userName string) (*entities.InfoResponse, error) {
	user, err := i.GetUserByName(userName)
	if err != nil {
		return nil, err
	}

	return &entities.InfoResponse{
		Coins: user.Coins,
		Inventory: user.Inventory,
		CoinHistory: user.CoinHistory,
	}, nil
}

func (i *InmemoryRepo) Auth(userName, password string) (*entities.User, error) {
	return nil, nil
}