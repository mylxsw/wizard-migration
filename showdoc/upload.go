package showdoc

import (
	"fmt"
	"github.com/pkg/errors"
)

const UploadFileFields = "file_id, display_name, file_type, file_size, page_id, item_id, addtime, real_url"

// GetAttachmentsInPage 查询文档包含的所有附件
func (m *PageModel) GetAttachmentsInPage(pageID int64) ([]UploadFile, error) {
	rows, err := m.db.Query(fmt.Sprintf("SELECT %s FROM upload_file WHERE page_id = ?", UploadFileFields), pageID)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	uploadFiles := make([]UploadFile, 0)
	for rows.Next() {
		var uf UploadFile
		if err := rows.Scan(
			&uf.FileID,
			&uf.DisplayName,
			&uf.FileType,
			&uf.FileSize,
			&uf.PageID,
			&uf.ItemID,
			&uf.AddTime,
			&uf.RealURL,
		); err != nil {
			return nil, errors.Wrap(err, "scan query results failed")
		}

		uploadFiles = append(uploadFiles, uf)
	}

	return uploadFiles, nil
}
