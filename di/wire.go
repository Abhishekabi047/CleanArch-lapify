package di

import (
	"project/config"
	server "project/delivery"
	"project/delivery/handlers"
	repository "project/repository/admin"
	"project/repository/infrastructure"
	repositorya "project/repository/product"
	repositoryb "project/repository/cart"
	repositoryc "project/repository/user"
	repositoryd "project/repository/order"
	usecase "project/usecase/admin"
	usecasea "project/usecase/product"
	usecaseb "project/usecase/cart"
	usecasec "project/usecase/user"
	usecased "project/usecase/order"

)

func InitializeAPI(config *config.Config) error {
	DB, err := infrastructure.ConnectDb(config.DB)
	if err != nil {
		return err
	}

	adminRepository := repository.NewAdminRepository(DB)
	adminUseCase := usecase.NewAdmin(adminRepository)
	

	productRepository := repositorya.NewProductRepository(DB)
	productUseCase:=usecasea.NewProduct(productRepository,&config.S3aws)

	
	
	cartrepository:=repositoryb.NewCartRepository(DB)
	cartUseCase:=usecaseb.NewCart(cartrepository,productRepository)


	UserRepository:=repositoryc.NewUserRepository(DB)
	userUseCase:=usecasec.NewUser(UserRepository,&config.Otp)


	orderrepository:=repositoryd.NewOrderRepository(DB)
	orderUseCase:=usecased.NewOrder(orderrepository,cartrepository,UserRepository,productRepository,&config.Razopay)

	adminHandler := handlers.NewAdminHandler(adminUseCase,productUseCase)
	userHandler:=handlers.NewUserhandler(userUseCase,productUseCase,cartUseCase)
	orderHandler:=handlers.NewOrderHandler(orderUseCase,config.Razopay)


	err = server.NewServerHttp(adminHandler, userHandler, orderHandler)

	return err
}
