package notifnewdish

import (
	dishmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/dish_model"
	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	notificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/notification_model"
	"github.com/gin-gonic/gin"
)

func NotifNewDish(c *gin.Context) (int, string) {
	var dish dishmodel.Dish
	dish.DecodeFromContext(c)

	var notif notificationmodel.Notification
	notif.GetTemplateByTag("email")
	notif.Dish = dish

	var email emailmodel.Email
	emails, code := email.GetAllFromTable()
	if code != 200 {
		return code, "Error get all emails from table"
	}

	for _, email := range emails {
		notif.Email = email
		notif.Send()
	}
	return 200, "Success broadcast"
}
