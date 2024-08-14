package filters

import (
	"blog/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleFilter struct {
}

func (*ArticleFilter) UserID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (*ArticleFilter) StartDate(db *gorm.DB, date *time.Time, _ any) *gorm.DB {
	return db.Where("created_at >= ?", date)
}

func (*ArticleFilter) EndDate(db *gorm.DB, date *time.Time, _ any) *gorm.DB {
	return db.Where("created_at <= ?", date)
}

func (*ArticleFilter) Order(db *gorm.DB, order string, _ any) *gorm.DB {
	orderBy := clause.OrderBy{}
	if orderByClause, ok := db.Statement.Clauses[orderBy.Name()]; ok {
		orderByClause.Expression = nil
		db.Statement.Clauses[orderBy.Name()] = orderByClause
	}

	switch order {
	case "latest":
		return db.Order("id DESC")
	case "like":
		return db.Order("like_count DESC")
	case "comment":
		return db.Order("comment_count DESC")
	}

	return db.Order("hotness DESC").Order("id DESC")
}

func (*ArticleFilter) Tag(db *gorm.DB, tag *models.Tag, _ any) *gorm.DB {
	subQuery := db.Session(&gorm.Session{NewDB: true}).
		Table("article_tags").
		Select("article_id").
		Scopes(New(&ArticleTagFilter{}, tag))

	return db.Where("id IN (?)", subQuery)
}
