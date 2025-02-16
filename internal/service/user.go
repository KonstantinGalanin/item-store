package service

import (
	"fmt"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
)

//go:generate mockgen -source=user.go -destination=../repository/user_repo_mock.go -package=repository
type UserRepo interface {
	BuyItem(userID, itemID int) error
	SendCoin(fromUserID, toUserID int, amount int) error
	Auth(userName, password string) (*entities.User, error)
	GetUserID(userName string) (int, error)
	GetItemID(itemName string) (int, error)
	GetCoinsInfo(userID int) (int, error)
	GetInventoryInfo(userID int) ([]*entities.Item, error)
	GetReceiveInfo(userID int) ([]*entities.ReceiveOperation, error)
	GetSentInfo(userID int) ([]*entities.SentOperation, error)
}

type UserService struct {
	UserRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService{
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) BuyItem(userName, itemName string) error {
	userID, err := u.UserRepo.GetUserID(userName)
	if err != nil {
		return err
	}
	itemID, err := u.UserRepo.GetItemID(itemName)
	if err != nil {
		return err
	}
	if err := u.UserRepo.BuyItem(userID, itemID); err != nil {
		return err
	}
	return nil
}

func (u *UserService) SendCoin(fromUser, toUser string, amount int) error {
	fromUserID, err := u.UserRepo.GetUserID(fromUser)
	if err != nil {
		return err
	}
	toUserID, err := u.UserRepo.GetUserID(toUser)
	if err := u.UserRepo.SendCoin(fromUserID, toUserID, amount); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetInfo(userName string) (*entities.InfoResponse, error) {
	userID, err := u.UserRepo.GetUserID(userName)
	if err != nil {
		return nil, err
	}

	coins, err := u.UserRepo.GetCoinsInfo(userID)
	if err != nil {
		return nil, err
	}

	inventory, err := u.UserRepo.GetInventoryInfo(userID)
	if err != nil {
		return nil, err
	}

	receives, err := u.UserRepo.GetReceiveInfo(userID)
	if err != nil {
		return nil, err
	}

	sents, err := u.UserRepo.GetSentInfo(userID)
	if err != nil {
		return nil, err
	}

	info := &entities.InfoResponse{
		Coins: coins,
		Inventory: inventory,
		CoinHistory: entities.CoinHistory{
			Received: receives,
			Sent: sents,
		},
	}

	return info, nil
}

func (u *UserService) Auth(userName, password string) (*entities.User, error) {
	user, err := u.UserRepo.Auth(userName, password)
	if err != nil {
		fmt.Println("auth service error", err)
		return nil, err
	}

	return user, nil
}
