package main

import (
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/permissions"
	"github.com/plally/dgcommand"
)

var defaults = map[string]bool {
	"commands.sub.delete": false,
	"commands.sub.deleteid": false,
	"commands.sub.add": false,
	"commands.perms.grant": false,
	"commands.perms.deny": false,

}
func registerCommandGroupPermissions(permName string, store permissions.Store, group *dgcommand.CommandGroup) {
	for name, handler := range group.Commands {
		newPermName := permName + "." + name
		switch handler := handler.(type) {
		case *dgcommand.CommandGroup:
			registerCommandGroupPermissions(newPermName, store, handler)
		case *dgcommand.Command:
			defaultValue, ok := defaults[newPermName]
			if !ok {
				defaultValue = true
			}
			store.NewPermission(newPermName, defaultValue)
			handler.Use(middleware.RequireBotPermission(newPermName))
		}
	}
}