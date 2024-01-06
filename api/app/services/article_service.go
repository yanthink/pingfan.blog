package services

import (
	"blog/app"
	"blog/app/events"
	"blog/app/filters"
	"blog/app/helpers"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/pagination"
	"blog/config"
	"fmt"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"time"
)

type articleService struct {
}

func (*articleService) Paginate(paginator pagination.Pager) (articles models.Articles, count int64) {
	tx := app.DB.
		Model(articles).
		Select("ID", "Title", "TextContent", "Preview", "ViewCount", "LikeCount", "CommentCount", "FavoriteCount", "CreatedAt", "UserID").
		Order("hotness DESC").
		Order("id DESC").
		Scopes(filters.New(&filters.ArticleFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &articles); err != nil {
		panic(err)
	}

	return
}

func (*articleService) SearchPaginate(keyword string, filters, fields, sortFields []string, size, page int) (articles models.Articles, count int64) {
	from := helpers.Max(page-1, 0) * size

	if len(fields) == 0 {
		fields = []string{"title", "content"}
	}

	must := []map[string]any{
		{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": fields,
			},
		},
	}

	for _, filter := range filters {
		var startAt time.Time

		switch filter {
		case "1month":
			startAt = datetime.BeginOfDay(time.Now().AddDate(0, -1, 0))
		case "3months":
			startAt = datetime.BeginOfDay(time.Now().AddDate(0, -3, 0))
		case "6months":
			startAt = datetime.BeginOfDay(time.Now().AddDate(0, -6, 0))
		case "1year":
			startAt = datetime.BeginOfDay(time.Now().AddDate(-1, 0, 0))
		case "2years":
			startAt = datetime.BeginOfDay(time.Now().AddDate(-2, 0, 0))
		case "3years":
			startAt = datetime.BeginOfDay(time.Now().AddDate(-3, 0, 0))
		}

		if !startAt.IsZero() {
			must = append(must, map[string]any{
				"range": map[string]any{
					"created_at": map[string]any{
						"gte":    startAt,
						"format": time.RFC3339,
					},
				},
			})
		}
	}

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": must,
			},
		},
		"from": from,
		"size": size,
		"highlight": map[string]any{
			"fields": map[string]any{
				"title":   map[string]any{},
				"content": map[string]any{},
			},
			"fragment_size": 200,
		},
		"sort": sortFields,
	}

	result, err := Zincsearch.Search(query)
	if err != nil {
		panic(err)
	}

	count = result.Hits.Total.Value

	var ids []int64
	for _, hit := range result.Hits.Hits {
		id, _ := strconv.ParseInt(hit.ID, 10, 64)
		if id > 0 {
			articles = append(articles, &models.Article{
				ID:         id,
				Highlights: hit.Highlight,
			})
			ids = append(ids, id)
		}
	}

	if len(ids) > 0 {
		rows, _ := app.DB.
			Model(articles).
			Select("ID", "Title", "TextContent", "Preview", "ViewCount", "LikeCount", "CommentCount", "FavoriteCount", "CreatedAt", "UserID").
			Where("id IN ?", ids).Rows()

		idToArticle := map[int64]*models.Article{}

		for rows.Next() {
			var article models.Article
			if err = app.DB.ScanRows(rows, &article); err == nil {
				idToArticle[article.ID] = &article
			}
		}

		for i, article := range articles {
			articles[i] = idToArticle[article.ID]
			articles[i].Highlights = article.Highlights
		}
	}

	return
}

func (*articleService) CursorPaginate(paginator pagination.CursorPager) (articles models.Articles, cursor *pagination.Cursor) {
	tx := app.DB.
		Model(articles).
		Select("ID", "Title", "TextContent", "Preview", "ViewCount", "CommentCount", "CreatedAt").
		Scopes(filters.New(&filters.ArticleFilter{}, paginator))

	var err error

	if _, cursor, err = paginator.Paginate(tx, &articles); err != nil {
		panic(err)
	}

	return
}

