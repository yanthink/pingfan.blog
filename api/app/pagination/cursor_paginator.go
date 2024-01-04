package pagination

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"
)

type CursorPaginator struct {
	paginator.Config `form:"-" json:"-"`
	Cursor           *string `form:"cursor" json:"cursor,omitempty"`
	Limit            int     `form:"limit,default=10" json:"limit,omitempty" binding:"omitempty,min=1,max=100"`
	PageSize         *int    `form:"pageSize" json:"pageSize,omitempty" binding:"omitempty,min=1,max=100"`
}

type Cursor struct {
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

func (cp *CursorPaginator) Paginate(db *gorm.DB, dest any) (result *gorm.DB, cursor *Cursor, err error) {
	if cp.PageSize != nil {
		cp.Limit = *cp.PageSize
	}

	cp.Config.Limit = cp.Limit

	if cp.Cursor != nil {
		qc, _ := cursorDecode(cp.Cursor)
		if qc.After != nil {
			cp.Config.After = *qc.After
		} else if qc.Before != nil {
			cp.Config.Before = *qc.Before
		}
	}

	p := paginator.New(cp)

	var pc paginator.Cursor

	if result, pc, err = p.Paginate(db, dest); err != nil {
		return
	}

	if result.Error != nil {
		err = result.Error
		return
	}

	cursor = &Cursor{}

	if pc.After != nil {
		*pc.After, err = cursorEncode(&Cursor{After: pc.After})
		cursor.After = pc.After
	}

	if pc.Before != nil {
		*pc.Before, err = cursorEncode(&Cursor{Before: pc.Before})
		cursor.Before = pc.Before
	}

	if cursor.After == nil && cursor.Before == nil {
		cursor = nil
	}

	return
}

func cursorEncode(cursor *Cursor) (str string, err error) {
	var jsonStr []byte

	if jsonStr, err = json.Marshal(cursor); err == nil {
		str = base64.StdEncoding.EncodeToString(jsonStr)
	}

	return
}

func cursorDecode(str *string) (cursor Cursor, err error) {
	var jsonStr []byte

	if jsonStr, err = base64.StdEncoding.DecodeString(*str); err == nil {
		err = json.Unmarshal(jsonStr, &cursor)
	}

	return
}
