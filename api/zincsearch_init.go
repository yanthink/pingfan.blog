package main

import (
	"blog/app"
	"blog/app/models"
	"blog/app/services"
	"blog/bootstrap"
	"blog/config"
	"time"
)

func main() {
	bootstrap.SetupLogger()

	createIndex()
	importArticles()
}

func createIndex() {
	// https://zincsearch-docs.zinc.dev/api/index/create/#update-a-exists-index
	result, err := services.Zincsearch.CreateIndex(map[string]any{
		"name":         config.Zincsearch.Index,
		"storage_type": "disk",
		"shard_num":    config.Zincsearch.ShardNum,
		"mappings": map[string]any{
			"properties": map[string]any{
				"title": map[string]any{
					"type":          "text",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"content": map[string]any{
					"type":          "text",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"view_count": map[string]any{
					"type":          "numeric",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
				"like_count": map[string]any{
					"type":          "numeric",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
				"comment_count": map[string]any{
					"type":          "numeric",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
				"favorite_count": map[string]any{
					"type":          "numeric",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
				"hotness": map[string]any{
					"type":          "numeric",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
				"created_at": map[string]any{
					"type":          "date",
					"format":        time.RFC3339,
					"time_zone":     "+08:00",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false,
				},
			},
		},
		"settings": map[string]any{
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"default": map[string]any{
						"tokenizer": "gse_search",
					},
				},
			},
		},
	})

	if err != nil {
		app.Logger.Sugar().Debugf("创建索引失败：%v", err)
		return
	}

	app.Logger.Sugar().Debugf("创建索引成功 ==== %+v", result)
}

func importArticles() {
	bootstrap.SetupDatabase()

	rows, _ := app.DB.Model(&models.Articles{}).Rows()
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		// ScanRows 方法用于将一行记录扫描至结构体
		if err := app.DB.ScanRows(rows, &article); err == nil {
			err = services.Article.ToZincsearch(&article)

			if err != nil {
				app.Logger.Sugar().Debugf("文档创建失败：%v", err)
			} else {
				app.Logger.Sugar().Debug(article.ID, " 导入成功！")
			}
		}
	}

	app.Logger.Debug("导入完成！")
}
