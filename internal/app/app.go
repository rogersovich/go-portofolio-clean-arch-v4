package app

import (
	"github.com/rogersovich/go-portofolio-clean-arch-v4/config"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/app/router"
)

func Run() {
	db := config.InitDB()
	r := router.SetupRouter(db)
	r.Run()
}
