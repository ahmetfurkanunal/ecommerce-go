package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ecommerce/models"
)

type SQLRepos struct {
	Users    UserRepository
	Products ProductRepository
	Carts    CartRepository
}

func NewSQLRepos(db *sql.DB) *SQLRepos {
	return &SQLRepos{
		Users:    &sqlUserRepo{db: db},
		Products: &sqlProductRepo{db: db},
		Carts:    &sqlCartRepo{db: db},
	}
}

func withTimeout(db *sql.DB) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

//// USER ////

type sqlUserRepo struct {
	db *sql.DB
}

func (r *sqlUserRepo) Create(user *models.User) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users(name, email, password) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Password,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *sqlUserRepo) Update(user *models.User) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		"UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		user.Name, user.Email, user.Password, user.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlUserRepo) GetAll() ([]models.User, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *sqlUserRepo) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	var u models.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, password FROM users WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqlUserRepo) GetByID(id int64) (*models.User, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	var u models.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, password FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//// PRODUCT ////

type sqlProductRepo struct {
	db *sql.DB
}

func (r *sqlProductRepo) Create(p *models.Product) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		"INSERT INTO products(name, price, category) VALUES (?, ?, ?)",
		p.Name, p.Price, p.Category,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *sqlProductRepo) Update(p *models.Product) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		"UPDATE products SET name = ?, price = ?, category = ? WHERE id = ?",
		p.Name, p.Price, p.Category, p.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlProductRepo) Delete(id int64) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		"DELETE FROM products WHERE id = ?",
		id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlProductRepo) GetAll() ([]models.Product, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, price, category FROM products",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *sqlProductRepo) GetByID(id int64) (*models.Product, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	var p models.Product
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, price, category FROM products WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//// CART ////

type sqlCartRepo struct {
	db *sql.DB
}

func (r *sqlCartRepo) AddItem(userID int64, item models.CartItem) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO cart_items(user_id, product_id, quantity, price)
         VALUES (?, ?, ?, ?)
         ON CONFLICT(user_id, product_id)
         DO UPDATE SET quantity = cart_items.quantity + excluded.quantity`,
		userID, item.ProductID, item.Quantity, item.Price,
	)
	return err
}

func (r *sqlCartRepo) GetCart(userID int64) (*models.Cart, error) {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	rows, err := r.db.QueryContext(ctx,
		`SELECT product_id, quantity, price
         FROM cart_items
         WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart := &models.Cart{
		UserID: userID,
		Items:  []models.CartItem{},
	}

	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.ProductID, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *sqlCartRepo) ClearCart(userID int64) error {
	ctx, cancel := withTimeout(r.db)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		"DELETE FROM cart_items WHERE user_id = ?",
		userID,
	)
	return err
}
