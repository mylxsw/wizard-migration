package wizard

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          int64
	Name        string
	Description string
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SortLevel   int64
	CatalogId   int64
}

type Page struct {
	ID          int64
	PID         int64
	Title       string
	Description string
	Content     string
	ProjectID   int64
	UserID      int64
	Type        int64
	Status      int64
	HistoryID   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SortLevel   int64
}

type PageHistory struct {
	ID          int64
	PageID      int64
	PID         int64
	Title       string
	Description string
	Content     string
	ProjectID   int64
	Type        int64
	Status      int64
	UserID      int64
	OperatorID  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SortLevel   int64
}

func insert(db *sql.DB, sqlStat string, params []interface{}) (int64, error) {
	//log.Printf("sql -> %s", fmt.Sprintf(""+strings.ReplaceAll(sqlStat, "?", "%v"), params...))
	stmt, err := db.Prepare(sqlStat)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
