package permissions

import (
	"context"
	"github.com/plally/dgcommand"
)

type UserPerms map[string]bool

func (p UserPerms) Has(name string) bool {
	value := p[name]
	return value
}

// stores all permissions for the app
type Store interface {
	NewPermission(name string, defaultValue bool) error
	GetPermissions(GuildID string, userID string) (UserPerms, error)
	SetPermission(GuildID string, userID string, permName string, value bool) error
}
func FromContext(ctx dgcommand.CommandContext) (perms UserPerms) {
	perms, ok := ctx.Value("userPermissions").(UserPerms)
	if !ok {
		store, _ := ctx.Value("permissionsStore").(Store)
		perms, _ = store.GetPermissions(ctx.Message.GuildID, ctx.Message.Author.ID)
		context.WithValue(ctx, "userPermissions", perms)
		return perms
	}
	return perms
}
