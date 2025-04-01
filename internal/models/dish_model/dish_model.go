package dishmodel

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Dish struct {
	Id       int
	Name     string
	Recipe   string
	Type     Type
	Category Category
	Tag      Tag
}

type Category struct {
	Id           int
	CategoryDish string
}

type Tag struct {
	Id      int
	TagDish string
}

type Type struct {
	Id       int
	TypeDish string
}

func (dish *Dish) DecodeFromContext(c *gin.Context) error {
	if err := c.ShouldBindJSON(&dish); err != nil {
		logrus.Error("Error decode JSON: ", err)
		return err
	}
	return nil
}
