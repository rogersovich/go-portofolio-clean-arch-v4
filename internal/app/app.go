package app

import (
	"os"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/config"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/app/router"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

func Run() {
	utils.InitLogger()

	db := config.InitDB()
	r := router.SetupRouter(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4000"
	}
	r.Run(":" + port)
}
