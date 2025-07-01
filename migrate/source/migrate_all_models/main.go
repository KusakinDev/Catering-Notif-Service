package main

import (
	loggerconfig "github.com/KusakinDev/Catering-Notif-Service/internal/config/logger"
	"github.com/KusakinDev/Catering-Notif-Service/internal/database"
	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	templatemodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/template_model"
)

func main() {
	loggerconfig.Init()

	var db database.DataBase
	db.InitDB()

	var email emailmodel.Email
	email.MigrateToDB(db)

	var template templatemodel.Template
	template.MigrateToDB(db)

	var template1 templatemodel.Template
	template1.LoadNewDishTemplate()
	template1.AddToTable()

	var template2 templatemodel.Template
	template2.LoadNewMenuTemplate()
	template2.AddToTable()

	var template3 templatemodel.Template
	template3.LoadResetTemplate()
	template3.AddToTable()

	db.CloseDB()
}
