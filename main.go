package main

/*
 *  Main file for clinancial
 *  Copyright (C) 2017 Arthur Mendes
 */

import (
	"fmt"
	"os"
	"time"
)

type CCommandFunc func([]string)

type CCommand struct {
	name     string
	desc     string
	function CCommandFunc
}

var commands = make([]CCommand, 0)

func printHelp() {
	fmt.Println(" clinancial - a command-line financial manager")
	fmt.Println("")
	fmt.Println(" Commands: ")

	for _, c := range commands {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc)
	}
}

func main() {

	commands = append(commands,
		CCommand{name: "help", desc: "Print this help text",
			function: _printHelp},
		CCommand{name: "account", desc: "Manages accounts",
			function: manageAccounts},
		CCommand{name: "argprint", desc: "Test argument printing",
			function: testArgs})

	// Check command
	if len(os.Args) <= 1 {
		printHelp()
		return
	}

	fmt.Println(" Please note that the interface might be not fully functional")
	for _, c := range commands {
		if c.name == os.Args[1] {
			c.function(os.Args[1:])
			return
		}
	}

	fmt.Println("No command named " + os.Args[1])
}

func _printHelp(args []string) {
	printHelp()
}

func testArgs(args []string) {
	fmt.Println("")
	for _, s := range args {
		fmt.Print(s, " - ")
	}

	fmt.Println("")
}

func manageAccounts(args []string) {
	if len(args) < 3 {
		fmt.Println("Expected format: " + args[0] + " [create|view|delete] account_name")
		return
	}

	operation := args[1]
	acc_name := args[2]

	if operation == "create" {
		a := &Account{id: uint(time.Now().Unix()),
			name: acc_name}
		a.Create()
		fmt.Printf("Account %s created (id %d)\n",
			a.GetName(), a.GetID())

		fmt.Println(" -- note that it will not persist... --- ")
		return
	}

	if operation == "view" {
		fmt.Println("view operation is not supported")
		return
	}

	if operation == "delete" {
		fmt.Println("delete operation is not supported")
	}
}
