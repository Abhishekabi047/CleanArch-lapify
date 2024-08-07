package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"
	"project/config"
	"project/delivery/models"
	"project/domain/entity"
	"project/domain/utils"
	interfaces "project/repository/interfaceRepository"
	"project/usecase/interfaceUseCase"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo interfaces.ProductRepository
	s3          config.S3Bucket
}

func NewProduct(productRepo interfaces.ProductRepository, s3 *config.S3Bucket) interfaceUseCase.ProductUseCase {
	return &ProductUseCase{productRepo: productRepo, s3: *s3}
}

func (pu *ProductUseCase) ExecuteProductList(page, limit int) ([]models.ProductWithQuantityResponse, error) {
	offset := (page - 1) * limit
	productlist, err := pu.productRepo.GetAllProducts(offset, limit)
	if err != nil {
		return nil, err
	} else {
		return *productlist, nil
	}
}

func (pu *ProductUseCase) ExecuteProductDetails(id int) (*entity.Product, *entity.ProductDetails, *entity.Inventory, error) {
	product, err := pu.productRepo.GetProductById(id)
	if err != nil {
		return nil, nil, nil, err
	}
	productdetails, err := pu.productRepo.GetProductDetailsById(id)
	if err != nil {
		return nil, nil, nil, err
	}
	inventory, err := pu.productRepo.GetInventoryByID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	return product, productdetails, inventory, nil
}

func (pu *ProductUseCase) ExecuteCreateProduct(product entity.Product, image *multipart.FileHeader) (int, error) {
	validate := validator.New()
	validate.RegisterValidation("positive", PositiveNumeric)
	if err := validate.Struct(product); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return 0, err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			case "positive":
				errorMsg += fmt.Sprintf("%s should be a positive numeric value; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return 0, fmt.Errorf(errorMsg)
	}
	err := pu.productRepo.GetProductByName(product.Name)
	if err == nil {
		return 0, errors.New("Product already exists")
	}
	newprod := &entity.Product{
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
		Size:     product.Size,
	}
	switch newprod.Size {
	case "128gb":
		fmt.Println("1")
	case "256gb":
		fmt.Println("1")
	case "512gb":
		fmt.Println("1")
	case "1tb":
		fmt.Println("1")
	default:
		return 0, errors.New("size must be any od this value 128gb,256gb,512gb,1tb")

	}

	sess := utils.CreateSession(pu.s3)
	// fmt.Println("sess", sess)

	ImageURL, err := utils.UploadImageToS3(image, sess)
	if err != nil {
		fmt.Println("err:", err)
		return 0, err
	}
	fmt.Println("err:", err)
	newprod.ImageURL = ImageURL
	fmt.Println("image,", ImageURL)
	productid, err := pu.productRepo.CreateProduct(newprod)
	if err != nil {
		return 0, err
	} else {
		return productid, nil
	}
}

func (pu *ProductUseCase) ExecuteEditProduct1(product entity.Product, proddet entity.ProductDetails, inventory entity.Inventory, image *multipart.FileHeader) error {
	err := pu.productRepo.GetProductByName(product.Name)
	if err == nil {
		return errors.New("product already exists")
	}
	existingProduct, err := pu.productRepo.GetProductById(product.ID)
	sess := utils.CreateSession(pu.s3)
	// fmt.Println("sess", sess)
	if image != nil {
		ImageURL, err := utils.UploadImageToS3(image, sess)
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		existingProduct.Name = product.Name
		existingProduct.Price = product.Price
		existingProduct.Category = product.Category
		existingProduct.Size = product.Size
		existingProduct.ImageURL = ImageURL
		err1 := pu.productRepo.UpdateProduct(existingProduct)
		if err1 != nil {
			return err1
		}
	} else {
		existingProduct.Name = product.Name
		existingProduct.Price = product.Price
		existingProduct.Category = product.Category
		existingProduct.Size = product.Size
		err1 := pu.productRepo.UpdateProduct(existingProduct)
		if err1 != nil {
			return err1
		}
	}

	err2 := pu.productRepo.UpdateProductdetails(&proddet)
	if err2 != nil {
		return err2
	}
	err3 := pu.productRepo.UpdateInventory(&inventory)
	if err3 != nil {
		return err3
	}
	return nil
}

