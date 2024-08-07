package routes

import (
	"project/delivery/handlers"

	m "project/delivery/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	r.Use(m.CorsMiddleware)
	r.POST("/admin/login", adminHandler.AdminLoginWithPassword)
	r.GET("/admin/home", m.AdminRetreiveToken, adminHandler.Home)

	r.GET("admin/users", m.AdminRetreiveToken, adminHandler.UsersList)
	r.PUT("/admin/users/toggle-permission/:id", m.AdminRetreiveToken, adminHandler.TogglePermission)
	r.GET("admin/search/users", m.AdminRetreiveToken, adminHandler.SearchUsers)

	r.GET("/admin/categories", m.AdminRetreiveToken, adminHandler.AllCategory)
	r.POST("/admin/categories", m.AdminRetreiveToken, adminHandler.CreateCategory)
	r.PUT("/admin/categories/:id", m.AdminRetreiveToken, adminHandler.EditCategory)
	r.DELETE("/admin/categories/:id", m.AdminRetreiveToken, adminHandler.DeleteCategory)

	r.GET("/admin/products", m.AdminRetreiveToken, adminHandler.AdminProductlist)
	r.POST("/admin/products", m.AdminRetreiveToken, adminHandler.CreateProduct)
	r.PUT("/admin/products/stocks/:id", m.AdminRetreiveToken, adminHandler.AddStock)
	r.GET("/admin/products/details/:productid", m.AdminRetreiveToken, adminHandler.ProductDetailsAdmin)

	r.PATCH("/admin/products/:id", m.AdminRetreiveToken, adminHandler.EditProduct)
	r.DELETE("admin/bin_products/:id", m.AdminRetreiveToken, adminHandler.DeleteProduct)
	r.DELETE("admin/products/:id", m.AdminRetreiveToken, adminHandler.PermanentDeleteProd)

	r.POST("/admin/coupons", m.AdminRetreiveToken, adminHandler.AddCoupon)
	r.GET("/admin/coupons", m.AdminRetreiveToken, adminHandler.AllCoupons)
	r.DELETE("/admin/coupons", m.AdminRetreiveToken, adminHandler.DeleteCoupon)

	r.POST("/admin/banner", m.AdminRetreiveToken, adminHandler.AddBanner)
	r.GET("/admin/banner", m.AdminRetreiveToken, adminHandler.GetBanner)
	r.GET("/admin/bannerid/:id", m.AdminRetreiveToken, adminHandler.GetBannerById)
	r.DELETE("/admin/banner/:id", m.AdminRetreiveToken, adminHandler.DeleteBanner)



	// r.POST("/admin/offer", m.AdminRetreiveToken, adminHandler.AddOffer)
	// r.GET("/admin/offer", m.AdminRetreiveToken, adminHandler.AllOffer)

	r.GET("/admin/stockless/products", m.AdminRetreiveToken, adminHandler.StocklessProducts)

	r.GET("/admin/user-address/:id", m.AdminRetreiveToken, adminHandler.UserAddress)

	r.POST("/admin/product/offer", m.AdminRetreiveToken, adminHandler.AddProductOffer)
	r.POST("/admin/category/offer", m.AdminRetreiveToken, adminHandler.AddCategoryOffer)
	r.POST("/admin/logout", adminHandler.Logout)
	return r
}
