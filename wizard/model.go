package wizard

import (
	"database/sql"
	"strings"
	"time"
)

const TypeMarkdown int64 = 1
const StatusNormal int64 = 1

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
	ID                 int64
	PID                int64
	Title              string
	Description        string
	Content            string
	ProjectID          int64
	UserID             int64
	LastModifiedUserID int64
	Type               int64
	Status             int64
	HistoryID          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SortLevel          int64
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

type Attachment struct {
	ID        int64
	Name      string
	Path      string
	UserID    int64
	PageID    int64
	ProjectID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func insert(db *sql.Tx, sqlStat string, params []interface{}) (int64, error) {
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

func placeholders(fields string) string {
	return strings.Repeat("?, ", len(strings.Split(fields, ","))-1) + "?"
}
