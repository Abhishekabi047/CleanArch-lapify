package interfaces

import (
	"project/delivery/models"
	"project/domain/entity"
	"time"
)

type OrderRepository interface {
	Create(*entity.Order) (int, error)
	CreateInvoice(*entity.Invoice) (*entity.Invoice, error)
	CreateOrderItems([]entity.OrderItem) error
	DetailedOrderDetails(int) (models.CombinedOrderDetails, error)
	GetAllOrderItems(int) ([]entity.OrderItem, error)
	GetAllOrderList(int, int) ([]entity.Order, error)
	GetAllOrders(int, int, int) ([]entity.Order, error)
	GetByDate(time.Time, time.Time) (*entity.SalesReport, error)
	GetByPaymentMethod(time.Time, time.Time, string) (*entity.SalesReport, error)
	GetByRazorId(string) (*entity.Order, error)
	GetByStatus(int, int, string) ([]entity.Order, error)
	GetOrderById( int) (*entity.Order, error)
	SavePayment(*entity.Charge) ( error)
	Update( *entity.Order) error
	UpdateInvoice( *entity.Invoice) error
	UpdateUserWallet( *entity.User) error
}
