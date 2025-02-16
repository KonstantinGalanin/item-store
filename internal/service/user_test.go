package service

import (
	"errors"
	"testing"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	repository "github.com/KonstantinGalanin/itemStore/internal/repository/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepo(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userName := "test_user"
		itemName := "test_item"
		userID := 1
		itemID := 100

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetItemID(itemName).Return(itemID, nil)
		mockRepo.EXPECT().BuyItem(userID, itemID).Return(nil)

		err := userService.BuyItem(userName, itemName)
		assert.NoError(t, err)
	})

	t.Run("get user id error", func(t *testing.T) {
		userName := "test_user"
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(0, someError)

		err := userService.BuyItem(userName, "test_item")
		assert.Error(t, err)
		assert.Equal(t, someError, err)
	})

	t.Run("get item id error", func(t *testing.T) {
		userName := "test_user"
		itemName := "test_item"
		userID := 1
		someError := errors.New("item not found")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetItemID(itemName).Return(0, someError)

		err := userService.BuyItem(userName, itemName)
		assert.Error(t, err)
		assert.Equal(t, someError, err)
	})

	t.Run("buy item error", func(t *testing.T) {
		userName := "test_user"
		itemName := "test_item"
		userID := 1
		itemID := 100
		someError := errors.New("failed to buy item")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetItemID(itemName).Return(itemID, nil)
		mockRepo.EXPECT().BuyItem(userID, itemID).Return(someError)

		err := userService.BuyItem(userName, itemName)
		assert.Error(t, err)
		assert.Equal(t, someError, err)
	})
}

func TestSendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepo(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		fromUser := "alice"
		toUser := "bob"
		amount := 50
		fromUserID := 1
		toUserID := 2

		mockRepo.EXPECT().GetUserID(fromUser).Return(fromUserID, nil)
		mockRepo.EXPECT().GetUserID(toUser).Return(toUserID, nil)
		mockRepo.EXPECT().SendCoin(fromUserID, toUserID, amount).Return(nil)

		err := userService.SendCoin(fromUser, toUser, amount)
		assert.NoError(t, err)
	})

	t.Run("get user id error", func(t *testing.T) {
		fromUser := "alice"
		toUser := "bob"
		amount := 50
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(fromUser).Return(0, someError)

		err := userService.SendCoin(fromUser, toUser, amount)
		assert.Error(t, err)
		assert.Equal(t, someError, err)
	})

	t.Run("send coin id error", func(t *testing.T) {
		fromUser := "alice"
		toUser := "bob"
		amount := 50
		fromUserID := 1
		toUserID := 2
		someError := errors.New("failed to send coins")

		mockRepo.EXPECT().GetUserID(fromUser).Return(fromUserID, nil)
		mockRepo.EXPECT().GetUserID(toUser).Return(toUserID, nil)
		mockRepo.EXPECT().SendCoin(fromUserID, toUserID, amount).Return(someError)

		err := userService.SendCoin(fromUser, toUser, amount)
		assert.Error(t, err)
		assert.Equal(t, someError, err)
	})

}

func TestGetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepo(ctrl)
	userService := NewUserService(mockRepo)
	t.Run("success", func(t *testing.T) {
		userName := "test_user"
		userID := 1
		coins := 100
		inventory := []*entities.Item{{ItemType: "cup", Quantity: 2}}
		receives := []*entities.ReceiveOperation{{FromUser: "Konstantin", Amount: 50}}
		sents := []*entities.SentOperation{{ToUser: "Masha", Amount: 20}}

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetCoinsInfo(userID).Return(coins, nil)
		mockRepo.EXPECT().GetInventoryInfo(userID).Return(inventory, nil)
		mockRepo.EXPECT().GetReceiveInfo(userID).Return(receives, nil)
		mockRepo.EXPECT().GetSentInfo(userID).Return(sents, nil)

		info, err := userService.GetInfo(userName)
		assert.NoError(t, err)
		assert.Equal(t, coins, info.Coins)
		assert.Equal(t, inventory, info.Inventory)
		assert.Equal(t, receives, info.CoinHistory.Received)
		assert.Equal(t, sents, info.CoinHistory.Sent)
	})

	t.Run("get user id error", func(t *testing.T) {
		userName := "test_user"
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(0, someError)

		info, err := userService.GetInfo(userName)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, someError, err)
	})

	t.Run("get coins info error", func(t *testing.T) {
		userName := "test_user"
		userID := 1
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetCoinsInfo(userID).Return(0, someError)

		info, err := userService.GetInfo(userName)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, someError, err)
	})

	t.Run("get inventory info error", func(t *testing.T) {
		userName := "test_user"
		userID := 1
		coins := 100
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetCoinsInfo(userID).Return(coins, nil)
		mockRepo.EXPECT().GetInventoryInfo(userID).Return(nil, someError)

		info, err := userService.GetInfo(userName)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, someError, err)
	})

	t.Run("get receive info error", func(t *testing.T) {
		userName := "test_user"
		userID := 1
		coins := 100
		inventory := []*entities.Item{{ItemType: "cup", Quantity: 2}}
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetCoinsInfo(userID).Return(coins, nil)
		mockRepo.EXPECT().GetInventoryInfo(userID).Return(inventory, nil)
		mockRepo.EXPECT().GetReceiveInfo(userID).Return(nil, someError)

		info, err := userService.GetInfo(userName)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, someError, err)
	})

	t.Run("get sent info error", func(t *testing.T) {
		userName := "test_user"
		userID := 1
		coins := 100
		inventory := []*entities.Item{{ItemType: "cup", Quantity: 2}}
		receives := []*entities.ReceiveOperation{{FromUser: "Konstantin", Amount: 50}}
		someError := errors.New("user not found")

		mockRepo.EXPECT().GetUserID(userName).Return(userID, nil)
		mockRepo.EXPECT().GetCoinsInfo(userID).Return(coins, nil)
		mockRepo.EXPECT().GetInventoryInfo(userID).Return(inventory, nil)
		mockRepo.EXPECT().GetReceiveInfo(userID).Return(receives, nil)
		mockRepo.EXPECT().GetSentInfo(userID).Return(nil, someError)

		info, err := userService.GetInfo(userName)
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, someError, err)
	})

}

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepo(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("success", func(t *testing.T) {
		userName := "test_user"
		password := "secure_password"
		expectedUser := &entities.User{ID: 1, Username: userName}

		mockRepo.EXPECT().Auth(userName, password).Return(expectedUser, nil)

		user, err := userService.Auth(userName, password)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("error", func(t *testing.T) {
		userName := "test_user"
		password := "wrong_password"
		someError := errors.New("authentication failed")

		mockRepo.EXPECT().Auth(userName, password).Return(nil, someError)

		user, err := userService.Auth(userName, password)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, someError, err)
	})
}
