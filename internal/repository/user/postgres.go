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
	InitBalance = 100
)

var (
	ErrNotEnoughBalance = errors.New("Not enough balance to buy item")
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
		return 0, fmt.Errorf("get user id error: %w", err)
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

	fmt.Println("balance:", balance)
	fmt.Println("price:", price)

	if balance < price {
		return ErrNotEnoughBalance
	}

	if _, err = tx.Exec(BuyItem, price, userID); err != nil {
		return fmt.Errorf("buy item error: %w", err)
	}

	if _, err := tx.Exec(AddToInventory, userID, itemID); err != nil {
		return fmt.Errorf("add to inventory error: %w", err)
	}

	return tx.Commit()
}

func (u *UserPostgresRepo) GetInfo(userID int) (*entities.InfoResponse, error) {
	rows, err := u.DB.Query(GetInventory, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &entities.InfoResponse{
		Inventory: make([]*entities.Item, 0),
	}

	for rows.Next() {
		var itemName string
		var quantity int

		if err := rows.Scan(&itemName, &quantity); err != nil {
			return nil, err
		}

		response.Inventory = append(response.Inventory, &entities.Item{
			ItemType: itemName,
			Quantity: quantity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
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

func (u *UserPostgresRepo) SendCoin(fromUserID, toUserID string, amount int) error {
	return nil
}