func (*articleService) GetByID(id int64, columns ...any) (article *models.Article) {
	if len(columns) == 0 {
		columns = []any{"*"}
	}

	abortIf(app.DB.Select(columns[0], columns[1:]...).First(&article, id).RowsAffected == 0, "数据不存在", responses.CodeModelNotFound, http.StatusNotFound)

	return
}

func (s *articleService) Add(article *models.Article) *models.Article {
	// Omit("Tags.*") 跳过关联的 upsert，这样不会插入新的数据到 tags 表
	if err := app.DB.Omit("Tags.*").Create(&article).Error; err != nil {
		panic(err)
	}

	_ = events.Fire(&events.ArticleChanged{Article: article})

	_ = s.ToZincsearch(article)

	return article
}

func (s *articleService) Update(id int64, article *models.Article, original *models.Article) *models.Article {
	lockKey := fmt.Sprintf("%s_article_update_lock:%d", config.Redis.Prefix, id)
	mutex := app.Redsync.NewMutex(lockKey, redsync.WithTries(1))
	abortIf(mutex.Lock() != nil, "请勿重复操作")
	defer mutex.Unlock()

	article.ID = id

	// Omit("Tags.*") 跳过关联的 upsert，这样不会插入新的数据到 tags 表
	err := app.DB.Model(article).Omit("UpdatedAt", "Tags.*").Association("Tags").Replace(article.Tags)
	if err != nil {
		panic(err)
	}

	err = app.DB.Model(article).Where("id = ?", id).Omit(clause.Associations).Updates(article).Error
	if err != nil {
		panic(err)
	}

	_ = events.Fire(&events.ArticleChanged{Original: original, Article: article})

	_ = s.ToZincsearch(article)

	return article
}

func (s *articleService) Like(id, userId int64) *models.Article {
	lockKey := fmt.Sprintf("%s_like_lock:%d", config.Redis.Prefix, userId)
	mutex := app.Redsync.NewMutex(lockKey, redsync.WithTries(1))
	abortIf(mutex.Lock() != nil, "请勿重复操作")
	defer mutex.Unlock()

	article := s.GetByID(id)
	like := models.Like{ArticleID: article.ID, UserID: userId}

	if app.DB.Unscoped().Where(&like).FirstOrInit(&like).RowsAffected == 0 {
		app.DB.Create(&like)
	} else {
		if like.DeletedAt != nil { // 恢复
			like.DeletedAt = nil
		} else { // 删除
			like.DeletedAt = &gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
		}

		app.DB.Unscoped().Select("DeletedAt").Updates(&like)
	}

	article.HasLiked = like.DeletedAt == nil

	if article.HasLiked {
		article.LikeCount += 1
		app.DB.Model(&article).UpdateColumn("like_count", gorm.Expr("like_count + 1"))
	} else {
		article.LikeCount -= 1
		app.DB.Model(&article).UpdateColumn("like_count", gorm.Expr("like_count - 1"))
	}

	_ = events.Fire(&events.ArticleLiked{Like: &like, Article: article})

	return article
}

func (s *articleService) Favorite(id, userId int64) *models.Article {
	lockKey := fmt.Sprintf("%s_favorite_lock:%d", config.Redis.Prefix, userId)
	mutex := app.Redsync.NewMutex(lockKey, redsync.WithTries(1))
	abortIf(mutex.Lock() != nil, "请勿重复操作")
	defer mutex.Unlock()

	article := s.GetByID(id)
	favorite := models.Favorite{ArticleID: article.ID, UserID: userId}

	_ = app.DB.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "ArticleID"}, {Name: "UserID"}},
			DoUpdates: clause.Assignments(map[string]any{
				"deleted_at": gorm.Expr("IF(ISNULL(deleted_at), now(), null)"),
				"updated_at": gorm.Expr("now()"),
			}),
		}).
		Create(&favorite)

	article.WithUserHasFavorited(userId)

	if article.HasFavorited {
		article.FavoriteCount += 1
		app.DB.Model(&article).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1"))
	} else {
		article.FavoriteCount -= 1
		app.DB.Model(&article).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1"))
	}

	return article
}

