package main

/*
 *  Types for basic financial transactions and things
 *
 *  Copyright (C) 2017 Arthur Mendes
 */

import (
	"time"
)

/*
 *   Base account
 *   Contains information about money storage
 *   An account is where you do transactions. You always send/receive money
 *   to/from an account, regardless of what account is (a bank account, the bank
 *   or a business)
 */
type BaseAccount interface {
	/* Account ID */
	GetID() uint

	/* Account Name */
	GetName() string
	SetName(s string)

	/* Get creation date */
	GetCreationDate() time.Time

	/* Create an account in a database */
	Create() error

	/* Update account info in the database */
	Update() error

	/* Get accounts in the database */
	GetbyID(id uint) error
	GetbyName(name string) error

	/* Get actual value for that month/year */
	GetValue(month, year uint) (float32, error)

	/* Add a register to an account */
	AddRegister(f *FinancialRegister) error

	/* Remove a register from an account */
	RemoveRegister(f *FinancialRegister) error

	/* Get register from an account */
	GetRegisterbyID(id uint) (*FinancialRegister, error)
	GetRegistersbyDatePeriod(start, end time.Time) ([]*FinancialRegister, error)
}

type AccountError struct {
	err  string
	code int
}

func (a *AccountError) Error() string {
	return a.err
}

/*
 *   A financial register.
 *   Contains information about a single transaction
 */
type FinancialRegister struct {
	id    uint
	name  string
	time  time.Time
	value float32
	from  BaseAccount
	to    BaseAccount
}
