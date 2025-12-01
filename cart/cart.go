package cart

import (
	"ecommerce/models"
	"ecommerce/repository"
)

func CalculateCartTotal(c *models.Cart) float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func CheckoutCart(repo repository.CartRepository, userID int64) (float64, error) {
	cart, err := repo.GetCart(userID)
	if err != nil {
		return 0, err
	}
	if len(cart.Items) == 0 {
		return 0, repository.ErrEmptyCart
	}

	total := CalculateCartTotal(cart)
	_ = repo.ClearCart(userID)
	return total, nil
}
