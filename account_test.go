package main

/*
 *  Tests for the account type methods
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"strconv"
	"testing"
	"time"
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

func TestCreateAndGetRegister(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	a.AddRegister(&FinancialRegister{id: 1, name: "Test", time: time.Now(),
		value: 50, from: b, to: a})
	a.AddRegister(&FinancialRegister{id: 2, name: "Test", time: time.Now(),
		value: 30, from: a, to: b})

	r := a.GetRegisterbyID(1)

	if r == nil {
		t.Error("wrong value, got nil, should be 1")
		DropDatabase()
		return
	}

	if r.id != 1 {
		t.Error("wrong value, got " + strconv.Itoa(int(r.id)) + ", should be 1")
	}

	DropDatabase()
}

func TestCreateAndRemoveRegister(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	a.AddRegister(&FinancialRegister{id: 1, name: "Test", time: time.Now(),
		value: 50, from: b, to: a})
	a.AddRegister(&FinancialRegister{id: 2, name: "Test", time: time.Now(),
		value: 30, from: a, to: b})

	r := a.GetRegisterbyID(1)
	if r == nil {
		t.Error("wrong value, got nil, should be 1")
		DropDatabase()
		return
	}

	a.RemoveRegister(r)
	r = a.GetRegisterbyID(1)
	if r != nil {
		t.Error("wrong value, should be nil ")
	}

	DropDatabase()
}

func TestGetPrice(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	a.AddRegister(&FinancialRegister{id: 1, name: "Test", time: time.Now(),
		value: 50, from: b, to: a})
	a.AddRegister(&FinancialRegister{id: 2, name: "Test", time: time.Now(),
		value: 30, from: a, to: b})
	a.AddRegister(&FinancialRegister{id: 3, name: "Test", time: time.Now(),
		value: 130, from: b, to: a})

	tm := time.Now().Month()
	ty := time.Now().Year()

	price := a.GetValue(uint(tm), uint(ty))
	if price != float32(150) {
		t.Error("wrong value, got " + strconv.FormatFloat(
			float64(price), 'f', -1, 64) + ", should be 150.0")
	}

	DropDatabase()
}

func TestGetRegisterByDate(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	a.AddRegister(&FinancialRegister{id: 18, name: "Test",
		time:  time.Date(2000, 10, 1, 0, 0, 0, 0, time.Now().Location()),
		value: 50, from: b, to: a})
	a.AddRegister(&FinancialRegister{id: 2, name: "Test", time: time.Now(),
		value: 30, from: a, to: b})
	a.AddRegister(&FinancialRegister{id: 3, name: "Test", time: time.Now(),
		value: 130, from: b, to: a})

	regs := a.GetRegistersbyDatePeriod(
		time.Date(2000, 9, 20, 0, 0, 0, 0, time.Now().Location()),
		time.Date(2000, 10, 20, 0, 0, 0, 0, time.Now().Location()))

	if len(regs) != 1 {
		t.Error("wrong len, got " + strconv.Itoa(
			len(regs)) + ", should be 1")
		DropDatabase()
		return
	}

	if regs[0].id != 18 {
		t.Error("wrong len, got " + strconv.Itoa(
			int(regs[0].id)) + ", should be 18")
	}

	DropDatabase()
}
