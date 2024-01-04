package main

import (
	"blog/app"
	"blog/app/models"
	"blog/app/services"
	"blog/bootstrap"
)

func main() {
	bootstrap.SetupLogger()
	bootstrap.SetupDatabase()
	bootstrap.SetupRedis()

	var articles models.Articles

	app.DB.Model(articles).Find(&articles)

	for _, article := range articles {
		app.DB.Select("TextContent").Omit("UpdatedAt").Save(&article)
		_ = services.Article.ToZincsearch(article)
	}
}
