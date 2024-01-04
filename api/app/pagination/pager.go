package pagination

import "gorm.io/gorm"

type Pager interface {
	Paginate(db *gorm.DB, dest any) (result *gorm.DB, count int64, err error)
}
