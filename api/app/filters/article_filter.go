package filters

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type ArticleFilter struct {
}

func (f *ArticleFilter) UserID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (f *ArticleFilter) TagID(db *gorm.DB, id int64) *gorm.DB {
	return db.
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", id)
}

func (f *ArticleFilter) TagIDs(db *gorm.DB, id string) *gorm.DB {
	ids := strings.Split(id, ",")

	return db.
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id IN ?", ids)
}

func (f *ArticleFilter) StartDate(db *gorm.DB, date *time.Time) *gorm.DB {
	return db.Where("created_at >= ?", date)
}

func (f *ArticleFilter) EndDate(db *gorm.DB, date *time.Time) *gorm.DB {
	return db.Where("created_at <= ?", date)
}

func (f *ArticleFilter) Order(db *gorm.DB, order string) *gorm.DB {
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
