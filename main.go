package main

/*
 *  Main file for clinancial
 *  Copyright (C) 2017 Arthur M
 */



import (
	"fmt"
)

func main() {
	f := FinancialRegister{name: "Comprei Pão", value: 38.20};
	
	fmt.Printf("Hello world (%s => %.2f)\n", f.name, f.value);

}
