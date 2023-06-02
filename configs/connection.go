package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	modelTask "mbf5923.com/todo/domain/task/models"
	modelUser "mbf5923.com/todo/domain/user/models"
	util "mbf5923.com/todo/utils"
	"os"
	"strconv"
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

func RedisConnection() *redis.Client {
	redisUri := fmt.Sprintf("%s:%s", util.GodotEnv("REDIS_HOST"), util.GodotEnv("REDIS_PORT"))
	db, err := strconv.ParseInt(util.GodotEnv("REDIS_DB"), 10, 32)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: util.GodotEnv("REDIS_PASSWORD"),
		DB:       int(db),
	})
	_, err = client.Ping(client.Context()).Result()

	if err != nil {
		defer logrus.Info("Connection to Redis Failed")
		logrus.Fatal(err.Error())
	}
	if os.Getenv("GO_ENV") != "production" {
		logrus.Info("Connection to Redis Successfully")
	}
	return client
}
