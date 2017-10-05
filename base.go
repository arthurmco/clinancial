package main

/*
 *  Types for basic financial transactions and things
 *
 *  Copyright (C) 2017 Arthur M
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

	/* Create an account in a database */
	Create()

	/* Get actual value for that month/year */
	GetValue(month, year uint) float32

	/* Add a register to an account */
	AddRegister(f *FinancialRegister)

	/* Remove a register from an account */
	RemoveRegister(f* FinancialRegister)

	/* Get register from an account */
	GetRegisterbyID(id uint) *FinancialRegister
	GetRegistersbyDatePeriod(start, end time.Time) []*FinancialRegister
};

/*
 *   A financial register. 
 *   Contains information about a single transaction
 */
type FinancialRegister struct {
	id uint
	name string
	time time.Time
	value float32
	from BaseAccount
	to BaseAccount
};


