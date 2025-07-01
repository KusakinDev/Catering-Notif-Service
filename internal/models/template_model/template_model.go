package templatemodel

import (
	"os"

	"github.com/KusakinDev/Catering-Notif-Service/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Template struct {
	Id          int    `gorm:"primaryKey;autoIncrement"`
	Template    string `gorm:"type:text"`
	Description string `gorm:"type:varchar(50)"`
}

func (template *Template) LoadResetTemplate() error {
	content := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				.card {
					width: 500px;
					margin: 0 auto;
					padding: 20px;
					border: 1px solid #ccc;
					border-radius: 8px;
					box-shadow: 0 2px 4px rgba(0,0,0,0.1);
				}
				.card h1, .card h2, .card p {
					text-align: center;
					margin: 10px 0;
				}
			</style>
		</head>
		<body>
			<div class="card">
				<h1>Catering service</h1>
				<h3>Password reset code: {{.Code}}</h3>
				<p>send by notification service</p>
			</div>
		</body>
		</html>`
	template.Template = string(content)
	template.Description = "reset_code"
	return nil
}

func (template *Template) LoadNewMenuTemplate() error {
	content := `<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Обновление меню</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333;
					background-color: #f8f8f8;
					padding: 20px;
				}
				.container {
					background-color: #fff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					max-width: 600px;
					margin: 0 auto;
				}
				h2, h3 {
					color: #d9534f; 
				}
				a {
					color: #007bff;
					text-decoration: none;
				}
				a:hover {
					text-decoration: underline;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Catering Notification Service</h2>

				<h3>Сообщение для сотрудников</h3>
				<p>В нашем ресторане <strong>обновление меню</strong>. Новое меню доступно по ссылке ниже, и мы рады поделиться с вами этими вкусными новинками:</p>

				<p><a href="{{ .Message }}" target="_blank">Посмотреть новое меню</a></p>

				<p>
					Если у вас есть вопросы или комментарии, пожалуйста, свяжитесь с нами.
				</p>

				<p class="small-text">
					Спасибо, что выбираете нас!<br>
					<strong>Catering Service Team</strong><br>
				</p>
			</div>
		</body>
		</html>`
	template.Template = string(content)
	template.Description = "new_menu"
	return nil
}

func (template *Template) LoadNewDishTemplate() error {
	content := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Новое блюдо в нашем меню</title>
			<style>
				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					background-color: #f9f9f9;
					margin: 0;
					padding: 20px;
					color: #333;
				}
				.container {
					background-color: #ffffff;
					padding: 30px;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					max-width: 600px;
					margin: 30px auto;
				}
				.header {
					text-align: center;
					padding: 10px 0;
					border-bottom: 2px solid #eaeaea;
					margin-bottom: 20px;
				}
				.header h1 {
					color: #333;
				}
				.content {
					margin-bottom: 20px;
				}
				.dish-details {
					border: 1px solid #e0e0e0;
					border-radius: 5px;
					padding: 20px;
					background-color: #f8f8f8;
				}
				.label {
					font-weight: bold;
					color: #007bff;
				}
				.footer {
					text-align: center;
					color: #777;
					font-size: 14px;
					margin-top: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Catering service notification</h1>
					<h2>Приветствуем вас!</h2>
				</div>
				<div class="content">
					<p>В меню нашего ресторана появилось новое блюдо. Более детально можно будет ознакомиться в рабочем меню.</p>
					<div class="dish-details">
						<p><span class="label">Тег:</span> {{.Tag.TagDish}}</p>
						<p><span class="label">Тип блюда:</span> {{.Type.TypeDish}}</p>
						<p><span class="label">Название блюда:</span> {{.Name}}</p>
						<p><span class="label">Категория блюда:</span> {{.Category.CategoryDish}}</p>
					</div>
				</div>
				<div class="footer">
					<p>Спасибо, что работаете у нас!</p>
				</div>
			</div>
		</body>
		</html>`
	template.Template = string(content)
	template.Description = "new_dish"
	return nil
}

func (template *Template) LoadFromFile(filename string, templateName string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		logrus.Errorln(err)
	}
	template.Template = string(content)
	template.Description = templateName
	return nil
}

func (template *Template) DecodeFromContext(c *gin.Context) error {
	if err := c.ShouldBindJSON(&template); err != nil {
		logrus.Error("Error decode JSON: ", err)
		return err
	}
	return nil
}

func (template *Template) AddToTable() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Create(&template).Error
	if err != nil {
		db.CloseDB()
		logrus.Error("Error add to table: ", err)
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) GetFromTableById() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.First(&template).Error
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) GetAllFromTable() ([]Template, int) {
	var db database.DataBase
	db.InitDB()

	var templates []Template

	err := db.Connection.Find(&templates).Error
	if err != nil {
		db.CloseDB()
		logrus.Errorln("Error get all templates from table")
		return []Template{}, 503
	}
	db.CloseDB()
	return templates, 200
}

func (template *Template) GetFromTableByDescription() int {
	var db database.DataBase
	db.InitDB()

	err := db.Connection.Where("description = ?", template.Description).First(&template).Error
	if err != nil {
		db.CloseDB()
		return 503
	}

	db.CloseDB()
	return 200
}

func (template *Template) MigrateToDB(db database.DataBase) error {
	err := db.Connection.AutoMigrate(&Template{})
	if err != nil {
		logrus.Errorln("Error migrate template model :")
		return err
	}
	logrus.Infoln("Success migrate template model :")
	return nil
}
