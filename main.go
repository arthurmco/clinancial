package main

/*
 *  Main file for clinancial
 *  Copyright (C) 2017 Arthur Mendes
 */

import (
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
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
	fmt.Printf(" Usage: %s [command] [commandargs...]\n", os.Args[0])
	fmt.Println("")
	fmt.Println(" Commands: ")

	for _, c := range commands {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc)
	}
}

func main() {
	SetDatabasePath(os.Getenv("HOME") + "/.config/clinancial.db")

	/* Use the enviroment variable to set the db path, if present */
	if os.Getenv("CLINANCIAL_DB") != "" {
		SetDatabasePath(os.Getenv("CLINANCIAL_DB"))
	}

	commands = append(commands,
		CCommand{name: "help", desc: "Print this help text",
			function: _printHelp},
		CCommand{name: "account", desc: "Manages accounts",
			function: manageAccounts},
		CCommand{name: "register",
			desc: "Manages financial registers, i.e transactions",
			function: manageRegisters},
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


func manageRegisters(args []string) {
	if len(args) < 2 {
		fmt.Println("Expected format: " + args[0] + " [create|view]")
		return
	}

	operation := args[1]
	
	if operation == "create" {
		var acname string
		var acval float32
		var acfrom, acto *Account
		accready := false

		accounts, aerr := GetAllAccounts()
		if aerr != nil {
			panic(aerr)
		}

		if len(accounts) == 0 {
			panic("No account created. \n"+
				"Please type "+os.Args[0]+" account create <acc> to create an account (named <acc>)")

			
		}

		accstrlist := make([]string, 0)
		for _, aval := range accounts {
			astr := fmt.Sprintf("%d: %s",
				aval.GetID(), aval.GetName())

			accstrlist = append(accstrlist, astr)
		}

		for !accready {
			// Request register name
			fmt.Print("Name: ")
			rd := bufio.NewReader(os.Stdin)
			
			acname, err := rd.ReadString('\n')
			acname = acname[0:len(acname)-1]
			var num int = 0

			if err != nil {
				panic(err)
			}

			// Request value
			valready := false
			for !valready {
				fmt.Print("Value ($): ")
				num, err = fmt.Scanf("%f", &acval)

				if num <= 0 {
					fmt.Fprint(os.Stderr,
						"Invalid price format")
				}

				if err != nil {
					panic(err)
				}
				valready = true
			}

			// Request src account
			var fromacc int
			fromready := false
			for !fromready {
				fmt.Println("Choose the origin account (the one to be debited)")
				fmt.Println("Available ones: " + strings.Join(accstrlist, ", "))
				fmt.Print("Number: ")
				num, err = fmt.Scanf("%d", &fromacc)

				for _, acc := range accounts {
					if acc.GetID() == uint(fromacc) {
						acfrom = acc
						fromready = true
						break
					}
				}

				if !fromready {
					fmt.Fprintf(os.Stderr,
						"This account does not exist")
					
				}
			}


			// Request dest account
			var toacc int
			toready := false
			for !toready {
				fmt.Println("Choose the destiny account (the one to be credited)")
				fmt.Println("Available ones: " + strings.Join(accstrlist, ", "))
				fmt.Print("Number: ")
				num, err = fmt.Scanf("%d", &toacc)

				for _, acc := range accounts {
					if acc.GetID() == uint(toacc) {
						acto = acc
						toready = true
						break
					}
				}

				if !toready {
					fmt.Fprintf(os.Stderr,
						"This account does not exist")
					
				}
			}


			strfrom := acfrom.GetName()
			strto := acto.GetName()
			fmt.Printf("Creating register '%s' with value %.2f, from account %s to account %s"+
				"\n\tConfirm (Y/N) or Ctrl+C to exit\n", acname, acval, strfrom, strto)

			res := "N"
			fmt.Scanf("%s", &res)

				
			if res == "Y" || res == "y" {
				accready = true
			}
		}
		
		freg := &FinancialRegister{name: acname, value: acval,
			from: acfrom, to: acto, time: time.Now()}
		err := acfrom.AddRegister(freg)

		if err != nil {
			panic(err)
		}
		
		
	}

}


func manageAccounts(args []string) {
	if len(args) < 2 {
		fmt.Println("Expected format: " + args[0] + " [create|view|delete]")
		return
	}

	operation := args[1]

	if operation == "create" {
		if len(args) < 3 {
			fmt.Println("Expected format: " + args[0] + " create <account_name>")
			return
		}

		acc_name := strings.TrimSpace(args[2])
		a := &Account{id: uint(time.Now().Unix()),
			name: acc_name}
		a.Create()
		fmt.Printf("Account %s created (id %d)\n",
			a.GetName(), a.GetID())

		return
	}

	if operation == "view" {
		acc, err := GetAllAccounts()
		if err != nil {
			fmt.Print("fatal: ")
			panic(err)
		}

		if len(acc) == 0 {
			fmt.Println("\t\tNo accounts registered")
			return
		}

		fmt.Printf("          id        |  price  | creation date \n")
		fmt.Printf("====================|=========|===============\n")
		
		tm := time.Now().Month()
		ty := time.Now().Year()
		for _, val := range acc {

			price, _ := val.GetValue(uint(tm), uint(ty))
			datefmt := val.GetCreationDate().Format("2006-01-02")
			
			fmt.Printf(" %-18s | %7.2f | %s\n", val.GetName(),
				price, datefmt)

		}

		fmt.Println("")
		return
	}

	if operation == "delete" {
		fmt.Println("delete operation is not supported")
	}
}
