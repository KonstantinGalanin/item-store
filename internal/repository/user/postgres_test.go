package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/KonstantinGalanin/itemStore/internal/entities"
	"github.com/KonstantinGalanin/itemStore/internal/utils"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	InternalTestError = errors.New("internal error")
)

func TestNewUserPostgresRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserPostgresRepo(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)
}

func TestGetItemID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}
	t.Run("success", func(t *testing.T) {
		expectedID := 1
		itemName := "cup"
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)

		mock.ExpectQuery("SELECT id FROM items WHERE name = (.+);").
			WithArgs(itemName).
			WillReturnRows(rows)

		itemID, err := repo.GetItemID(itemName)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, itemID)
	})

	t.Run("db error", func(t *testing.T) {
		itemName := "table"

		mock.ExpectQuery("SELECT id FROM items WHERE name = (.+);").
			WithArgs(itemName).
			WillReturnError(sql.ErrNoRows)

		itemID, err := repo.GetItemID(itemName)
		assert.Error(t, err)
		assert.Equal(t, 0, itemID)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})

}

func TestGetUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}

	t.Run("success", func(t *testing.T) {
		expectedID := 1
		username := "Konstantin"
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)

		mock.ExpectQuery("SELECT id FROM users WHERE username = (.+);").
			WithArgs(username).
			WillReturnRows(rows)

		userID, err := repo.GetUserID(username)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, userID)
	})

	t.Run("error no users", func(t *testing.T) {
		username := "noUser"

		mock.ExpectQuery("SELECT id FROM users WHERE username = (.+);").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetUserID(username)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, utils.ErrNoUser))
	})

	t.Run("internal error", func(t *testing.T) {
		username := "noUser"

		mock.ExpectQuery("SELECT id FROM users WHERE username = (.+);").
			WithArgs(username).
			WillReturnError(InternalTestError)

		_, err := repo.GetUserID(username)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, InternalTestError))
	})
}

func TestBuyItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}

	t.Run("success", func(t *testing.T) {
		userID := 1
		itemID := 2
		balance := 100
		price := 50
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT balance FROM users WHERE id = (.+);`).WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(balance))
		mock.ExpectQuery(`SELECT price FROM items WHERE id = (.+);`).WithArgs(itemID).
			WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
		mock.ExpectExec(`UPDATE users SET balance = balance - (.+) WHERE id = (.+);`).WithArgs(price, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`INSERT INTO purchases \(user_id, item_id, quantity\) VALUES \((.+), (.+), 1\) ON CONFLICT \(user_id, item_id\) DO UPDATE SET quantity = purchases\.quantity \+ 1;`).WithArgs(userID, itemID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.BuyItem(userID, itemID)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transation begin error", func(t *testing.T) {
		BeginTxError := fmt.Errorf("transaction begin error")
		mock.ExpectBegin().WillReturnError(BeginTxError)

		userID := 1
		itemID := 1
		err := repo.BuyItem(userID, itemID)
		assert.Error(t, err)
		assert.Equal(t, err, BeginTxError)
	})

	t.Run("get balance error", func(t *testing.T) {

		userID := 1
		itemID := 2

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT balance FROM users WHERE id = (.+);`).WithArgs(userID).
			WillReturnError(InternalTestError)
		mock.ExpectCommit()

		err := repo.BuyItem(userID, itemID)
		assert.Error(t, err)
	})

	t.Run("get price error", func(t *testing.T) {

		userID := 1
		itemID := 2
		balance := 100
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT balance FROM users WHERE id = (.+);`).WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(balance))
		mock.ExpectQuery(`SELECT price FROM items WHERE id = (.+);`).WithArgs(itemID).
			WillReturnError(InternalTestError)
		mock.ExpectCommit()

		err := repo.BuyItem(userID, itemID)
		assert.Error(t, err)
		// assert.True(t, errors.Is(err, InternalTestError))
	})
}

func TestGetCoinsInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}

	userID := 1
	expectedCoins := 100

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT balance FROM users WHERE id = (.+);").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(expectedCoins))

		coins, err := repo.GetCoinsInfo(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCoins, coins)
	})

	t.Run("error no user", func(t *testing.T) {
		mock.ExpectQuery("SELECT balance FROM users WHERE id = (.+);").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		coins, err := repo.GetCoinsInfo(userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get coin info error")
		assert.Equal(t, 0, coins)
	})
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}
	userID := 1
	expectedUser := &entities.User{
		ID:       userID,
		Username: "test_user",
		Coins:    500,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, balance FROM users WHERE id = (.+)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance"}).AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Coins))

		user, err := repo.GetUserByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("error no user", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, balance FROM users WHERE id = (.+)").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, utils.ErrNoUser))
	})

	t.Run("internal error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, balance FROM users WHERE id = (.+)").
			WithArgs(userID).
			WillReturnError(InternalTestError)

		user, err := repo.GetUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), InternalTestError.Error())
	})
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}

	username := "test_user"
	expectedUser := &entities.User{
		ID:       1,
		Username: username,
		Password: "hashed_password",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = (.+)").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password))

		user, err := repo.GetUserByUsername(username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("error no user", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = (.+)").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUsername(username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, utils.ErrNoUser))
	})

	t.Run("internal error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = (.+)").
			WithArgs(username).
			WillReturnError(fmt.Errorf("database error"))

		user, err := repo.GetUserByUsername(username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestSendCoin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserPostgresRepo{
		DB: db,
	}

	t.Run("success", func(t *testing.T) {
		fromUserID := 1
		toUserID := 2
		amount := 50

		fromUser := entities.User{
			ID:    fromUserID,
			Coins: 100,
		}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id, username, balance FROM users WHERE id = (.+);").
			WithArgs(fromUserID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "balance"}).
				AddRow(fromUser.ID, "user1", fromUser.Coins))

		mock.ExpectExec(`UPDATE users SET balance = balance \- (.+) WHERE id = (.+);`).
			WithArgs(amount, fromUserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`UPDATE users SET balance = balance \+ (.+) WHERE id = (.+);`).
			WithArgs(amount, toUserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO exchanges \(from_id, to_id, amount\) VALUES \((.+), (.+), (.+)\);`).
			WithArgs(fromUserID, toUserID, amount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.SendCoin(fromUserID, toUserID, amount)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
