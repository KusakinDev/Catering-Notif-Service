package notificationmodel

import (
	"bytes"
	"html/template"
	"log"
	"os"

	emailconfig "github.com/KusakinDev/Catering-Notif-Service/.env/email"
	dishmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/dish_model"
	emailmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/email_model"
	templatemodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/template_model"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Notification struct {
	Dish     dishmodel.Dish
	Email    emailmodel.Email
	Template templatemodel.Template
}

func (notif *Notification) GetTemplateByTag(tag string) int {
	notif.Template.Description = tag
	code := notif.Template.GetFromTableByDescription()
	if code != 200 {
		logrus.Errorln("Error find template")
		return 404
	}
	return 200
}

func (notif *Notification) Send() int {

	tmpl, err := template.New("email").Parse(notif.Template.Template)
	if err != nil {
		logrus.Errorln(err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, notif.Dish)
	if err != nil {
		logrus.Errorln(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailconfig.Email)
	m.SetHeader("To", notif.Email.Email)
	m.SetHeader("Subject", "Catering Service: Новое блюдо в нашем меню!")

	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(emailconfig.Host, emailconfig.Port, emailconfig.Email, emailconfig.Password)

	err = d.DialAndSend(m)
	if err != nil {
		logrus.Error("Error send email: ", err)
		return 404
	}

	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("Ошибка при создании файла: %v", err)
	}
	// Не забывайте закрывать файл в конце
	defer file.Close()

	// Запись строки в файл
	_, err = file.WriteString(body.String())
	if err != nil {
		log.Fatalf("Ошибка при записи в файл: %v", err)
	}

	return 200
}
