package main

import (
	"blog/app"
	"blog/app/models"
	"blog/bootstrap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bootstrap.SetupLogger()
	bootstrap.SetupDatabase()

	_ = app.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Tag{},
		&models.Like{},
		&models.Favorite{},
		&models.Comment{},
		&models.Upvote{},
		&models.Notification{},
	)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	app.DB.Create(&models.User{
		Name:     "admin",
		Password: string(hashedPassword),
		Role:     models.UserRoleManage,
	})
}
