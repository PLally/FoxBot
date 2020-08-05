package permissions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/FoxBot/help"
	"github.com/plally/FoxBot/permissions"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"sort"
	"strings"
)

func CommandGroup() *dgcommand.CommandGroup {
	var CommandGroup = dgcommand.Group()

	CommandGroup.Default(dgcommand.HandlerFunc(help.DefaultHelpHandler))

	CommandGroup.Desc("permissions?!!")

	CommandGroup.Command("grant <user> <perm>", func(ctx dgcommand.CommandContext) {
		store, _ := ctx.Value("permissionsStore").(permissions.Store)
		if len(ctx.Message.Mentions) <0 {
			return
		}
		user := ctx.Message.Mentions[0]
		perm := ctx.Args()[1]
		err := store.SetPermission(user.ID, perm, true)
		if err != nil {
			return
		}
		ctx.Reply(fmt.Sprintf("granted permission %v for %v", perm, user.Username))
	})

	CommandGroup.Command("deny <user> <perm>", func(ctx dgcommand.CommandContext) {
		store, _ := ctx.Value("permissionsStore").(permissions.Store)
		if len(ctx.Message.Mentions) <0 {
			return
		}
		user := ctx.Message.Mentions[0]
		perm := ctx.Args()[1]
		err := store.SetPermission(user.ID, perm, false)
		if err != nil {
			return
		}
		ctx.Reply(fmt.Sprintf("denied permission %v for %v", perm, user.Username))
	})

	CommandGroup.Command("list [user...]", func(ctx dgcommand.CommandContext) {
		var user *discordgo.User
		if len(ctx.Message.Mentions) == 0 {
			user = ctx.Message.Author
		} else {
			user = ctx.Message.Mentions[len(ctx.Message.Mentions)-1]
		}

		store, _ := ctx.Value("permissionsStore").(permissions.Store)
		perms, _ := store.GetPermissions(user.ID)
		var lines []string

		permsEmbed := embed.NewEmbed()
		for name, defaultValue := range perms {
			emoji := "❌"
			if defaultValue {
				emoji = "✅"
			}
			lines = append(lines, fmt.Sprintf("`%-30v` %v", name, emoji))
		}
		sort.Strings(lines)
		permsEmbed.SetTitle(fmt.Sprintf("Permissions for %v#%v", user.Username, user.Discriminator), "")
		permsEmbed.Description = strings.Join(lines, "\n")
		ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, permsEmbed.MessageEmbed)
	})

	return CommandGroup
}
