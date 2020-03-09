package showdoc

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const ItemFields = "item_id, item_name, item_description, uid, username, password, addtime, last_update_time, item_domain, item_type, is_archived, is_del"

type ItemModel struct {
	db *sql.DB
}

func NewItemModel(db *sql.DB) *ItemModel {
	return &ItemModel{db: db}
}

// GetItems 查询项目列表
func (m *ItemModel) GetItems() ([]Item, error) {
	rows, err := m.db.Query(fmt.Sprintf("SELECT %s FROM item", ItemFields))
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(
			&item.ItemId,
			&item.ItemName,
			&item.ItemDescription,
			&item.Uid,
			&item.UserName,
			&item.Password,
			&item.AddTime,
			&item.LastUpdateTime,
			&item.ItemDomain,
			&item.ItemType,
			&item.IsArchived,
			&item.IsDel,
		); err != nil {
			return nil, errors.Wrap(err, "scan query result failed")
		}

		items = append(items, item)
	}

	return items, nil
}
