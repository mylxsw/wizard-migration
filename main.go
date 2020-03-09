package main

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mylxsw/wizard-migration/showdoc"
	"github.com/mylxsw/wizard-migration/wizard"
	"io"
	"log"
	"strings"
	"time"
)

var Version = "1.0"
var GitCommit = "000000000000000000000000000"

var showdocDBConn string
var wizardDBConn string
var importUserID int64
var replaceUrl string
var replaceUrlTo string

func main() {

	flag.StringVar(&showdocDBConn, "showdoc_db", "/Users/mylxsw/codes/github/showdoc/Sqlite/showdoc.db.php", "ShowDoc 数据库文件路径")
	flag.StringVar(&wizardDBConn, "wizard_db", "root:@tcp(127.0.0.1:3306)/wizard_migration", "Wizard 数据库连接地址")
	flag.Int64Var(&importUserID, "import_user_id", 1, "导入后在 Wizard 中使用的 UID")
	flag.StringVar(&replaceUrl, "replace_url", "http://showdoc.local.yunsom.space/server/../Public/Uploads/", "替换地址，用于替换正文中的图片地址")
	flag.StringVar(&replaceUrlTo, "replace_url_to", "/storage/showdoc/", "替换地址，用于替换正文中的图片地址")

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

	tx, err := wizardDB.Begin()
	if err != nil {
		panic(err)
	}

	(func() {
		defer func() {
			if err := recover(); err != nil {
				_ = tx.Rollback()
				log.Printf("migrate to wizard failed: %v", err)
			}
		}()

		if err := migrate(showdocDB, tx); err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}
	})()

	log.Println("migration finished")
}

