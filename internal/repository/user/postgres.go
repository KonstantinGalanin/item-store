package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	"github.com/KonstantinGalanin/itemStore/internal/utils"
)

type UserPostgresRepo struct {
	DB *sql.DB
}

const (
	InitBalance = 1000
)

func NewUserPostgresRepo(db *sql.DB) *UserPostgresRepo {
	return &UserPostgresRepo{
		DB: db,
	}
}

func (u *UserPostgresRepo) GetItemID(itemName string) (int, error) {
	var itemID int
	row := u.DB.QueryRow(GetItemID, itemName)
	err := row.Scan(&itemID)
	if err != nil {
		return 0, fmt.Errorf("get item id error: %w", err)
	}

	return itemID, nil
}

func (u *UserPostgresRepo) GetUserID(username string) (int, error) {
	var userID int
	row := u.DB.QueryRow(GetUserID, username)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres get user id error: %w", utils.ErrNoUser)
		}
		return 0, fmt.Errorf("postgres get user id error: %w", err)
	}
	return userID, nil
}

func (u *UserPostgresRepo) BuyItem(userID, itemID int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int
	err = tx.QueryRow(GetBalance, userID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("get balance error: %w", err)
	}

	var price int
	err = tx.QueryRow(GetPrice, itemID).Scan(&price)
	if err != nil {
		return fmt.Errorf("get price error: %w", err)
	}

	if balance < price {
		return fmt.Errorf("buy item error: %w", utils.ErrNotEnoughBalance)
	}

	if _, err = tx.Exec(ReduceCoins, price, userID); err != nil {
		return fmt.Errorf("buy item error: %w", err)
	}

	if _, err := tx.Exec(AddToInventory, userID, itemID); err != nil {
		return fmt.Errorf("add to inventory error: %w", err)
	}

	return tx.Commit()
}

func (u *UserPostgresRepo) GetInventoryInfo(userID int) ([]*entities.Item, error) {
	rows, err := u.DB.Query(GetInventory, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventory := make([]*entities.Item, 0)
	for rows.Next() {
		var itemName string
		var quantity int

		if err := rows.Scan(&itemName, &quantity); err != nil {
			return nil, err
		}

		inventory = append(inventory, &entities.Item{
			ItemType: itemName,
			Quantity: quantity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inventory, nil
}

func (u *UserPostgresRepo) GetCoinsInfo(userID int) (int, error) {
	var coins int
	row := u.DB.QueryRow(GetCoins, userID)
	if err := row.Scan(&coins); err != nil {
		return 0, fmt.Errorf("get coin info error: %w", err)
	}

	return coins, nil
}

func (u *UserPostgresRepo) GetReceiveInfo(userID int) ([]*entities.ReceiveOperation, error) {
	rows, err := u.DB.Query(GetReceiveInfo, userID)
	if err != nil {
		return nil, fmt.Errorf("get receive info: %w", err)
	}

	receives := make([]*entities.ReceiveOperation, 0)
	for rows.Next() {
		var fromUser string
		var amount int
		if err := rows.Scan(&fromUser, &amount); err != nil {
			return nil, fmt.Errorf("get receive info: %w", err) 
		}

		receives = append(receives, &entities.ReceiveOperation{
			FromUser: fromUser,
			Amount: amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get receive info: %w", err) 
	}

	return receives, nil
}


func (u *UserPostgresRepo) GetSentInfo(userID int) ([]*entities.SentOperation, error) {
	rows, err := u.DB.Query(GetSentInfo, userID)
	if err != nil {
		return nil, fmt.Errorf("get sent info: %w", err)
	}

	receives := make([]*entities.SentOperation, 0)
	for rows.Next() {
		var toUser string
		var amount int
		if err := rows.Scan(&toUser, &amount); err != nil {
			return nil, fmt.Errorf("get sent info: %w", err) 
		}

		receives = append(receives, &entities.SentOperation{
			ToUser: toUser,
			Amount: amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get sent info: %w", err) 
	}

	return receives, nil
}

func (u *UserPostgresRepo) GetUserByUsername(username string) (*entities.User, error) {
	user := &entities.User{}

	row := u.DB.QueryRow(GetUser, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", utils.ErrNoUser)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
}

func (u *UserPostgresRepo) Auth(username, password string) (*entities.User, error) {
	user, err := u.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, utils.ErrNoUser) {
			_, err := u.DB.Exec(CreateUser, username, password, InitBalance)
			if err != nil {
				return nil, fmt.Errorf("postgres auth: %w", err)
			}
			user = &entities.User{
				Username: username,
			}

			return user, nil
		}

		return nil, fmt.Errorf("postgres auth: %w", err)
	}

	if user.Password != password {
		return nil, fmt.Errorf("postgres auth: %w", utils.ErrWrongPass)
	}

	return &entities.User{
		Username: username,
	}, nil
}

func (u *UserPostgresRepo) GetUserByID(userID int) (*entities.User, error) {
	user := &entities.User{}

	row := u.DB.QueryRow(GetUserByID, userID)
	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", utils.ErrNoUser)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
}

func (u *UserPostgresRepo) SendCoin(fromUserID, toUserID int, amount int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	fromUser, err := u.GetUserByID(fromUserID)
	if err != nil {
		return err
	}
	if fromUser.Coins < amount {
		return fmt.Errorf("send coin error: %w", utils.ErrNotEnoughBalance)
	}

	_, err = tx.Exec(ReduceCoins, amount, fromUserID)
	if err != nil {
		return fmt.Errorf("send coin error: %w", err)
	}

	_, err = tx.Exec(AddCoins, amount, toUserID)
	if err != nil {
		return fmt.Errorf("send coin error: %w", err)
	}

	_, err = tx.Exec(AddExchangeRecord, fromUserID, toUserID, amount)
	if err != nil {
		return fmt.Errorf("send coin error: %w", err)
	}

	tx.Commit()

	return nil
}
