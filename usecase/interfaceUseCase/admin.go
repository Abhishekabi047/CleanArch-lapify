package interfaceUseCase

import "project/domain/entity"

type AdminUseCase interface {
	ExecuteAdminDashBoard() (*entity.AdminDashboard, error)
	ExecuteAdminLoginWithPassword( string,  string) (int, error)
	ExecuteStocklessProducts() (*[]entity.Inventory, error)
	ExecuteTogglePermission( int) error
	ExecuteUserSearch( int,  int,  string) ([]entity.User, error)
	ExecuteUsersList( int,  int) ([]entity.User, error)
}
