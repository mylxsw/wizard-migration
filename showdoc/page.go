package showdoc

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const PageFields = "page_id, author_uid, author_username, item_id, cat_id, page_title, page_content, s_number, addtime, page_comments, is_del"
const PageHistoryFields = "page_history_id, page_id, author_uid, author_username, item_id, cat_id, page_title, page_content, s_number, addtime, page_comments"

type PageModel struct {
	db *sql.DB
}

func NewPageModel(db *sql.DB) *PageModel {
	return &PageModel{db: db}
}

// GetPagesInItemAndCatalog 查询目录下所有的文档
func (m *PageModel) GetPagesInItemAndCatalog(itemId int64, catalogId int64) ([]Page, error) {
	rows, err := m.db.Query(fmt.Sprintf("SELECT %s FROM page WHERE item_id=? AND cat_id=?", PageFields), itemId, catalogId)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	pages := make([]Page, 0)
	for rows.Next() {
		var page Page
		if err := rows.Scan(
			&page.PageId,
			&page.AuthorId,
			&page.AuthorUserName,
			&page.ItemId,
			&page.CatId,
			&page.PageTitle,
			&page.PageContent,
			&page.SNumber,
			&page.AddTime,
			&page.PageComments,
			&page.IsDel,
		); err != nil {
			return nil, errors.Wrap(err, "scan query result failed")
		}

		pages = append(pages, page)
	}

	return pages, nil
}

// GetPageHistories 查询文档历史记录
func (m *PageModel) GetPageHistories(pageId int64) ([]PageHistory, error) {
	rows, err := m.db.Query(fmt.Sprintf("SELECT %s FROM page_history WHERE page_id=?", PageHistoryFields), pageId)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	histories := make([]PageHistory, 0)
	for rows.Next() {
		var history PageHistory
		if err := rows.Scan(
			&history.PageHistoryId,
			&history.PageId,
			&history.AuthorUid,
			&history.AuthorUserName,
			&history.ItemId,
			&history.CatId,
			&history.PageTitle,
			&history.PageContent,
			&history.SNumber,
			&history.AddTime,
			&history.PageComments,
		); err != nil {
			return nil, errors.Wrap(err, "scan query result failed")
		}

		histories = append(histories, history)
	}

	return histories, nil
}
