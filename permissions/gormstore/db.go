package gormstore

import (
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
	GuildID        string `gorm:"UniqueIndex:idx_user_guild_perm"`
	UserID         string `gorm:"UniqueIndex:idx_user_guild_perm"`
	Value          bool
	Permission     Permission `gorm:"ForeignKey:PermissionName;References:Name"`
	PermissionName string     `gorm:"UniqueIndex:idx_user_guild_perm"`
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
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"default_value"}),
	}).Create(&perm).Error
	return perm, err
}

func getUserPermissions(db *gorm.DB, guildId string, userID string) ([]UserPermission, error) {
	perms := []UserPermission{}
	err := db.Joins("Permission").Find(&perms, UserPermission{GuildID: guildId, UserID: userID}).Error
	return perms, err
}

func setUserPermission(db *gorm.DB, guildId string, userID string, permName string, value bool) error {
	userPerms := UserPermission{
		UserID:         userID,
		PermissionName: permName,
		GuildID:        guildId,
		Value:          value,
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "permission_name"}, {Name: "guild_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&userPerms).Error
}

type GormPermissionsStore struct {
	db               *gorm.DB
	permissionsCache []Permission
}

func New(db *gorm.DB) GormPermissionsStore {
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

func (store GormPermissionsStore) GetPermissions(guildID string, userID string) (permslib.UserPerms, error) {
	perms, err := getUserPermissions(store.db, guildID, userID)
	permMap := permslib.UserPerms(map[string]bool{})

	for _, perm := range store.getDefaultPermissions() {
		permMap[perm.Name] = perm.DefaultValue
	}

	for _, perm := range perms {
		permMap[perm.PermissionName] = perm.Value
	}

	return permMap, err
}

func (store GormPermissionsStore) SetPermission(guildID string, userID string, permName string, value bool) error {
	err := setUserPermission(store.db, guildID, userID, permName, value)
	return err
}
