package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ecommerce/models"
	"ecommerce/repository"
)

type API struct {
	users    repository.UserRepository
	products repository.ProductRepository
	carts    repository.CartRepository
}

func NewAPI(u repository.UserRepository, p repository.ProductRepository, c repository.CartRepository) *API {
	return &API{users: u, products: p, carts: c}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ---- Kullanıcı İşlemleri ----

func (api *API) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if u.Email == "" || u.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}
	if err := api.users.Create(&u); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	u.Password = ""
	writeJSON(w, http.StatusCreated, u)
}

// PUT /users/{id}
func (api *API) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	u.ID = id
	if err := api.users.Update(&u); err != nil {
		if err == repository.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	u.Password = ""
	writeJSON(w, http.StatusOK, u)
}

// POST /users/login
func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	user, err := api.users.GetByEmail(input.Email)
	if err != nil || user.Password != input.Password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	user.Password = ""
	writeJSON(w, http.StatusOK, user)
}

// GET /users
func (api *API) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users, err := api.users.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	writeJSON(w, http.StatusOK, users)
}

// ---- Ürün İşlemleri ---- (Create/List/Delete/Update aynı mantıkla buraya koyabilirsin; sen zaten yazmıştın)
