package showdoc

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

const CatalogFields = "cat_id, cat_name, item_id, s_number, addtime, parent_cat_id, level"

type CatalogModel struct {
	db *sql.DB
}

func NewCatalogModel(db *sql.DB) *CatalogModel {
	return &CatalogModel{db: db}
}

// GetCatalogTreeInItem 查询项目下的目录树
func (m *CatalogModel) GetCatalogTreeInItem(itemId int64) (CatalogTree, error) {
	catalogs, err := m.GetCatalogsInItem(itemId)
	if err != nil {
		return CatalogTree{}, err
	}

	return m.traverseCatalogsToTree(CatalogTree{
		ID:          0,
		Name:        "",
		Level:       0,
		SubCatalogs: make([]CatalogTree, 0),
	}, catalogs), nil
}

// traverseCatalogsToTree 遍历目录树，创建时树形结构
func (m *CatalogModel) traverseCatalogsToTree(tree CatalogTree, catalogs []Catalog) CatalogTree {
	subCatalogs, restCatalogs := m.searchAndRemoveCatalogByParentID(catalogs, tree.ID)
	for _, sc := range subCatalogs {
		tree.SubCatalogs = append(tree.SubCatalogs, m.traverseCatalogsToTree(CatalogTree{
			ID:          sc.CatId,
			Name:        sc.CatName,
			Level:       sc.Level - 1,
			SubCatalogs: make([]CatalogTree, 0),
		}, restCatalogs))
	}

	return tree
}

// searchAndRemoveCatalogByParentID 在目录列表中搜索子目录
func (m *CatalogModel) searchAndRemoveCatalogByParentID(catalogs []Catalog, parentID int64) (result []Catalog, rest []Catalog) {
	for _, cat := range catalogs {
		if cat.ParentCatId == parentID {
			result = append(result, cat)
		} else {
			rest = append(rest, cat)
		}
	}

	return
}

// GetCatalogsInItem 查询项目下所有的目录列表
func (m *CatalogModel) GetCatalogsInItem(itemId int64) ([]Catalog, error) {
	rows, err := m.db.Query(fmt.Sprintf("SELECT %s FROM catalog WHERE item_id=?", CatalogFields), itemId)
	if err != nil {
		return nil, errors.Wrap(err, "query failed")
	}

	catalogs := make([]Catalog, 0)
	for rows.Next() {
		var catalog Catalog

		if err := rows.Scan(
			&catalog.CatId,
			&catalog.CatName,
			&catalog.ItemId,
			&catalog.SNumber,
			&catalog.AddTime,
			&catalog.ParentCatId,
			&catalog.Level,
		); err != nil {
			return nil, errors.Wrap(err, "scan query result failed")
		}

		catalogs = append(catalogs, catalog)
	}

	return catalogs, nil
}
