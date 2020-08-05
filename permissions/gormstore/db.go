package gormstore

import (
	"fmt"
	permslib "github.com/plally/FoxBot/permissions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Permission struct {
	Name         string `gorm:"primarykey"`
	DefaultValue bool
}

type UserPermission struct {
	gorm.Model
	Snowflake    string `gorm:"UniqueIndex:idx_snowflake_permission_name"`
	Value bool
	Permission   Permission `gorm:"ForeignKey:PermissionName;References:Name"`
	PermissionName string `gorm:"UniqueIndex:idx_snowflake_permission_name"`
}

func getUserPermissions(db *gorm.DB, snowflake string) ([]UserPermission, error) {
	perms := []UserPermission{}
	err := db.Joins("Permission").Find(&perms, UserPermission{Snowflake: snowflake}).Error
	return perms, err
}

func getPermissions(db *gorm.DB) ([]Permission, error) {
	perms := []Permission{}
	err := db.Find(&perms).Error
	return perms, err
}

func createPermission(db *gorm.DB, name string, defaultValue bool) (Permission, error) {
	perm := Permission{
		Name:         name,
		DefaultValue: defaultValue,
	}
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates:  clause.AssignmentColumns([]string{"default_value"}),
	}).Create(&perm).Error
	return perm, err
}

func setUserPermission(db *gorm.DB, snowflake string, permName string, value bool) (error){
	userPerms := UserPermission{
		Snowflake:      snowflake,
		PermissionName: permName,
		Value:          value,
	}

	fmt.Println(snowflake, permName)
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "snowflake"}, {Name: "permission_name"}},
		DoUpdates:  clause.AssignmentColumns([]string{"value"}),
	}).Create(&userPerms).Error
}

type GormPermissionsStore struct {
	db *gorm.DB
	permissionsCache []Permission
}

func New(db *gorm.DB) GormPermissionsStore{
	_ = db.AutoMigrate(&Permission{})
	_ = db.AutoMigrate(&UserPermission{})
	return GormPermissionsStore{
		db: db,
	}
}

func (store GormPermissionsStore) NewPermission(name string, defaultValue bool) error {
	_, err := createPermission(store.db, name, defaultValue)
	return err
}

func (store GormPermissionsStore) getDefaultPermissions() []Permission {
	if len(store.permissionsCache) == 0 {
		perms, _ := getPermissions(store.db)
		store.permissionsCache = perms
	}

	return store.permissionsCache
}

func (store GormPermissionsStore) GetPermissions(snowflake string) (permslib.UserPerms, error) {
	perms, err := getUserPermissions(store.db, snowflake)
	permMap := permslib.UserPerms(map[string]bool{})

	for _, perm := range store.getDefaultPermissions() {
		permMap[perm.Name] = perm.DefaultValue
	}

	for _, perm := range perms {
		permMap[perm.PermissionName] = perm.Value
	}

	return permMap, err
}

func (store GormPermissionsStore) SetPermission(snowflake string, permName string, value bool) error {
	err := setUserPermission(store.db, snowflake, permName, value)
	return err
}
