package interfaces

import "project/domain/entity"

type CartRepository interface {
	AddProductToWishlist(*entity.WishList) error
	Create(int) (*entity.Cart, error)
	CreateCartItem(*entity.CartItem) error
	GetAllCartItems(int) ([]entity.CartItem, error)
	GetByName(string, int) (*entity.CartItem, error)
	GetByType(int, string) (*entity.UserAddress, error)
	GetByUserid(int) (*entity.Cart, error)
	GetProductsFromWishlist(int, int) (bool, error)
	GetWishlist(int) (*[]entity.WishList, error)
	RemoveCartItem(*entity.CartItem) error
	RemoveCartItems(int) error
	RemoveFromWishlist(int, int) error
	UpdateCart(*entity.Cart) error
	UpdateCartItem(*entity.CartItem) error
}
