package main

/*
 *  Main file for clinancial
 *  Copyright (C) 2017 Arthur M
 */



import (
	"fmt"
	"os"
)

type CCommandFunc func([]string)

type CCommand struct {
	name string
	desc string
	function CCommandFunc	
}

var commands = make([]CCommand, 0);

func printHelp() {
	fmt.Println(" clinancial - a command-line financial manager");
	fmt.Println("");
	fmt.Println(" Commands: ");

	for _, c := range(commands) {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc);
	}
}

func main() {

	commands = append(commands,
		CCommand{name: "help", desc: "Print this help text",
			function: _printHelp},
		CCommand{name: "argprint", desc: "Test argument printing",
			function: testArgs});
		

	// Check command
	if len(os.Args) <= 1 {
		printHelp();
		return;
	}

	for _, c := range(commands) {
		if c.name == os.Args[1] {
			c.function(os.Args[2:]);
			return;
		}
	}

	fmt.Println("No command named "+ os.Args[1]);
}

func _printHelp(args []string) {
	printHelp();
}

func testArgs(args []string) {
	fmt.Println("");
	for _, s := range args {
		fmt.Print(s, " - ");
	}

	fmt.Println("");
}
