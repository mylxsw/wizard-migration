package wizard

import (
	"database/sql"
	"fmt"
	"strings"
)

const ProjectFields = "id, name, description, visibility, user_id, created_at, updated_at, deleted_at, sort_level, catalog_id"

type ProjectModel struct {
	db *sql.DB
}

func NewProjectModel(db *sql.DB) *ProjectModel {
	return &ProjectModel{db: db}
}

func (m *ProjectModel) CreateProject(project Project) (int64, error) {
	sqlStat := fmt.Sprintf("INSERT INTO wz_projects (%s) VALUES (%s)", ProjectFields, strings.Repeat("?, ", len(strings.Split(ProjectFields, ","))-1)+"?")

	return insert(m.db, sqlStat, []interface{}{
		project.ID,
		project.Name,
		project.Description,
		1,
		project.UserID,
		project.CreatedAt,
		project.UpdatedAt,
		nil,
		project.SortLevel,
		project.CatalogId,
	})
}
