package interfaces

import (
	"project/delivery/models"
	"project/domain/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AllCategory() (*[]entity.Category, error)
	BeginTransaction() *gorm.DB
	CreateCategory( *entity.Category) (int, error)
	CreateCoupon( *entity.Coupon) error
	CreateInventory( *entity.Inventory) error
	CreateOffer( *entity.Offer) error
	CreateProduct( *entity.Product) (int, error)
	CreateProductDetails( *entity.ProductDetails) error
	DecreaseProductQuantity( *entity.Inventory) error
	DeleteCategory( int) error
	DeleteCoupon( string) error
	DeleteProduct(*entity.Product) error
	DeleteProductByCategory( int) error
	DeleteProductId( int) error
	GetAllCoupons() (*[]entity.Coupon, error)
	GetAllOffers() ([]entity.Offer, error)
	GetAllProducts( int,  int) (*[]models.ProductWithQuantityResponse, error)
	GetCategoryById( int) (*entity.Category, error)
	GetCategoryByName( string) error
	GetCouponByCategory( string) (*entity.Coupon, error)
	GetCouponByCode( string) (*entity.Coupon, error)
	GetInventoryByID( int) (*entity.Inventory, error)
	GetOfferByPrize( int) (*[]entity.Offer, error)
	GetProductById( int) (*entity.Product, error)
	GetProductByName( string) error
	GetProductDetailsById(int) (*entity.ProductDetails, error)
	GetProductsByCategory( int,  int,  int) ([]entity.Product, error)
	GetProductsByCategoryoffer( int) ([]entity.Product, error)
	GetProductsByFilter( int,  int,  int,  string) ([]entity.Product, error)
	GetProductsBySearch( int,  int, string) ([]entity.Product, error)
	UpdateCategory( *entity.Category) error
	UpdateCouponCount( *entity.Coupon) error
	UpdateCouponCounts( *entity.Coupon) error
	UpdateCouponUsage(*entity.UsedCoupon) error
	UpdateInventory( *entity.Inventory) error
	UpdateProduct( *entity.Product) error
	UpdateProductdetails( *entity.ProductDetails) error
	GetProductDescriptionByID( int) (*entity.ProductDetails, error) 
	PermanentDelete( int) error
	GetAllProductsSearch(int,  int, string) (*[]models.ProductWithQuantityResponse, error)
}
