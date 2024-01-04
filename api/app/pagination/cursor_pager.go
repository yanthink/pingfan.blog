package pagination

import "gorm.io/gorm"

type CursorPager interface {
	Paginate(db *gorm.DB, dest any) (result *gorm.DB, cursor *Cursor, err error)
}
