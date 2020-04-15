package commands

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/spf13/viper"
	"strings"
)

type helpGroup dgcommand.CommandGroup


func (g helpGroup) helpCommand(ctx dgcommand.Context) {
	ctx = ctx.(*dgcommand.DiscordContext)
	builder := strings.Builder{}
	builder.WriteString("```\n")
	args := strings.Split(ctx.Args()[0], " ")
	commandGroup := dgcommand.CommandGroup(g)

	handler := getRequestedHandler(&commandGroup, args)

	if handler == nil {
		ctx.Reply(fmt.Sprintf("No help found for %v", strings.Join(ctx.Args(), " ")))
		return
	}

	switch handler := handler.(type) {
	case *dgcommand.CommandGroup:
		builder.WriteString(args[len(args)-1])
		builder.WriteByte('\n')
		builder.WriteString(handler.Description)
		builder.WriteByte('\n')
		builder.WriteByte('\n')
		builder.WriteString(getGroupHelp(handler))
	case *dgcommand.Command:
		builder.WriteString(handler.String())
		builder.WriteString("\n")
		builder.WriteString(handler.Description)
	}

	builder.WriteString(fmt.Sprintf(
		"\n\ntype %vhelp <command> to view help for a specific command\n",
		viper.GetString("prefix"),
		))
	builder.WriteString("```")

	ctx.Reply(builder.String())
}

func getGroupHelp(group *dgcommand.CommandGroup) string {
	var builder strings.Builder
	for name, handler := range group.Commands {
		switch handler := handler.(type) {
		case *dgcommand.CommandGroup:
			builder.WriteString(name+" <subcommand>")
		case *dgcommand.Command:
			builder.WriteString(handler.String())
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}
func getRequestedHandler(group *dgcommand.CommandGroup, args []string) dgcommand.Handler {
	if len(args) == 0 || args[0] == "" {
		return group
	}

	var nextHandler dgcommand.Handler = nil
	for i:=0; i<len(args); i++ {
		next := args[i]

		var ok bool
		nextHandler, ok = group.Commands[next]
		if !ok {
			return nil
		}

		switch nextHandler := nextHandler.(type) {
		case *dgcommand.CommandGroup:
			group = nextHandler
		case *dgcommand.Command:
			return nextHandler
		}
	}

	return nextHandler
}
