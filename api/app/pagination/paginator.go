package pagination

import (
	"gorm.io/gorm"
	"reflect"
)

type Paginator struct {
	Page     int  `form:"page,default=1" json:"page,default=1" binding:"min=1"`
	Current  *int `form:"current" json:"current" binding:"omitempty,min=1"`
	Limit    int  `form:"limit,default=10" json:"limit,default=10" binding:"omitempty,min=1,max=100"`
	PageSize *int `form:"pageSize" json:"pageSize,omitempty" binding:"omitempty,min=1,max=100"`
}

func (p *Paginator) Paginate(db *gorm.DB, dest any) (result *gorm.DB, count int64, err error) {
	if p.Current != nil {
		p.Page = *p.Current
	}

	if p.PageSize != nil {
		p.Limit = *p.PageSize
	}

	limit := p.Limit
	offset := (p.Page - 1) * limit

	result = db.Offset(offset).Limit(limit).Find(dest)

	if err = result.Error; err != nil {
		return
	}

	// dest 必须是指针类型
	elems := reflect.ValueOf(dest).Elem()

	if elems.Kind() == reflect.Slice {
		if size := elems.Len(); limit > 0 && size > 0 && size < limit {
			count = int64(size + offset)
			return
		}
	}

	db.Offset(-1).Limit(-1).Count(&count)
	return
}
