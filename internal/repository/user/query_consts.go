package repository

var (
	GetBalance = "SELECT balance FROM users WHERE id = $1;"
	GetPrice = "SELECT price FROM items WHERE id = $1;"
	BuyItem = "UPDATE users SET balance = balance - $1 WHERE id = $2;"
	AddRecord = "INSERT INTO purchases (user_id, item_id) VALUES ($1, $2);" //
	GetItemID = "SELECT id FROM items WHERE name = $1;"
	GetUserID = "SELECT id FROM users WHERE username = $1;"
	AddToInventory = "INSERT INTO purchases (user_id, item_id, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = purchases.quantity + 1;"
	GetInventory = "SELECT items.name, purchases.quantity FROM purchases JOIN items ON purchases.item_id = items.id WHERE purchases.user_id = $1;"
	CheckExists = "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);"
	GetUser     = "SELECT id, username, password FROM users WHERE username = $1;"
	CreateUser = "INSERT INTO users (username, password, balance) VAlUES ($1, $2, $3);"
)