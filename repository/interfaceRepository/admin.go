package interfaces

import "project/domain/entity"

type AdminRepository interface {
	Create( *entity.Admin) error
	CreateOtpKey( string,  string) error
	GetAllUsers( int,  int) ([]entity.User, error)
	GetByEmail( string) (*entity.Admin, error)
	GetById( int) (*entity.User, error)
	GetByPhone( string) (*entity.Admin, error)
	GetOrderByStatus() (int, int, error)
	GetOrders() (int, int, error)
	GetProducts() (int, int, error)
	GetRevenue() (int, error)
	GetUsers() (int, int, error)
	GetUsersBySearch( int,  int,  string) ([]entity.User, error)
	GetstocklessProducts() (*[]entity.Inventory, error)
	Update( *entity.User) error
	
}
