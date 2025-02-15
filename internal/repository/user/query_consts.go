package repository

var (
	GetBalance = "SELECT balance FROM users WHERE id = $1;"
	GetPrice = "SELECT price FROM items WHERE id = $1;"
	AddRecord = "INSERT INTO purchases (user_id, item_id) VALUES ($1, $2);" //
	GetItemID = "SELECT id FROM items WHERE name = $1;"
	GetUserID = "SELECT id FROM users WHERE username = $1;"
	AddToInventory = "INSERT INTO purchases (user_id, item_id, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = purchases.quantity + 1;"
	GetInventory = "SELECT items.name, purchases.quantity FROM purchases JOIN items ON purchases.item_id = items.id WHERE purchases.user_id = $1;"
	CheckExists = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);"
	GetUser     = "SELECT id, username, password FROM users WHERE username = $1;"
	GetUserByID     = "SELECT id, username, balance FROM users WHERE id = $1;"
	CreateUser = "INSERT INTO users (username, password, balance) VAlUES ($1, $2, $3);"
	ReduceCoins = "UPDATE users SET balance = balance - $1 WHERE id = $2;"
	AddCoins = "UPDATE users SET balance = balance + $1 WHERE id = $2;"
	AddExchangeRecord = "INSERT INTO exchanges (from_id, to_id, amount) VALUES ($1, $2, $3);"
	GetCoins = "SELECT balance FROM users WHERE id = $1;"
	GetReceiveInfo = "SELECT from_id, amount FROM exchanges WHERE to_id = $1;"
	GetSentInfo = "SELECT to_id, amount FROM exchanges WHERE from_id = $1"
)