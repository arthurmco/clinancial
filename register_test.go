package main

/*
 *  Tests for the register methods
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"strconv"
	"testing"
	"time"
)

func TestCreateAndGetRegister(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	err := a.AddRegister(&FinancialRegister{id: 1, name: "Test",
		time: time.Now(), value: 50, from: b, to: a})
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	err = a.AddRegister(&FinancialRegister{id: 2, name: "Test",
		time: time.Now(), value: 30, from: a, to: b})
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	r, eerr := a.GetRegisterbyID(1)
	if eerr != nil {
		t.Error(eerr)
		DropDatabase()
		return
	}

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

	err := a.AddRegister(&FinancialRegister{id: 1, name: "Test",
		time: time.Now(), value: 50, from: b, to: a})
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	err = a.AddRegister(&FinancialRegister{id: 2, name: "Test",
		time: time.Now(), value: 30, from: a, to: b})
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	r, eerr := a.GetRegisterbyID(1)
	if eerr != nil {
		t.Error(eerr)
		DropDatabase()
		return
	}

	if r == nil {
		t.Error("wrong value, got nil, should be 1")
		DropDatabase()
		return
	}

	err = a.RemoveRegister(r)
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	r, err = a.GetRegisterbyID(1)
	if err == nil {
		t.Error(err)
		DropDatabase()
		return
	}

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

	price, err := a.GetValue(uint(tm), uint(ty))
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	if price != float32(150) {
		t.Error("wrong value, got " + strconv.FormatFloat(
			float64(price), 'f', -1, 64) + ", should be 150.0")
	}

	DropDatabase()
}

func TestGetRegisterByDate(t *testing.T) {
	a := createTestAccount(1)
	b := createTestAccount(2)

	a.AddRegister(&FinancialRegister{id: 1, name: "Test", time: time.Now(),
		value: 30, from: a, to: b})
	a.AddRegister(&FinancialRegister{id: 2, name: "Test",
		time:  time.Date(2000, 10, 1, 0, 0, 0, 0, time.Now().Location()),
		value: 50, from: b, to: a})
	a.AddRegister(&FinancialRegister{id: 3, name: "Test", time: time.Now(),
		value: 130, from: b, to: a})

	regs, err := a.GetRegistersbyDatePeriod(
		time.Date(2000, 9, 20, 0, 0, 0, 0, time.Now().Location()),
		time.Date(2000, 10, 20, 0, 0, 0, 0, time.Now().Location()))
	if err != nil {
		t.Error(err)
		DropDatabase()
		return
	}

	if len(regs) != 1 {
		t.Error("wrong len, got " + strconv.Itoa(
			len(regs)) + ", should be 1")
		DropDatabase()
		return
	}

	if regs[0].id != 2 {
		t.Error("wrong id, got " + strconv.Itoa(
			int(regs[0].id)) + ", should be 2")
	}

	DropDatabase()
}