// UpdateHotnessInDecayHours 更新衰减时间内的文章热度
func (s *articleService) UpdateHotnessInDecayHours() int64 {
	return app.DB.
		Model(&models.Article{}).
		Where("created_at > ?", time.Now().Add(-models.ArticleInitialHotnessDecayHours*time.Hour)).
		UpdateColumn("hotness", gorm.Expr(fmt.Sprintf(
			"%d * view_count + %d * like_count + %d * comment_count + %d * favorite_count + %d * EXP(-%.8f * TIMESTAMPDIFF(HOUR, created_at, NOW()))",
			models.ArticleViewHotnessWeight,
			models.ArticleLikeHotnessWeight,
			models.ArticleCommentHotnessWeight,
			models.ArticleFavoriteHotnessWeight,
			models.ArticleInitialHotness,
			models.ArticleHotnessDecayFactor,
		))).RowsAffected
}

// SyncViewCountAndUpdateHotness 同步浏览量并且更新热度
func (s *articleService) SyncViewCountAndUpdateHotness(date time.Time) {
	field := models.NewRedisIncrField("articles", "view_count", date)

	var ids []string

	field.Scan(func(id, value string) {
		field.Del(id)

		app.DB.
			Table(field.Table).
			Where("id = ?", id).
			UpdateColumns(map[string]any{
				field.Field: gorm.Expr(fmt.Sprintf("%s + %s", field.Field, value)),
				"hotness": gorm.Expr(fmt.Sprintf(
					"%d * view_count + %d * like_count + %d * comment_count + %d * favorite_count + %d * EXP(-%.8f * TIMESTAMPDIFF(HOUR, created_at, NOW()))",
					models.ArticleViewHotnessWeight,
					models.ArticleLikeHotnessWeight,
					models.ArticleCommentHotnessWeight,
					models.ArticleFavoriteHotnessWeight,
					models.ArticleInitialHotness,
					models.ArticleHotnessDecayFactor,
				)),
			})

		ids = append(ids, id)
	})

	if date.Before(datetime.BeginOfDay(time.Now())) {
		field.Flush()

		// 更新衰减时间内的文章热度
		go func() {
			rows, _ := app.DB.Model(&models.Article{}).
				Where("created_at > ?", time.Now().Add(-models.ArticleInitialHotnessDecayHours*time.Hour)).
				Rows()

			defer rows.Close()

			for rows.Next() {
				var article models.Article
				if err := app.DB.ScanRows(rows, &article); err == nil {
					_ = s.ToZincsearch(&article)
				}
			}
		}()
	}

	// 更新浏览文章
	go func() {
		for _, block := range slice.Chunk(ids, 100) {
			rows, _ := app.DB.Model(&models.Article{}).
				Where("id IN ?", block).
				Where("created_at <= ?", time.Now().Add(-models.ArticleInitialHotnessDecayHours*time.Hour)).
				Rows()

			for rows.Next() {
				var article models.Article
				if err := app.DB.ScanRows(rows, &article); err == nil {
					_ = s.ToZincsearch(&article)
				}
			}

			rows.Close()
		}
	}()
}

func (s *articleService) ToZincsearch(article *models.Article) error {
	_, err := Zincsearch.UpdateDoc(strconv.FormatInt(article.ID, 10), map[string]any{
		"title":          article.Title,
		"content":        article.TextContent,
		"view_count":     article.ViewCount,
		"like_count":     article.LikeCount,
		"comment_count":  article.CommentCount,
		"favorite_count": article.FavoriteCount,
		"hotness":        article.Hotness,
		"created_at":     article.CreatedAt,
	})

	return err
}
