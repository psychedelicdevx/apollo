package main

import (
	"apollo/internal/label"
	"apollo/internal/platform/database"
	"apollo/internal/user"
	"apollo/internal/wallet"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()

	db.AutoMigrate(&user.User{}, &wallet.WalletSnapshot{}, &wallet.TokenCache{}, &label.Label{})

	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	labelRepo := label.NewRepository(db)
	labelService := label.NewService(labelRepo)
	if err := labelService.Seed(); err != nil {
		log.Println("label seed:", err)
	}
	labelHandler := label.NewHandler(labelService)
	labelHandler.RegisterRoutes(r)

	walletClient := wallet.NewClient()
	walletRepo := wallet.NewRepository(db)
	walletService := wallet.NewService(walletClient, walletRepo, labelService)
	walletHandler := wallet.NewHandler(walletService)
	walletHandler.RegisterRoutes(r, userHandler.QuotaMiddleware())

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
