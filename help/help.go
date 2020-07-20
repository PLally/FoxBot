package help

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/spf13/viper"
	"strings"
)

func GetHelp(commandGroup *dgcommand.CommandGroup, args []string) string {
	builder := strings.Builder{}

	builder.WriteString("```\n")
	handler := getRequestedHandler(commandGroup, args)

	if handler == nil {
		return fmt.Sprintf("No help found" )
	}

	switch handler := handler.(type) {
	case *dgcommand.CommandGroup:
		groupName := args[len(args)-1]
		groupName = strings.Trim(groupName, " ")
		builder.WriteString(groupName)
		if len(groupName) > 0 {
			builder.WriteString(" <subcommand>")
		}
		builder.WriteByte('\n')
		builder.WriteString(handler.Description)
		builder.WriteByte('\n')
		builder.WriteByte('\n')
		builder.WriteString("Commands:\n")
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
	return builder.String()
}

func getGroupHelp(group *dgcommand.CommandGroup) string {
	var builder strings.Builder
	for name, handler := range group.Commands {
		switch handler := handler.(type) {
		case *dgcommand.CommandGroup:
			builder.WriteString(name)
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
	for i := 0; i < len(args); i++ {
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

func DefaultHelpHandler(ctx dgcommand.CommandContext) {
	paths, ok := ctx.Value("handlerPath").([]string)
	if !ok { return }

	handler, ok := ctx.Value("rootHandler").(*dgcommand.CommandGroup)
	if !ok { return }

	ctx.Reply(GetHelp(handler, paths))
}