//	func PositiveNumeric(fl validator.FieldLevel) bool {
//		value, err := strconv.ParseInt(fl.Field().String(),10, 64)
//		if err != nil {
//			return false
//		}
//		return value >= 0
//	}
func (pu *ProductUseCase) ExecuteAddBanner(image *multipart.FileHeader, name string) error {
	sess := utils.CreateSession(pu.s3)

	ImageURL, err := utils.UploadImageToS3(image, sess)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}
	newBanner := &entity.Banner{
		Name:     name,
		ImageURL: ImageURL,
	}
	err1 := pu.productRepo.CreateBanner(newBanner)
	if err1 != nil {
		return err
	}
	return nil

}

func (pu *ProductUseCase) ExecuteGetAllBanner() ([]entity.Banner, error) {
	res, err := pu.productRepo.GetAllBanner()
	if err != nil {
		return nil, err
	}
	return *res, nil

}

func (pu *ProductUseCase) ExecuteDeleteBanner(id int) error {
	err := pu.productRepo.DeleteBanner(id)
	if err != nil {
		return err
	}
	return nil

}

func (pu *ProductUseCase) ExecuteGetBannerById(id int) (*entity.Banner, error) {
	res, err := pu.productRepo.GetBannerById(id)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func PositiveNumeric(fl validator.FieldLevel) bool {
	value := fl.Field().Int()

	return value >= 0
}

func (pu *ProductUseCase) ExecuteCreateProductDetails(details entity.ProductDetails) error {
	validate := validator.New()
	if err := validate.Struct(details); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	productDetails := &entity.ProductDetails{
		ProductID:     details.ProductID,
		Description:   details.Description,
		Specification: details.Specification,
	}
	err := pu.productRepo.CreateProductDetails(productDetails)
	if err != nil {
		return errors.New("creating details failed")
	} else {
		return nil
	}
}

func (pt *ProductUseCase) ExecuteEditProduct(product models.EditProduct, id int, image *multipart.FileHeader) error {
	validate := validator.New()
	validate.RegisterValidation("positive", PositiveNumeric)
	if err := validate.Struct(product); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			case "positive":
				errorMsg += fmt.Sprintf("%s should be a positive numeric value; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	existingProduct, err := pt.productRepo.GetProductById(id)
	sess := utils.CreateSession(pt.s3)
	// fmt.Println("sess", sess)
	if image != nil {
		ImageURL, err := utils.UploadImageToS3(image, sess)
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		existingProduct.Name = product.Name
		existingProduct.Price = product.Price
		existingProduct.Category = product.Category
		existingProduct.Size = product.Size
		existingProduct.ImageURL = ImageURL
		err1 := pt.productRepo.UpdateProduct(existingProduct)
		if err1 != nil {
			return err1
		}
	} else {
		existingProduct.Name = product.Name
		existingProduct.Price = product.Price
		existingProduct.Category = product.Category
		existingProduct.Size = product.Size
		err1 := pt.productRepo.UpdateProduct(existingProduct)
		if err1 != nil {
			return err1
		}
	}

	des, err := pt.productRepo.GetProductDescriptionByID(id)
	if err != nil {
		return err
	}
	des.Description = product.Description
	des.Specification = product.Specification

	err2 := pt.productRepo.UpdateProductdetails(des)
	if err2 != nil {
		return err2
	}

	inventory, err := pt.productRepo.GetInventoryByID(id)
	if err != nil {
		return err
	}
	inventory.Quantity = product.Quantity
	err3 := pt.productRepo.UpdateInventory(inventory)
	if err3 != nil {
		return err3
	}

	return nil

}

func (de *ProductUseCase) ExecuteDeleteProduct(id int) error {
	result, err := de.productRepo.GetProductById(id)
	if err != nil {
		return err
	}
	result.Removed = !result.Removed
	err1 := de.productRepo.UpdateProduct(result)
	if err1 != nil {
		return errors.New("product deleted")
	}
	return nil
}

func (de *ProductUseCase) ExecutePermanentDeleteProduct(id int) error {
	err := de.productRepo.PermanentDelete(id)
	if err != nil {
		return errors.New("product deleted")
	}
	return nil
}

func (pu *ProductUseCase) ExecuteCreateCategory(category entity.Category) (int, error) {
	validate := validator.New()
	if err := validate.Struct(category); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return 0, err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return 0, fmt.Errorf(errorMsg)
	}
	err := pu.productRepo.GetCategoryByName(category.Name)
	if err == nil {
		return 0, errors.New("category already exists")
	}
	newcat := &entity.Category{
		Name:        category.Name,
		Description: category.Description,
	}
	categoryid, err := pu.productRepo.CreateCategory(newcat)
	if err != nil {
		return 0, errors.New("category not created")
	} else {
		return categoryid, nil
	}
}

func (pt *ProductUseCase) ExecuteGetAllCategory() (*[]entity.Category, error) {
	category, err := pt.productRepo.AllCategory()
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (pt *ProductUseCase) ExecuteEditCategory(category entity.Category, id int) error {
	validate := validator.New()
	if err := validate.Struct(category); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	existingCat, err := pt.productRepo.GetCategoryById(id)
	if err != nil {
		return err
	}

	existingCat.Name = category.Name
	existingCat.Description = category.Description

	err = pt.productRepo.UpdateCategory(existingCat)
	if err != nil {
		return err
	}

	return nil
}

func (pu *ProductUseCase) ExecuteDeleteCategory(Id int) error {

	category, err := pu.productRepo.GetCategoryById(Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category does not exist")
		}
		return err
	}

	err = pu.productRepo.DeleteCategory(category.ID)
	if err != nil {
		return err
	}
	err1 := pu.productRepo.DeleteProductByCategory(category.ID)
	if err1 != nil {
		return err
	}

	return nil
}
func (pu *ProductUseCase) ExecuteGetCategory(category entity.Category) (int, error) {
	name, err := pu.productRepo.GetCategoryById(category.ID)
	if err != nil {
		return 0, errors.New("error getting category")
	}
	return name.ID, err
}
func (pu *ProductUseCase) ExecuteCreateInventory(inventory entity.Inventory) error {
	validate := validator.New()
	if err := validate.Struct(inventory); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}

	err := pu.productRepo.CreateInventory(&inventory)
	if err != nil {
		return errors.New("Creating inventory failed")
	} else {
		return nil
	}

}

