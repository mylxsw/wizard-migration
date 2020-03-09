package wizard

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const PageFields = "id, pid, title, description, content, project_id, user_id, type, status, last_modified_uid, history_id, deleted_at, created_at, updated_at, sort_level"
const AttachmentFields = "id, name, path, user_id, page_id, project_id, created_at, updated_at, deleted_at"
const PageHistoryFields = "id, page_id, pid, title, description, content, project_id, type, status, user_id, operator_id, created_at, updated_at, sort_level"

type PageModel struct {
	db *sql.Tx
}

func NewPageModel(db *sql.Tx) *PageModel {
	return &PageModel{db: db}
}

// CreatePage 新增文档
func (m *PageModel) CreatePage(page Page) (int64, error) {
	sqlStat := fmt.Sprintf("INSERT INTO wz_pages (%s) VALUES (%s)", PageFields, placeholders(PageFields))

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
		page.LastModifiedUserID,
		page.HistoryID,
		nil,
		page.CreatedAt,
		page.UpdatedAt,
		page.SortLevel,
	})
}

// UpdatePageHistoryID 更新文档历史记录ID
func (m *PageModel) UpdatePageHistoryID(pageID int64, historyID int64) (int64, error) {
	sqlStat := "UPDATE wz_pages SET history_id=? WHERE id=?"
	stmt, err := m.db.Prepare(sqlStat)
	if err != nil {
		return 0, errors.Wrap(err, "prepare statement failed")
	}

	result, err := stmt.Exec(historyID, pageID)
	if err != nil {
		return 0, errors.Wrap(err, "execute update sql failed")
	}

	return result.RowsAffected()
}

// AddHistory 添加文档历史
func (m *PageModel) AddHistory(history PageHistory) (int64, error) {
	sqlStat := fmt.Sprintf("INSERT INTO wz_page_histories (%s) VALUES (%s)", PageHistoryFields, placeholders(PageHistoryFields))
	return insert(m.db, sqlStat, []interface{}{
		history.ID,
		history.PageID,
		history.PID,
		history.Title,
		history.Description,
		history.Content,
		history.ProjectID,
		history.Type,
		history.Status,
		history.UserID,
		history.OperatorID,
		history.CreatedAt,
		history.UpdatedAt,
		history.SortLevel,
	})
}

// AddAttachment 添加附件
func (m *PageModel) AddAttachment(file Attachment) (int64, error) {
	sqlStat := fmt.Sprintf("INSERT INTO wz_attachments (%s) VALUES (%s)", AttachmentFields, placeholders(AttachmentFields))
	return insert(m.db, sqlStat, []interface{}{
		file.ID,
		file.Name,
		file.Path,
		file.UserID,
		file.PageID,
		file.ProjectID,
		file.CreatedAt,
		file.UpdatedAt,
		nil,
	})
}
