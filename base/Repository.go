package base

import (
	"log"
	"reflect"

	"github.com/Metadiv-Atomic-Engine/sql"
	"gorm.io/gorm"
)

type Repository[T any] struct {
}

func (r *Repository[T]) FindByID(tx *gorm.DB, id uint, workspaceID ...uint) *T {
	cls := make([]*sql.Clause, 0)
	cls = append(cls, sql.Eq("id", id))
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq("workspace_id", r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq("workspace_id", 0),
		))
	}
	e, err := sql.FindOne[T](tx, sql.And(cls...))
	if err != nil {
		log.Println(err)
		return nil
	}
	return e
}

func (r *Repository[T]) FindOne(tx *gorm.DB, clause *sql.Clause, workspaceID ...uint) *T {
	cls := make([]*sql.Clause, 0)
	if clause != nil {
		cls = append(cls, clause)
	}
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq("workspace_id", r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq("workspace_id", 0),
		))
	}
	var clause2 *sql.Clause
	if len(cls) > 0 {
		clause2 = sql.And(cls...)
	}
	e, err := sql.FindOne[T](tx, clause2)
	if err != nil {
		log.Println(err)
		return nil
	}
	return e
}

func (r *Repository[T]) FindAll(tx *gorm.DB, clause *sql.Clause, tableName string, workspaceID ...uint) []T {
	cls := make([]*sql.Clause, 0)
	if clause != nil {
		cls = append(cls, clause)
	}
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq(r.workspaceIDField(tableName), r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq(r.workspaceIDField(tableName), 0),
		))
	}
	var clause2 *sql.Clause
	if len(cls) > 0 {
		clause2 = sql.And(cls...)
	}
	es, err := sql.FindAll[T](tx, clause2)
	if err != nil {
		log.Println(err)
		return nil
	}
	return es
}

func (r *Repository[T]) FindAllComplex(tx *gorm.DB,
	clause *sql.Clause, page *sql.Pagination, sort *sql.Sorting, tableName string, workspaceID ...uint) ([]T, *sql.Pagination) {
	cls := make([]*sql.Clause, 0)
	if clause != nil {
		cls = append(cls, clause)
	}
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq(r.workspaceIDField(tableName), r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq(r.workspaceIDField(tableName), 0),
		))
	}
	var clause2 *sql.Clause
	if len(cls) > 0 {
		clause2 = sql.And(cls...)
	}
	es, page, err := sql.FindAllComplex[T](tx, clause2, page, sort)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return es, page
}

func (r *Repository[T]) Count(tx *gorm.DB, clause *sql.Clause, tableName string, workspaceID ...uint) int64 {
	cls := make([]*sql.Clause, 0)
	if clause != nil {
		cls = append(cls, clause)
	}
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq(r.workspaceIDField(tableName), r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq(r.workspaceIDField(tableName), 0),
		))
	}
	var clause2 *sql.Clause
	if len(cls) > 0 {
		clause2 = sql.And(cls...)
	}
	count, err := sql.Count[T](tx, clause2)
	if err != nil {
		log.Println(err)
		return 0
	}
	return count
}

func (r *Repository[T]) DeleteBy(tx *gorm.DB, clause *sql.Clause, workspaceID ...uint) bool {
	cls := make([]*sql.Clause, 0)
	if clause != nil {
		cls = append(cls, clause)
	}
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		cls = append(cls, sql.Or(
			sql.Eq("workspace_id", r.DetermineWorkspaceID(workspaceID...)),
			sql.Eq("workspace_id", 0),
		))
	}
	var clause2 *sql.Clause
	if len(cls) > 0 {
		clause2 = sql.And(cls...)
	}
	err := sql.DeleteBy[T](tx, clause2)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (r *Repository[T]) Save(tx *gorm.DB, e *T, workspaceID ...uint) *T {
	wID := r.DetermineWorkspaceID(workspaceID...)
	_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
	if ok {
		reflect.ValueOf(e).Elem().FieldByName("WorkspaceID").SetUint(uint64(wID))
	}
	e, err := sql.Save[T](tx, e)
	if err != nil {
		log.Println(err)
		return nil
	}

	return e
}

func (r *Repository[T]) SaveAll(tx *gorm.DB, es []T, workspaceID ...uint) []T {
	wID := r.DetermineWorkspaceID(workspaceID...)
	for i := range es {
		_, ok := reflect.ValueOf(new(T)).Interface().(IModelWorkspace)
		if ok {
			reflect.ValueOf(&es[i]).Elem().FieldByName("WorkspaceID").SetUint(uint64(wID))
		}
	}
	es, err := sql.SaveAll[T](tx, es)
	if err != nil {
		log.Println(err)
		return nil
	}
	return es
}

func (r *Repository[T]) Delete(tx *gorm.DB, e *T) (ok bool) {
	err := sql.Delete[T](tx, e)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (r *Repository[T]) DeleteAll(tx *gorm.DB, es []T) (ok bool) {
	err := sql.DeleteAll[T](tx, es)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (r *Repository[T]) DetermineWorkspaceID(workspaceID ...uint) uint {
	var wID uint
	if len(workspaceID) > 0 {
		wID = workspaceID[0]
	} else {
		wID = 1
	}
	return wID
}

func (r *Repository[T]) workspaceIDField(tableName string) string {
	if tableName != "" {
		return tableName + ".workspace_id"
	}
	return "workspace_id"
}
