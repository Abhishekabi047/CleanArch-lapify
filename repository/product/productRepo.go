package repository

import (
	"errors"
	"fmt"
	"project/delivery/models"
	"project/domain/entity"
	interfaces "project/repository/interfaceRepository"
	"time"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
	return &ProductRepository{db}
}

func (pr *ProductRepository) GetAllProducts(offset, limit int) (*[]models.ProductWithQuantityResponse, error) {
	var productsWithQuantity []models.ProductWithQuantityResponse

	rows, err := pr.db.
		Table("products").
		Select("products.id, products.name, products.price, products.offer_prize, products.size, products.category, products.image_url,products.wish_listed, inventories.quantity").
		Joins("JOIN inventories ON products.id = inventories.product_id").
		Offset(offset).
		Limit(limit).
		Where("products.removed = ?", false).
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

func (pr *ProductRepository) GetAllProductsByCategory(offset, limit int, category int) (*[]models.ProductWithQuantityResponse, error) {
	var productsWithQuantity []models.ProductWithQuantityResponse

	query := pr.db.
		Table("products").
		Select("products.id, products.name, products.price, products.offer_prize, products.size, products.category, products.image_url, products.wish_listed, inventories.quantity").
		Joins("JOIN inventories ON products.id = inventories.product_id").
		Offset(offset).
		Limit(limit).
		Where("products.removed = ?", false)

	if category != 0 {
		query = query.Where("products.category = ?", category)
	}

	rows, err := query.Rows()
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

func (pr *ProductRepository) GetAllProductsSearch(offset, limit int, search string) (*[]models.ProductWithQuantityResponse, error) {
	var productsWithQuantity []models.ProductWithQuantityResponse

	rows, err := pr.db.
		Table("products").
		Select("products.id, products.name, products.price, products.offer_prize, products.size, products.category, products.image_url,products.wish_listed, inventories.quantity").
		Joins("JOIN inventories ON products.id = inventories.product_id").
		Offset(offset).
		Limit(limit).
		Where("products.removed = ?", false).
		Where("name iLIKE ?", "%"+search+"%").
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

func (au *ProductRepository) GetProductDetailsById(id int) (*entity.ProductDetails, error) {
	var productDetails entity.ProductDetails
	result := au.db.Where("product_id=?", id).Find(&productDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &productDetails, nil
}

func (pr *ProductRepository) GetProductById(id int) (*entity.Product, error) {
	var product entity.Product
	result := pr.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &product, nil
}

func (pr *ProductRepository) GetBannerById(id int) (*entity.Banner, error) {
	var banner entity.Banner
	result := pr.db.First(&banner, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &banner, nil
}

func (pr *ProductRepository) GetAllBanner() (*[]entity.Banner, error) {
	var banner []entity.Banner
	err := pr.db.Find(&banner).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return &banner, nil
}

func (cn *ProductRepository) DeleteBanner(Id int) error {
	err := cn.db.Delete(&entity.Banner{}, Id).Error
	if err != nil {
		return errors.New("Coudnt delete")
	}
	return nil
}

func (pn *ProductRepository) GetProductByName(name string) error {
	var prodname entity.Product
	result := pn.db.Where(&entity.Product{Name: name}).First(&prodname)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil

}
func (ct *ProductRepository) CreateProduct(product *entity.Product) (int, error) {
	if err := ct.db.Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (up *ProductRepository) UpdateProduct(product *entity.Product) error {
	return up.db.Save(product).Error
}

func (up *ProductRepository) UpdateProductdetails(product *entity.ProductDetails) error {
	return up.db.Save(product).Error
}

func (dp *ProductRepository) DeleteProduct(product *entity.Product) error {
	return dp.db.Delete(product).Error
}
func (up *ProductRepository) CreateProductDetails(details *entity.ProductDetails) error {
	return up.db.Create(details).Error
}

func (up *ProductRepository) CreateBanner(details *entity.Banner) error {
	return up.db.Create(details).Error
}

func (cc *ProductRepository) CreateCategory(category *entity.Category) (int, error) {
	if err := cc.db.Create(category).Error; err != nil {
		return 0, errors.New("error creating category")
	}
	return category.ID, nil
}

func (cc *ProductRepository) AllCategory() (*[]entity.Category, error) {
	var category []entity.Category
	err := cc.db.Find(&category).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return &category, nil
}

func (uc *ProductRepository) UpdateCategory(category *entity.Category) error {
	return uc.db.Save(category).Error
}

func (cn *ProductRepository) GetCategoryByName(name string) error {
	var prodname entity.Category
	result := cn.db.Where("name=?", name).First(&prodname)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return errors.New("failed to get category by name: %v")
	}
	return nil

}
func (cr *ProductRepository) GetCategoryById(id int) (*entity.Category, error) {
	var category entity.Category
	result := cr.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &category, nil
}

func (cn *ProductRepository) DeleteCategory(Id int) error {

	err := cn.db.Delete(&entity.Category{}, Id).Error
	if err != nil {
		return errors.New("Coudnt delete")
	}
	return nil

}

func (cn *ProductRepository) PermanentDelete(id int) error {
	tx := cn.db.Begin()

	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec("DELETE FROM products WHERE id=?", id).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete from product: " + err.Error())
	}

	if err := tx.Exec("DELETE FROM product_details WHERE product_id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete from product_details: " + err.Error())
	}

	if err := tx.Exec("DELETE FROM inventories WHERE product_id = ?", id).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete from inventory: " + err.Error())
	}

	// Commit the transaction if all deletions were successful
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (cn *ProductRepository) DeleteProductByCategory(Id int) error {

	err := cn.db.Where("category=?", Id).Delete(&entity.Product{}).Error
	if err != nil {
		return errors.New("Coudnt delete")
	}
	return nil

}

func (pr *ProductRepository) CreateInventory(inventory *entity.Inventory) error {
	return pr.db.Create(inventory).Error
}

func (pr *ProductRepository) CreateCoupon(coupon *entity.Coupon) error {
	if err := pr.db.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetAllCoupons() (*[]entity.Coupon, error) {
	var coupon []entity.Coupon
	currenttime := time.Now()
	err := pr.db.Where("validuntil > ?", currenttime).Find(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (pr *ProductRepository) GetCouponByCode(code string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := pr.db.Where("code=?", code).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (pr *ProductRepository) UpdateCouponCount(coupon *entity.Coupon) error {
	return pr.db.Save(coupon).Error
}

func (repo *ProductRepository) UpdateCouponCounts(coupon *entity.Coupon) error {
	// Assuming 'db' is your GORM instance
	err := repo.db.Model(&entity.Coupon{}).Where("id = ?", coupon.ID).Update("used_count", coupon.UsedCount).Error
	return err
}

func (pr *ProductRepository) UpdateCouponUsage(coupon *entity.UsedCoupon) error {
	if err := pr.db.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) CreateOffer(offer *entity.Offer) error {
	if err := pr.db.Create(offer).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetOfferByPrize(Prize int) (*[]entity.Offer, error) {
	offers := &[]entity.Offer{}
	err := pr.db.Where("min_prize=?", Prize).Find(&offers).Error
	if err != nil {
		return nil, err
	} else if offers == nil {
		return nil, err
	}
	return offers, nil
}

func (pr *ProductRepository) DecreaseProductQuantity(product *entity.Inventory) error {
	exisitingproduct := &entity.Inventory{}
	err := pr.db.Where("product_category=? AND product_id=?", product.ProductCategory, product.ProductId).First(exisitingproduct).Error
	if err != nil {
		return err
	}
	if exisitingproduct.Quantity == 0 {
		return errors.New("out of stock")
	}
	newQuantity := exisitingproduct.Quantity - product.Quantity

	if newQuantity < 0 {
		return fmt.Errorf("There is only %d quantity avialable", exisitingproduct.Quantity)
	}

	err = pr.db.Model(exisitingproduct).Update("quantity", newQuantity).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetCouponByCategory(category string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := pr.db.Where("category=?", category).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (pr *ProductRepository) DeleteCoupon(code string) error {
	err := pr.db.Where("code=?", code).Delete(&entity.Coupon{}).Error
	if err != nil {
		return errors.New("coudnt delete")
	}
	return nil
}

func (ar *ProductRepository) GetProductsBySearch(offset, limit int, search string) ([]entity.Product, error) {
	var products []entity.Product

	err := ar.db.Select("id, name, price, category, image_url, size").Where("name iLIKE ?", "%"+search+"%").Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return products, nil
}

func (ar *ProductRepository) GetProductsByCategory(offset, limit, id int) ([]entity.Product, error) {
	var product []entity.Product

	err := ar.db.Where("category=? AND removed =?", id, false).Offset(offset).Limit(limit).Find(&product).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return product, nil
}
func (ar *ProductRepository) GetProductsByFilter(minPrize, maxPrize, category int, size string) ([]models.ProductWithQuantityResponse, error) {
	var productsWithQuantity []models.ProductWithQuantityResponse

	query := ar.db.
		Table("products").
		Select("products.id, products.name, products.price, products.offer_prize, products.size, products.category, products.image_url, products.wish_listed, inventories.quantity").
		Joins("JOIN inventories ON products.id = inventories.product_id")

	if size != "" {
		query = query.Where("products.size = ?", size)
	}
	if minPrize > 0 {
		query = query.Where("products.price >= ?", minPrize)
	}
	if maxPrize > 0 {
		query = query.Where("products.price <= ?", maxPrize)
	}
	if category > 0 {
		query = query.Where("products.category = ?", category)
	}
	query = query.Where("products.removed = ?", false)

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productWithQuantity models.ProductWithQuantityResponse
		if err := ar.db.ScanRows(rows, &productWithQuantity); err != nil {
			return nil, err
		}
		productsWithQuantity = append(productsWithQuantity, productWithQuantity)
	}

	return productsWithQuantity, nil
}

func (pr *ProductRepository) GetAllOffers() ([]entity.Offer, error) {
	var offer []entity.Offer
	currenttime := time.Now()
	err := pr.db.Where("valid_until > ?", currenttime).Find(&offer).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return offer, nil

}

func (ar *ProductRepository) GetProductsByCategoryoffer(id int) ([]entity.Product, error) {
	var product []entity.Product

	err := ar.db.Where("category=? AND removed =?", id, false).Find(&product).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return product, nil
}

func (pr *ProductRepository) BeginTransaction() *gorm.DB {
	return pr.db.Begin()
}
func (dp *ProductRepository) DeleteProductId(id int) error {
	var product *entity.Product
	return dp.db.Where("id=?", id).Delete(&product).Error
}

func (dp *ProductRepository) GetInventoryByID(id int) (*entity.Inventory, error) {
	var prod entity.Inventory
	err := dp.db.Where("product_id=?", id).First(&prod).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return &prod, nil
}

func (up *ProductRepository) UpdateInventory(inventory *entity.Inventory) error {
	return up.db.Save(inventory).Error
}

func (dp *ProductRepository) GetProductDescriptionByID(id int) (*entity.ProductDetails, error) {
	var prod entity.ProductDetails
	err := dp.db.Where("product_id=?", id).First(&prod).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return &prod, nil
}
