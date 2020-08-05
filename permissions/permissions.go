package permissions

import (
	"context"
	"fmt"
)

type UserPerms map[string]bool

func (p UserPerms) Has(name string) bool {
	value := p[name]
	return value
}

// stores all permissions for the app
type Store interface {
	NewPermission(name string, defaultValue bool) error
	GetPermissions(userID string) (UserPerms, error)
	SetPermission(userID string, permName string, value bool) error
}

func GetPermissionsIdentifier(guildID string, userID string) string {
	return fmt.Sprintf("%v %v", guildID, userID)
}

func FromContext(ctx context.Context) (perms UserPerms) {
	perms, ok := ctx.Value("userPermissions").(UserPerms)
	if !ok {
		store, _ := ctx.Value("permissionsStore").(Store)
		authorID, ok := ctx.Value("permissionsSnowflake").(string)
		if !ok {
			return perms
		}

		perms, _ = store.GetPermissions(authorID)
		context.WithValue(ctx, "userPermissions", perms)
		return perms
	}
	return perms
}
