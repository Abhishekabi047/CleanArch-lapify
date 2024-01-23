package interfaceUseCase

import (
	"project/domain/entity"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type OrderUseCase interface {
	CreateInvoiceForFailedPayment(int, int) (*entity.Invoice, error)
	ExecutPrintInvoice(int, int) (*gofpdf.Fpdf, error)
	ExecuteAdminCancelOrder(int) error
	ExecuteAdminOrder(int, int) ([]entity.Order, error)
	ExecuteCancelOrder(int) error
	ExecuteCartit(int) (*entity.Cart, error)
	ExecuteInvoiceStripe(int, int) (*entity.Invoice, error)
	ExecuteOrderCod(int, int) (*entity.Invoice, error)
	ExecuteOrderHistory(int, int, int) ([]entity.Order, error)
	ExecuteOrderUpdate(int, string) error
	// Execute(int) (*entity.Order, error)
	ExecutePaymentWallet(int, int) (*entity.Invoice, error)
	ExecuteRazorPay(int, int) (string, int, error)
	ExecuteRazorPayVerification(string, string, string) (*entity.Invoice, error)
	ExecuteSalesReportByDate(time.Time, time.Time) (*entity.SalesReport, error)
	ExecuteSalesReportByPaymentMethod(time.Time, time.Time, string) (*entity.SalesReport, error)
	ExecuteSalesReportByPeriod(string) (*entity.SalesReport, error)
	UpdateInvoiceStatus(int, string) error
	UpdatedUser(int) (*entity.Order, error)
	ExecuteOrderid(int) (*entity.Order, error)
}
