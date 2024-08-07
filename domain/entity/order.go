package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model    `json:"-"`
	ID            int    `gorm:"primarykey" json:"orderid"`
	UserId        int    `json:"userid"`
	Addressid     int    `json:"addressid"`
	Total         int    `json:"total"`
	Status        string `json:"status"`
	PaymentMethod string `json:"paymentmethod"`
	PaymentStatus string `json:"payemntstatus"`
	PaymentId     string `json:"paymentid"`
}

type OrderItem struct {
	gorm.Model `json:"-"`
	OrderId    int `json:"orderid"`
	ProductId  int `json:"productid"`
	Category   int `json:"category"`
	Quantity   int `json:"quantity"`
	Prize      int `json:"prize"`
}
type Invoice struct {
	gorm.Model  `json:"-"`
	OrderId     int     `json:"orderid"`
	UserId      int     `json:"userid"`
	AddressType string  `json:"addresstype"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Payment     string  `json:"payment"`
	Status      string  `json:"status"`
	PaymentId   string  `json:"paymentid"`
	Remark      string  `json:"remark" gorm:"default lapify_festiv"`
}
type SalesReport struct {
	TotalSales   float64
	TotalOrders  int64
	AverageOrder float64
}
type Charge struct {
	gorm.Model   `json:"-"`
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receiptMail"`
	ProductName  string `json:"productName"`
}
