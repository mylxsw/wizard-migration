package showdoc

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const PageFields = "page_id, author_uid, author_username, item_id, cat_id, page_title, page_content, s_number, addtime, page_comments, is_del"

type PageModel struct {
	db *sql.DB
}

func NewPageModel(db *sql.DB) *PageModel {
	return &PageModel{db: db}
}

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
