package models

type EditUser struct {
	Name  string `json:"name"  `
	Email string `json:"email" `
}

type Signup struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" `
	Phone       string `json:"phone" validate:"required,numeric,len=10"`
	Password    string `json:"password" validate:"required,min=8"`
	ReferalCode string `json:"referalcode"`
}
type CombinedOrderDetails struct {
	OrderId       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	OrderStatus   string  `json:"order_status"`
	PaymentStatus bool    `json:"payment_status"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	State         string  `json:"state" validate:"required"`
	Pin           string  `json:"pin" validate:"required"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
	Address       string  `json:"address"`
}

type ProductWithQuantityResponse struct {
	ID         int    `json:"id" gorm:"column:id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	OfferPrize int    `json:"offerprice"  `
	Size       string `json:"size"`
	Category   int    `json:"category"`
	ImageURL   string `json:"image_url"`
	Quantity   int    `json:"quantity"`
	WishListed bool   `json:"wishlisted"`
}

type EditProduct struct {
	Name          string `json:"name" validate:"required" form:"name"`
	Price         int    `json:"price" validate:"required,number" form:"price"`
	Size          string `json:"size" validate:"required" form:"size"`
	Category      int    `form:"category" gorm:"foreignKey:ID;references:ID" validate:"required,numeric"`
	Description   string `json:"description" validate:"required" form:"description"`
	Specification string `json:"specification" validate:"required" form:"specification"`
	Quantity      int    `json:"quantity" validate:"required" form:"quantity"`
}

type AllUser struct {
	Id          int    `gorm:"primarykey"  json:"id"`
	Name        string `json:"name" validate:"required,alpha" `
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Wallet      int    `json:"wallet"`
	ReferalCode string `json:"referalcode"`
}
