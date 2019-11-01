package command

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
)

//TODO refactor command package to allow disabling and enabling modules per guild
//types
type Bot struct {
	*discordgo.Session
	EnabledModules []*Module
}

type TextCommand struct {
	Name     string
	Callback commandCallback
}

type Module struct {
	Name     string
	Commands []TextCommand
	OnEnable func(*Bot)
}

type TextCommandEvent struct {
	Command *TextCommand
	Message *discordgo.MessageCreate
	Args    string
}

type commandCallback func(*discordgo.Session, *TextCommandEvent) string

//
var Modules = make(map[string]*Module)

// bot methods
func NewBot(s *discordgo.Session) *Bot {
	bot := &Bot{
		s,
		make([]*Module, 0),
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		bot.CheckCommands(s, m)
	})

	return bot

}

func (b *Bot) EnableModule(name string) {
	m, ok := Modules[name]
	if !ok {
		log.Warnf("Module %s doesn't exist", name)
	}
	m.OnEnable(b)
	b.EnabledModules = append(b.EnabledModules, m)
}

func (b *Bot) DisableModule(name string) {
	length := len(b.EnabledModules)
	for i := 1; i < length; i++ {
		if b.EnabledModules[i].Name != name {
			continue
		}
		b.EnabledModules[i] = b.EnabledModules[length-1]
		b.EnabledModules[length-1] = nil
		return
	}
	log.Warnf("Module %s doesn't exist", name)
}

func (b *Bot) CheckCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandName := getFirstWord(m.Content)
	for _, module := range b.EnabledModules {
		if !isModuleEnabledInGuild(m.GuildID) {
			continue
		}
		for _, cmd := range module.Commands {
			if commandName != cmd.Name {
				continue
			}
			log.Infof("Executing command %s", cmd.Name)
			event := TextCommandEvent{
				Args:    m.Content[len(cmd.Name)+1:],
				Command: &cmd,
				Message: m,
			}
			reply := cmd.Callback(s, &event)
			if reply != "" {
				s.ChannelMessageSend(m.ChannelID, reply)
			}
		}
	}
}

// Module methods
func RegisterModule(name string) *Module {
	if _, ok := Modules[name]; ok {
		log.Errorf("Module %s already exists", name)
		return nil
	}
	m := &Module{
		name,
		make([]TextCommand, 0),
		func(*Bot) {},
	}
	Modules[name] = m
	log.Infof("Registered module %s", m.Name)
	return m
}

func (m *Module) RegisterCommandFunc(name string, callback commandCallback) TextCommand {
	cmd := TextCommand{
		Name:     name,
		Callback: callback,
	}
	return m.RegisterCommand(cmd)
}

func (m *Module) RegisterCommand(c TextCommand) TextCommand {
	log.Infof("Registered command %s in module %s", c.Name, m.Name)
	m.Commands = append(m.Commands, c)
	return c
}

//util
func isModuleEnabledInGuild(guildID string) bool {
	//TODO finish this function
	return true
}

func getFirstWord(s string) string {
	i := strings.IndexByte(s, " "[0])
	if i == -1 {
		return s
	} else {
		return s[:i]
	}
}
