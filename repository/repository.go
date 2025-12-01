package repository

import (
	"errors"
	"sync"

	"ecommerce/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrEmptyCart    = errors.New("cart is empty")
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetAll() ([]models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
}

type ProductRepository interface {
	Create(p *models.Product) error
	Update(p *models.Product) error
	Delete(id int64) error
	GetAll() ([]models.Product, error)
	GetByID(id int64) (*models.Product, error)
}

type CartRepository interface {
	AddItem(userID int64, item models.CartItem) error
	GetCart(userID int64) (*models.Cart, error)
	ClearCart(userID int64) error
}

// ---- USER ----

type inMemoryUserRepo struct {
	mu    sync.Mutex
	last  int64
	users map[int64]*models.User
}

func NewInMemoryUserRepo() UserRepository {
	return &inMemoryUserRepo{
		users: make(map[int64]*models.User),
	}
}

func (r *inMemoryUserRepo) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.last++
	user.ID = r.last
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return ErrNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) GetAll() ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]models.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, *u)
	}
	return result, nil
}

func (r *inMemoryUserRepo) GetByEmail(email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			tmp := *u
			return &tmp, nil
		}
	}
	return nil, ErrNotFound
}

func (r *inMemoryUserRepo) GetByID(id int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	tmp := *u
	return &tmp, nil
}

// ---- PRODUCT ----

type inMemoryProductRepo struct {
	mu       sync.Mutex
	last     int64
	products map[int64]*models.Product
}

func NewInMemoryProductRepo() ProductRepository {
	return &inMemoryProductRepo{
		products: make(map[int64]*models.Product),
	}
}

func (r *inMemoryProductRepo) Create(p *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.last++
	p.ID = r.last
	r.products[p.ID] = p
	return nil
}

func (r *inMemoryProductRepo) Update(p *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[p.ID]; !ok {
		return ErrNotFound
	}
	r.products[p.ID] = p
	return nil
}

func (r *inMemoryProductRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return ErrNotFound
	}
	delete(r.products, id)
	return nil
}

func (r *inMemoryProductRepo) GetAll() ([]models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]models.Product, 0, len(r.products))
	for _, p := range r.products {
		result = append(result, *p)
	}
	return result, nil
}

func (r *inMemoryProductRepo) GetByID(id int64) (*models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.products[id]
	if !ok {
		return nil, ErrNotFound
	}
	tmp := *p
	return &tmp, nil
}

// ---- CART ----

type inMemoryCartRepo struct {
	mu    sync.Mutex
	carts map[int64]*models.Cart
}

func NewInMemoryCartRepo() CartRepository {
	return &inMemoryCartRepo{
		carts: make(map[int64]*models.Cart),
	}
}

func (r *inMemoryCartRepo) AddItem(userID int64, item models.CartItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart, ok := r.carts[userID]
	if !ok {
		cart = &models.Cart{UserID: userID}
		r.carts[userID] = cart
	}

	for i, it := range cart.Items {
		if it.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			return nil
		}
	}
	cart.Items = append(cart.Items, item)
	return nil
}

func (r *inMemoryCartRepo) GetCart(userID int64) (*models.Cart, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart, ok := r.carts[userID]
	if !ok {
		return &models.Cart{UserID: userID, Items: []models.CartItem{}}, nil
	}
	tmp := *cart
	tmp.Items = append([]models.CartItem(nil), cart.Items...)
	return &tmp, nil
}

func (r *inMemoryCartRepo) ClearCart(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts, userID)
	return nil
}
