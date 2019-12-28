package command

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Command struct {
	Name     string
	Callback CommandCallback
}

type CommandContext struct {
	InvokedCommand *Command
	Message *discordgo.MessageCreate
	Args    []string
	Bot *Bot
}

type CommandCallback func(ctx *CommandContext) string

func parseCommand(cmd string) []string {
	i := 0
	var args []string
	start := 0
	end := 0
	for i < len(cmd) {
		char := cmd[i]
		if char == '"' {
			i++
			end = strings.Index( cmd[i:],`"`)
			if end == -1 {
				end = len(cmd) - 1
			} else {
				end = i+end
			}
			args = append(args, cmd[i:end])
			i = end+1
			start = i+1
		} else if char == ' ' {
			args = append(args, cmd[start:i])
			start = i+1
		}
		i++
	}
	args = append(args, cmd[start:i])
	return args
}