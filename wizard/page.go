package wizard

import (
	"database/sql"
	"fmt"
	"strings"
)

const PageFields = "id, pid, title, description, content, project_id, user_id, type, status, last_modified_uid, history_id, deleted_at, created_at, updated_at, sort_level"

type PageModel struct {
	db *sql.DB
}

func NewPageModel(db *sql.DB) *PageModel {
	return &PageModel{db: db}
}

func (m *PageModel) CreatePage(page Page) (int64, error) {
	sqlStat := fmt.Sprintf("INSERT INTO wz_pages (%s) VALUES (%s)", PageFields, strings.Repeat("?, ", len(strings.Split(PageFields, ","))-1)+" ?")

	return insert(m.db, sqlStat, []interface{}{
		page.ID,
		page.PID,
		page.Title,
		page.Description,
		page.Content,
		page.ProjectID,
		page.UserID,
		page.Type,
		page.Status,
		0,
		page.HistoryID,
		nil,
		page.CreatedAt,
		page.UpdatedAt,
		page.SortLevel,
	})
}
