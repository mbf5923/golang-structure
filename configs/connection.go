package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
	util "mbf5923.com/todo/utils"
	"os"
)

func Connection() *gorm.DB {
	databaseURI := make(chan string, 1)
	if os.Getenv("GO_ENV") != "production" {
		databaseURI <- util.GodotEnv("DATABASE_URI_DEV")
	} else {
		databaseURI <- os.Getenv("DATABASE_URI_PROD")
	}

	db, err := gorm.Open(mysql.Open(<-databaseURI), &gorm.Config{})

	if err != nil {
		defer logrus.Info("Connection to Database Failed")
		logrus.Fatal(err.Error())
	}

	if os.Getenv("GO_ENV") != "production" {
		logrus.Info("Connection to Database Successfully")
	}

	err = db.AutoMigrate(
		&modelUser.EntityUsers{},
		&modelTask.EntityTask{},
		//&model.EntityStudent{},
	)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return db
}
