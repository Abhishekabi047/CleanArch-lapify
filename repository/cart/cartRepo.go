package cart

import (
	"errors"
	"project/delivery/models"
	"project/domain/entity"
	interfaces "project/repository/interfaceRepository"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) interfaces.CartRepository {
	return &CartRepository{db}
}

func (cr *CartRepository) Create(userid int) (*entity.Cart, error) {
	cart := &entity.Cart{
		UserId: userid,
	}

	if err := cr.db.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}
func (cr *CartRepository) UpdateCart(cart *entity.Cart) error {
	return cr.db.Where("user_id=?", cart.UserId).Save(&cart).Error
}

func (cr *CartRepository) GetByUserid(userid int) (*entity.Cart, error) {
	var cart entity.Cart
	result := cr.db.Where("user_id=?", userid).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &cart, nil
}

func (cr *CartRepository) CreateCartItem(cartitem *entity.CartItem) error {
	if err := cr.db.Create(cartitem).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CartRepository) UpdateCartItem(cartitem *entity.CartItem) error {
	return cr.db.Save(cartitem).Error
}

func (cr *CartRepository) RemoveCartItem(cartitem *entity.CartItem) error {
	return cr.db.Where("product_name=?", cartitem.ProductName).Delete(&cartitem).Error
}

func (cr *CartRepository) GetByName(productName string, cartId int) (*entity.CartItem, error) {
	var cartitem entity.CartItem
	result := cr.db.Where("product_name=? AND cart_id=?", productName, cartId).First(&cartitem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &cartitem, nil
}

func (cr *CartRepository) GetAllCartItems(cartId int) ([]entity.CartItem, error) {
	var cartitems []entity.CartItem
	result := cr.db.Where("cart_id=?", cartId).Find(&cartitems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return cartitems, nil
}

func (cr *CartRepository) RemoveCartItems(cartid int) error {
	var cartitems entity.CartItem
	result := cr.db.Where("cart_id=?", cartid).Delete(&cartitems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func (cr *CartRepository) AddProductToWishlist(product *entity.WishList) error {
	if err := cr.db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CartRepository) GetProductsFromWishlist(id, userid int) (bool, error) {
	var product entity.WishList
	result := cr.db.Where(&entity.WishList{UserId: userid, ProductId: id}).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.New("error finding ticket")
	}
	return true, nil
}

func (cr *CartRepository) GetWishlist(userid int) (*[]entity.WishList, error) {
	var WishList []entity.WishList
	result := cr.db.Where("user_id=?", userid).Find(&WishList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &WishList, nil
}

func (pr *CartRepository) GetAllWishlist(offset, limit int) (*[]models.ProductWithQuantityResponse, error) {
	var productsWithQuantity []models.ProductWithQuantityResponse

	rows, err := pr.db.
		Table("products").
		Select("products.id, products.name, products.price, products.offer_prize, products.size, products.category, products.image_url,products.wish_listed, inventories.quantity").
		Joins("JOIN inventories ON products.id = inventories.product_id").
		Offset(offset).
		Limit(limit).
		Where("products.removed = ?", false).
		Where("products.wish_listed ?", true).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productWithQuantity models.ProductWithQuantityResponse
		if err := pr.db.ScanRows(rows, &productWithQuantity); err != nil {
			return nil, err
		}
		productsWithQuantity = append(productsWithQuantity, productWithQuantity)
	}

	return &productsWithQuantity, nil
}

func (cr *CartRepository) RemoveFromWishlist(id, userid int) error {
	product := &entity.WishList{
		UserId:    userid,
		ProductId: id,
	}
	return cr.db.Where("product_id=?", id).Delete(&product).Error
}

func (cr *CartRepository) GetByType(userid int, addresstype string) (*entity.UserAddress, error) {
	var address entity.UserAddress
	result := cr.db.Where(&entity.UserAddress{User_id: userid, Type: addresstype}).First(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &address, nil
}
