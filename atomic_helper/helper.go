package atomic_helper

import (
	"reflect"

	"github.com/Metadiv-Atomic-Engine/atomic/base"
	"github.com/Metadiv-Atomic-Engine/sql"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

/*
Handle one to many relationship
This function will handle the following cases:
1. Create new children
2. Update existing children
3. Delete children that are not in the list
*/
func HandleOne2Many[T any](tx *gorm.DB, parentID uint, parentIDFieldName string, children []T, workspaceID ...uint) (ok bool) {

	idFieldName := strcase.ToCamel(parentIDFieldName)
	idDBFieldName := strcase.ToSnake(parentIDFieldName)

	iChildren := make([]interface{}, len(children))
	for i := range children {
		iChildren[i] = &children[i]
	}

	exists, _ := sql.FindAll[T](tx, sql.And(
		sql.Eq(idDBFieldName, parentID),
		sql.Eq("workspace_id", determineWorkspaceID(workspaceID...)),
	))

	var existIDs = make(map[uint]T)
	for i := range exists {
		id := uint(reflect.ValueOf(exists[i]).FieldByName("ID").Uint())
		existIDs[id] = exists[i]
	}

	var toCreate = make([]T, 0)
	var toUpdate = make([]T, 0)
	var toDelete = make([]T, 0)

	var updatedIDs = make(map[uint]bool)

	for i := range children {
		id := uint(reflect.ValueOf(children[i]).FieldByName("ID").Uint())
		v := reflect.ValueOf(&children[i]).Elem().FieldByName(idFieldName)
		if v.Kind() == reflect.Ptr {
			reflect.ValueOf(&children[i]).Elem().FieldByName(idFieldName).Set(reflect.ValueOf(&parentID))
		} else {
			reflect.ValueOf(&children[i]).Elem().FieldByName(idFieldName).Set(reflect.ValueOf(parentID))
		}
		if id == 0 {
			toCreate = append(toCreate, children[i])
		} else {
			toUpdate = append(toUpdate, children[i])
			updatedIDs[id] = true
		}
	}

	for i := range exists {
		id := uint(reflect.ValueOf(exists[i]).FieldByName("ID").Uint())
		if id != 0 && !updatedIDs[id] {
			toDelete = append(toDelete, exists[i])
		}
	}

	tx1 := tx.Begin()

	if len(toCreate) > 0 {
		wID := determineWorkspaceID(workspaceID...)
		for i := range toCreate {
			_, ok := reflect.ValueOf(new(T)).Interface().(base.IModelWorkspace)
			if ok {
				reflect.ValueOf(&toCreate[i]).Elem().FieldByName("WorkspaceID").SetUint(uint64(wID))
			}
		}
		_, err := sql.SaveAll[T](tx1, toCreate)
		if err != nil {
			tx1.Rollback()
			return false
		}
	}

	if len(toUpdate) > 0 {
		wID := determineWorkspaceID(workspaceID...)
		for i := range toUpdate {
			_, ok := reflect.ValueOf(new(T)).Interface().(base.IModelWorkspace)
			if ok {
				reflect.ValueOf(&toUpdate[i]).Elem().FieldByName("WorkspaceID").SetUint(uint64(wID))
			}
		}
		_, err := sql.SaveAll[T](tx1, toUpdate)
		if err != nil {
			tx1.Rollback()
			return false
		}
	}

	if len(toDelete) > 0 {
		for i := range toDelete {
			err := sql.Delete[T](tx1.Unscoped(), &toDelete[i])
			if err != nil {
				tx1.Rollback()
				return false
			}
		}
	}

	tx1.Commit()
	return true
}

func determineWorkspaceID(workspaceID ...uint) uint {
	var wID uint
	if len(workspaceID) > 0 {
		wID = workspaceID[0]
	} else {
		wID = 1
	}
	return wID
}
