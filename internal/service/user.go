package service

import (
	"fmt"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
)

type UserRepo interface {
	BuyItem(userID, itemID int) error
	SendCoin(fromUserName, toUserName string, amount int) error
	GetInfo(userID int) (*entities.InfoResponse, error)
	Auth(userName, password string) (*entities.User, error)
	GetUserID(userName string) (int, error)
	GetItemID(itemName string) (int, error)
}

type UserService struct {
	UserRepo UserRepo
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
	if err := u.UserRepo.SendCoin(fromUser, toUser, amount); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetInfo(userName string) (*entities.InfoResponse, error) {
	userID, err := u.UserRepo.GetUserID(userName)
	if err != nil {
		return nil, err
	}
	coinHistory, err := u.UserRepo.GetInfo(userID)
	if err != nil {
		return nil, err
	}

	return coinHistory, nil
}

func (u *UserService) Auth(userName, password string) (*entities.User, error) {
	user, err := u.UserRepo.Auth(userName, password)
	if err != nil {
		fmt.Println("auth service error", err)
		return nil, err
	}

	fmt.Println("auth service no errir",user)

	return user, nil
}
