package cart

import (
	"testing"

	"ecommerce/models"
	"ecommerce/repository"
)

func TestCalculateCartTotal(t *testing.T) {
	c := &models.Cart{
		UserID: 1,
		Items: []models.CartItem{
			{ProductID: 1, Quantity: 2, Price: 10.0},
			{ProductID: 2, Quantity: 1, Price: 5.5},
		},
	}

	got := CalculateCartTotal(c)
	want := 25.5

	if got != want {
		t.Errorf("CalculateCartTotal() = %v, want %v", got, want)
	}
}

func TestCheckoutCart_EmptyCart(t *testing.T) {
	cartRepo := repository.NewInMemoryCartRepo()

	_, err := CheckoutCart(cartRepo, 1)
	if err == nil {
		t.Fatal("expected error for empty cart, got nil")
	}
	if err != repository.ErrEmptyCart {
		t.Fatalf("expected ErrEmptyCart, got %v", err)
	}
}

func TestCheckoutCart_Success(t *testing.T) {
	cartRepo := repository.NewInMemoryCartRepo()

	err := cartRepo.AddItem(1, models.CartItem{
		ProductID: 1,
		Quantity:  2,
		Price:     10.0,
	})
	if err != nil {
		t.Fatalf("unexpected error while adding item: %v", err)
	}

	total, err := CheckoutCart(cartRepo, 1)
	if err != nil {
		t.Fatalf("unexpected error in CheckoutCart: %v", err)
	}
	if total != 20.0 {
		t.Fatalf("expected total 20.0, got %v", total)
	}
}
