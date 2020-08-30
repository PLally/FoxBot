package permissions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/FoxBot/converters"
	"github.com/plally/FoxBot/help"
	"github.com/plally/FoxBot/permissions"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"sort"
	"strings"
)

func CommandGroup() *dgcommand.CommandGroup {
	var group = dgcommand.Group()

	group.Default(dgcommand.HandlerFunc(help.DefaultHelpHandler))

	group.Desc("permissions?!!")

	group.Command("grant <user> <perm>", func(ctx dgcommand.CommandContext) {
		store := ctx.Value("permissionsStore").(permissions.Store)
		if len(ctx.Message.Mentions) < 0 {
			return
		}
		member := converters.ParseMember(ctx, ctx.Value("member").(string))
		perm := ctx.Value("perm").(string)
		if member == nil {
			ctx.Reply("Invalid user")
			return
		}

		identifier := permissions.GetPermissionsIdentifier(ctx.Message.GuildID, member.User.ID)
		err := store.SetPermission(identifier, perm, true)
		if err != nil {
			return
		}
		ctx.Reply(fmt.Sprintf("granted permission %v for %v", perm, member.User.ID))
	})

	group.Command("deny <user> <perm>", func(ctx dgcommand.CommandContext) {
		store := ctx.Value("permissionsStore").(permissions.Store)
		if len(ctx.Message.Mentions) < 0 {
			return
		}
		member := converters.ParseMember(ctx, ctx.Value("member").(string))
		perm := ctx.Value("perm").(string)
		if member == nil {
			ctx.Reply("Invalid user")
			return
		}

		identifier := permissions.GetPermissionsIdentifier(ctx.Message.GuildID, member.User.ID)
		err := store.SetPermission(identifier, perm, false)
		if err != nil {
			return
		}
		ctx.Reply(fmt.Sprintf("denied permission %v for %v", perm, member.User.Username))
	})

	group.Command("list [user...]", func(ctx dgcommand.CommandContext) {
		userArg := ctx.Value("user").(string)
		var user *discordgo.User
		if len(userArg) == 0 {
			user = ctx.Message.Author
		} else {
			member := converters.ParseMember(ctx, userArg)
			user = member.User
		}

		identifier := permissions.GetPermissionsIdentifier(ctx.Message.GuildID, user.ID)

		store, _ := ctx.Value("permissionsStore").(permissions.Store)
		perms, _ := store.GetPermissions(identifier)
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

	return group
}