// migrate 执行迁移
func migrate(showdocDB *sql.DB, wizardDB *sql.Tx) error {
	showdocItemModel := showdoc.NewItemModel(showdocDB)
	showdocPageModel := showdoc.NewPageModel(showdocDB)
	showdocCatalogModel := showdoc.NewCatalogModel(showdocDB)

	wizardProjectModel := wizard.NewProjectModel(wizardDB)
	wizardPageModel := wizard.NewPageModel(wizardDB)

	items, err := showdocItemModel.GetItems()
	if err != nil {
		return err
	}

	// 遍历项目
	for _, item := range items {
		log.Printf("> %s", item.ItemName)

		// 创建项目
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
			return err
		}

		// 查询项目下的目录树
		tree, err := showdocCatalogModel.GetCatalogTreeInItem(item.ItemId)
		if err != nil {
			return err
		}

		// 遍历目录树
		traverseCatalogTree(tree, func(id int64, name string, level int64, pid int64) int64 {
			log.Printf("%s △ %s(%d)", strings.Repeat("     ", int(level)), name, id)

			var catalogID int64 = 0
			// 创建目录（在Wizard中作为Markdown文档）
			if id > 0 {
				catalogID, err = wizardPageModel.CreatePage(wizard.Page{
					PID:                pid,
					Title:              name,
					Description:        "",
					Content:            "",
					ProjectID:          projectID,
					UserID:             importUserID,
					LastModifiedUserID: importUserID,
					Type:               wizard.TypeMarkdown,
					Status:             wizard.StatusNormal,
					HistoryID:          0,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
					SortLevel:          0,
				})
				if err != nil {
					panic(err)
				}
			}

			// 查询目录下所有的文章
			pages, err := showdocPageModel.GetPagesInItemAndCatalog(item.ItemId, id)
			if err != nil {
				panic(err)
			}

			// 新增文章到当前目录
			for _, page := range pages {
				log.Printf("%s      - [%s]", strings.Repeat("     ", int(level)), page.PageTitle)
				pageID, err := wizardPageModel.CreatePage(wizard.Page{
					PID:                catalogID,
					Title:              page.PageTitle,
					Description:        page.PageComments,
					Content:            preProcess(page.PageContent, replaceUrl, replaceUrlTo),
					ProjectID:          projectID,
					UserID:             importUserID,
					LastModifiedUserID: importUserID,
					Type:               wizard.TypeMarkdown,
					Status:             wizard.StatusNormal,
					HistoryID:          0,
					CreatedAt:          time.Unix(page.AddTime, 0),
					UpdatedAt:          time.Unix(page.AddTime, 0),
					SortLevel:          page.SNumber,
				})
				if err != nil {
					panic(err)
				}

				// 添加附件
				attachments, err := showdocPageModel.GetAttachmentsInPage(page.PageId)
				if err != nil {
					panic(err)
				}
				for _, att := range attachments {
					log.Printf("%s           ☁️ [%s]", strings.Repeat("     ", int(level)), att.DisplayName)

					if _, err := wizardPageModel.AddAttachment(wizard.Attachment{
						Name:      att.DisplayName,
						Path:      preProcess(att.RealURL, replaceUrl, replaceUrlTo),
						UserID:    importUserID,
						PageID:    pageID,
						ProjectID: projectID,
						CreatedAt: time.Unix(att.AddTime, 0),
						UpdatedAt: time.Unix(att.AddTime, 0),
					}); err != nil {
						panic(err)
					}
				}

				// 添加历史记录
				histories, err := showdocPageModel.GetPageHistories(page.PageId)
				if err != nil {
					panic(err)
				}

				log.Printf("%s           ❄️️ 导入历史记录 (%d) 条", strings.Repeat("     ", int(level)), len(histories)+1)
				for _, his := range histories {
					pageContent, err := stringDecode(his.PageContent)
					if err != nil {
						log.Printf("decode history page content failed: %v", err)
						continue
					}

					if _, err := wizardPageModel.AddHistory(wizard.PageHistory{
						PageID:      pageID,
						PID:         catalogID,
						Title:       his.PageTitle,
						Description: his.PageComments,
						Content:     pageContent,
						ProjectID:   projectID,
						Type:        wizard.TypeMarkdown,
						Status:      wizard.StatusNormal,
						UserID:      importUserID,
						OperatorID:  importUserID,
						CreatedAt:   time.Unix(his.AddTime, 0),
						UpdatedAt:   time.Unix(his.AddTime, 0),
						SortLevel:   0,
					}); err != nil {
						panic(err)
					}
				}

				// 文档最新内容作为最后一条历史记录
				lastHistoryID, err := wizardPageModel.AddHistory(wizard.PageHistory{
					PageID:      pageID,
					PID:         catalogID,
					Title:       page.PageTitle,
					Description: page.PageComments,
					Content:     page.PageContent,
					ProjectID:   projectID,
					Type:        wizard.TypeMarkdown,
					Status:      wizard.StatusNormal,
					UserID:      importUserID,
					OperatorID:  importUserID,
					CreatedAt:   time.Unix(page.AddTime, 0),
					UpdatedAt:   time.Unix(page.AddTime, 0),
					SortLevel:   0,
				})
				if err != nil {
					panic(err)
				}
				// 更新文档最后历史记录ID
				if _, err := wizardPageModel.UpdatePageHistoryID(pageID, lastHistoryID); err != nil {
					panic(err)
				}
			}

			return catalogID
		}, 0)
	}

	return nil
}

// traverseCatalogTree 遍历目录树
func traverseCatalogTree(tree showdoc.CatalogTree, callback func(id int64, name string, level int64, pid int64) int64, pid int64) {
	newID := callback(tree.ID, tree.Name, tree.Level, pid)

	if len(tree.SubCatalogs) > 0 {
		for _, sc := range tree.SubCatalogs {
			traverseCatalogTree(sc, callback, newID)
		}
	}
}

// preProcess 预处理内容（替换 url 地址)
func preProcess(content string, replaceUrl string, replaceUrlTo string) string {
	return strings.ReplaceAll(content, replaceUrl, replaceUrlTo)
}

// stringDecode 字符串解码 zlib(base64)
func stringDecode(src string) (string, error) {
	pageContent, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	reader, err := zlib.NewReader(bytes.NewReader(pageContent))
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	_, _ = io.Copy(&content, reader)

	return content.String(), nil
}
