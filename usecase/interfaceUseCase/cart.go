package interfaceUseCase

import "project/domain/entity"

type CartUsecase interface {
	ExecuteAddToCart(int, int, int) error
	ExecuteAddWishlist(int, int) error
	ExecuteApplyCoupon(int, string) (int, error)
	ExecuteCart(int) (*entity.Cart, error)
	ExecuteCartItems(int) ([]entity.CartItem, error)
	ExecuteCartitem(int) (*[]entity.CartItem, error)
	ExecuteOfferCheck(int) (*[]entity.Offer, error)
	ExecuteRemoveCartItem(int, int) error
	ExecuteRemoveFromWishList(int, int) error
	ExecuteViewWishlist(int) ([]entity.WishList, error)
	ExecuteDeleteCartItem(int, int) error
}
