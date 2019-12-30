package main

import (
	"fmt"
	"github.com/plally/discord_modular_bot/command"
	"strings"
)

func ExampleFirstWord() {
	fmt.Println(command.GetFirstWord("this is a test"))
	// Output:
	// this
}

func ExampleParseCommand() {
	args := command.ParseCommand(`>ping test`)
	fmt.Println(strings.Join(args, ", "))

	args = command.ParseCommand(`cmd --flag1 something --flag2 "test nothing"`)

	fmt.Println(strings.Join(args, ", "))


	// Output:
	//>ping, test
	//cmd, --flag1, something, --flag2, test nothing
}
