package interfaceUseCase

import (
	"mime/multipart"
	"project/delivery/models"
	"project/domain/entity"

	"gorm.io/gorm"
)

type ProductUseCase interface {
	BeginTransaction() *gorm.DB
	ExecuteAddCoupon(*entity.Coupon) error
	ExecuteAddOffer( *entity.Offer) error
	ExecuteAddProductOffer(int,  int) (*entity.Product, error)
	ExecuteAddStock( int,  int) (*entity.Inventory, error)
	ExecuteAvailableCoupons() (*[]entity.Coupon, error)
	ExecuteCategoryOffer( int, int) ([]entity.Product, error)
	ExecuteCreateCategory(entity.Category) (int, error)
	ExecuteCreateInventory( entity.Inventory) error
	ExecuteCreateProduct(entity.Product,  *multipart.FileHeader) (int, error)
	ExecuteCreateProductDetails( entity.ProductDetails) error
	ExecuteDeleteCategory( int) error
	ExecuteDeleteCoupon( string) error
	ExecuteDeleteProduct( int) error
	ExecuteDeleteProductAdd( int) error
	ExecuteEditCategory(entity.Category,  int) error
	ExecuteEditProduct(models.EditProduct,  int) error
	ExecuteGetAllCategory() (*[]entity.Category, error)
	ExecuteGetCategory( entity.Category) (int, error)
	ExecuteGetCategoryId( int) (*entity.Category, error)
	ExecuteGetCouponByCode( string) (*entity.Coupon, error)
	ExecuteGetOffers() (*[]entity.Offer, error)
	ExecuteGetProductById( int) (*entity.Product, error)
	ExecuteProductByCategory( int,  int,  int) ([]entity.Product, error)
	ExecuteProductDetails(int) (*entity.Product, *entity.ProductDetails, error)
	ExecuteProductFilter(string,  int,  int,  int) ([]entity.Product, error)
	ExecuteProductList( int,  int) ([]models.ProductWithQuantityResponse, error)
	ExecuteProductSearch( int, int, string) ([]entity.Product, error)

}
