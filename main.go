package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mylxsw/wizard-migration/showdoc"
	"github.com/mylxsw/wizard-migration/wizard"
	"log"
	"strings"
	"time"
)

var Version = "1.0"
var GitCommit = "000000000000000000000000000"

var showdocDBConn string
var wizardDBConn string
var importUserID int64

func main() {

	flag.StringVar(&showdocDBConn, "showdoc_db", "/Users/mylxsw/codes/github/showdoc/Sqlite/showdoc.db.php", "ShowDoc 数据库文件路径")
	flag.StringVar(&wizardDBConn, "wizard_db", "root:@tcp(127.0.0.1:3306)/wizard_migration", "Wizard 数据库连接地址")
	flag.Int64Var(&importUserID, "import_user_id", 1, "导入后在 Wizard 中使用的 UID")

	flag.Parse()

	log.Printf("Version=%s. GitCommit=%s", Version, GitCommit)

	showdocDB, err := sql.Open("sqlite3", showdocDBConn)
	if err != nil {
		panic(err)
	}

	wizardDB, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", wizardDBConn))
	if err != nil {
		panic(err)
	}

	showdocItemModel := showdoc.NewItemModel(showdocDB)
	showdocPageModel := showdoc.NewPageModel(showdocDB)
	showdocCatalogModel := showdoc.NewCatalogModel(showdocDB)

	wizardProjectModel := wizard.NewProjectModel(wizardDB)
	wizardPageModel := wizard.NewPageModel(wizardDB)

	items, err := showdocItemModel.GetItems()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		log.Printf("> %s", item.ItemName)

		projectID, err := wizardProjectModel.CreateProject(wizard.Project{
			Name:        item.ItemName,
			Description: item.ItemDescription,
			UserID:      importUserID,
			CreatedAt:   time.Unix(item.AddTime, 0),
			UpdatedAt:   time.Unix(item.AddTime, 0),
			SortLevel:   0,
			CatalogId:   0,
		})
		if err != nil {
			panic(err)
		}

		tree, err := showdocCatalogModel.GetCatalogTreeInItem(item.ItemId)
		if err != nil {
			panic(err)
		}

		traverseCatalogTree(tree, func(id int64, name string, level int64, pid int64) int64 {
			log.Printf("%s △ %s(%d)", strings.Repeat("     ", int(level)), name, id)

			var catalogID int64 = 0
			if id > 0 {
				catalogID, err = wizardPageModel.CreatePage(wizard.Page{
					PID:         pid,
					Title:       name,
					Description: "",
					Content:     "",
					ProjectID:   projectID,
					UserID:      importUserID,
					Type:        1,
					Status:      1,
					HistoryID:   0,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					SortLevel:   0,
				})
				if err != nil {
					panic(err)
				}
			}

			pages, err := showdocPageModel.GetPagesInItemAndCatalog(item.ItemId, id)
			if err != nil {
				panic(err)
			}

			for _, page := range pages {
				log.Printf("%s      - [%s]", strings.Repeat("     ", int(level)), page.PageTitle)
				if _, err := wizardPageModel.CreatePage(wizard.Page{
					PID:         catalogID,
					Title:       page.PageTitle,
					Description: page.PageComments,
					Content:     page.PageContent,
					ProjectID:   projectID,
					UserID:      importUserID,
					Type:        1,
					Status:      1,
					HistoryID:   0,
					CreatedAt:   time.Unix(page.AddTime, 0),
					UpdatedAt:   time.Unix(page.AddTime, 0),
					SortLevel:   page.SNumber,
				}); err != nil {
					panic(err)
				}
			}

			return catalogID
		}, 0)
	}

}

func traverseCatalogTree(tree showdoc.CatalogTree, callback func(id int64, name string, level int64, pid int64) int64, pid int64) {
	newID := callback(tree.ID, tree.Name, tree.Level, pid)

	if len(tree.SubCatalogs) > 0 {
		for _, sc := range tree.SubCatalogs {
			traverseCatalogTree(sc, callback, newID)
		}
	}
}
