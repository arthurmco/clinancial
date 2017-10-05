package main

/*
 *  Tests for the account type methods
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"strconv"
	"testing"
)

func createTestAccount(id uint) *Account {
	SetDatabasePath("/tmp/clinancial.test")
	a := &Account{id: id, name: "Account" + strconv.Itoa(int(id))}
	a.Create()
	return a
}

func TestAccountPersistency(t *testing.T) {
	DropDatabase()
	a := createTestAccount(1)

	aa := &Account{}
	err := aa.GetbyID(1)

	if err != nil {
		DropDatabase()
		t.Fatal(err)
	}

	if aa.name != a.name {
		t.Error("ID: wrong value, got " + aa.GetName() + "|" +
			strconv.Itoa(int(aa.GetID())) +
			", should be " + a.GetName() + "|" +
			strconv.Itoa(int(a.GetID())))
		DropDatabase()
		return
	}

	aa = &Account{}
	err = aa.GetbyName("Account1")

	if err != nil {
		DropDatabase()
		t.Fatal(err)
	}

	if aa.name != a.name {
		t.Error("name: wrong value, got " + aa.GetName() + "|" +
			strconv.Itoa(int(aa.GetID())) +
			", should be " + a.GetName() + "|" +
			strconv.Itoa(int(a.GetID())))
	}

	DropDatabase()
}

func TestAccountAllAccounts(t *testing.T) {
	SetDatabasePath("/tmp/clinancial.test")
	acc, err := GetAllAccounts()
	if err != nil {
		DropDatabase()
		t.Fatal(err)
	}

	if len(acc) != 0 {
		DropDatabase()
		t.Error("empty set: expected 0, found " + strconv.Itoa(len(acc)))
		return
	}

	a := createTestAccount(1)
	b := createTestAccount(2)
	acc, err = GetAllAccounts()
	if err != nil {
		DropDatabase()
		t.Fatal(err)
	}

	if len(acc) != 2 {
		DropDatabase()
		t.Error("full set: expected 2, found " + strconv.Itoa(len(acc)))
		return
	}

	if acc[0].GetID() != a.GetID() {
		DropDatabase()
		t.Error("first item: expected 1, found " +
			strconv.Itoa(int(acc[0].GetID())))
		return

	}

	if acc[1].GetID() != b.GetID() {
		t.Error("second item: expected 2, found " +
			strconv.Itoa(int(acc[1].GetID())))

	}

	DropDatabase()
}