func (p *ProductUseCase) ExecuteAddCoupon(coupon *entity.Coupon) error {
	validate := validator.New()

	validate.RegisterValidation("positive", PositiveNumeric)
	if err := validate.Struct(coupon); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			case "positive":
				errorMsg += fmt.Sprintf("%s should be a positive numeric value; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}

	err := p.productRepo.CreateCoupon(coupon)
	if err != nil {
		return errors.New("creating coupon failed")
	} else {
		return nil
	}
}

func (p *ProductUseCase) ExecuteAddOffer(offer *entity.Offer) error {
	err := p.productRepo.CreateOffer(offer)
	if err != nil {
		return errors.New("error creating offer")
	} else {
		return nil
	}
}

func (p *ProductUseCase) ExecuteAvailableCoupons() (*[]entity.Coupon, error) {
	coupons, err := p.productRepo.GetAllCoupons()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	avialablecoup := []entity.Coupon{}
	for _, coupons := range *coupons {
		if coupons.UsageLimit != coupons.UsedCount {
			avialablecoup = append(avialablecoup, coupons)
		}
	}
	return &avialablecoup, nil
}

func (pu *ProductUseCase) ExecuteGetCategoryId(id int) (*entity.Category, error) {
	cat, err := pu.productRepo.GetCategoryById(id)
	if err != nil {
		return nil, errors.New("error getting category")
	}
	return cat, err
}

func (pu *ProductUseCase) ExecuteGetCouponByCode(code string) (*entity.Coupon, error) {
	coup, err := pu.productRepo.GetCouponByCode(code)
	if err != nil {
		return nil, err
	}
	return coup, nil
}
func (pu *ProductUseCase) ExecuteGetProductById(id int) (*entity.Product, error) {
	prod, err := pu.productRepo.GetProductById(id)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (pu *ProductUseCase) ExecuteDeleteCoupon(code string) error {
	err := pu.productRepo.DeleteCoupon(code)
	if err != nil {
		return err
	}
	return nil
}

func (pu *ProductUseCase) ExecuteProductSearch(page, limit int, search string) ([]models.ProductWithQuantityResponse, error) {
	offset := (page - 1) * limit
	products, err := pu.productRepo.GetAllProductsSearch(offset, limit, search)
	if err != nil {
		return nil, err
	}
	// var result []entity.Product
	// for _, product := range products {
	// 	result = append(result, entity.Product{
	// 		ID:       product.ID,
	// 		Name:     product.Name,
	// 		Price:    product.Price,
	// 		Category: product.Category,
	// 		ImageURL: product.ImageURL,
	// 		Size:     product.Size,
	// 	})
	// }

	return *products, nil
}
func (pu *ProductUseCase) ExecuteProductByCategory(page, limit, id int) (*[]models.ProductWithQuantityResponse, error) {
	offset := (page - 1) * limit
	products, err := pu.productRepo.GetAllProductsByCategory(offset, limit, id)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (pu *ProductUseCase) ExecuteProductFilter(size string, minPrize, maxPrize, category int) ([]models.ProductWithQuantityResponse, error) {
	products, err := pu.productRepo.GetProductsByFilter(minPrize, maxPrize, category, size)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (pu *ProductUseCase) ExecuteGetOffers() (*[]entity.Offer, error) {
	offers, err := pu.productRepo.GetAllOffers()
	if err != nil {
		return nil, err
	}
	avialableoffers := []entity.Offer{}
	for _, offers := range offers {
		if offers.UsageLimit != offers.UsedCount {
			avialableoffers = append(avialableoffers, offers)
		}
	}
	return &avialableoffers, nil
}

func (pu *ProductUseCase) ExecuteAddProductOffer(productid, offer int) (*entity.Product, error) {

	product, err := pu.productRepo.GetProductById(productid)
	if err != nil {
		return nil, err
	}
	if offer < 0 || offer > 100 {
		return nil, errors.New("Invalid offer percentage")
	}

	amount := float64(offer) / 100.0 * float64(product.Price)
	product.OfferPrize = product.Price - int(amount)
	err1 := pu.productRepo.UpdateProduct(product)
	if err1 != nil {
		return nil, err
	}
	return product, nil
}

func (pu *ProductUseCase) ExecuteCategoryOffer(catid, offer int) ([]entity.Product, error) {

	productlist, err := pu.productRepo.GetProductsByCategoryoffer(catid)
	if err != nil {
		return nil, err
	}
	if offer < 0 || offer > 100 {
		return nil, errors.New("Invalid offer percentage")
	}
	for i := range productlist {
		product := &(productlist)[i]

		amount := float64(offer) / 100.0 * float64(product.Price)
		product.OfferPrize = product.Price - int(amount)
		err := pu.productRepo.UpdateProduct(product)
		if err != nil {
			return nil, err
		}
	}
	return productlist, nil

}

func (pu *ProductUseCase) BeginTransaction() *gorm.DB {
	return pu.productRepo.BeginTransaction()
}

func (au *ProductUseCase) ExecuteDeleteProductAdd(id int) error {
	err := au.productRepo.DeleteProductId(id)
	if err != nil {
		return err
	}
	return nil
}

func (au *ProductUseCase) ExecuteAddStock(productId, stock int) (*entity.Inventory, error) {

	product, err := au.productRepo.GetInventoryByID(productId)
	if err != nil {
		return nil, err
	}
	product.Quantity = product.Quantity + stock
	err1 := au.productRepo.UpdateInventory(product)
	if err1 != nil {
		return nil, errors.New("error updating inventory")
	}
	return product, nil

}
