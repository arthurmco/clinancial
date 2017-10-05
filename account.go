package main

/*
 *  Defines most common account types
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"time"
)

/*
 *  A concrete account
 */

type Account struct {
	// ID
	id uint

	// Name
	name string

	// ctime
	creationDate time.Time

	// Transaction list
	// The transaction list is encoded in a map. The key is the  month and
	// year, in this way: 201701 to access data for jan 2017, etc.
	transactions map[uint][]*FinancialRegister
}

func (a *Account) GetID() uint {
	return a.id
}

func (a *Account) GetName() string {
	return a.name
}

func (a *Account) SetName(s string) {
	a.name = s
}

func (a *Account) GetValue(month, year uint) float32 {
	if month >= 12 {
		month = 1
		year++
	} else {
		month++
	}

	stopkey := (year * 100) + month

	// Get all registers until the one specified.
	var val float32 = 0.0
	for tkey, tlist := range a.transactions {
		if tkey >= stopkey {
			break
		}

		for _, freg := range tlist {
			v := float32(0.0)
			if a.id == freg.from.GetID() {
				v = -freg.value
			} else if a.id == freg.to.GetID() {
				v = freg.value
			}

			val += v
		}
	}

	return val
}

func (a *Account) AddRegister(f *FinancialRegister) {
	tkey := uint(f.time.Year()*100) + uint(f.time.Month())

	a.transactions[tkey] = append(a.transactions[tkey], f)
}

func (a *Account) RemoveRegister(f *FinancialRegister) {
	tkey := uint(f.time.Year()*100) + uint(f.time.Month())

	for idx, freg := range a.transactions[tkey] {
		if freg.id == f.id {
			a.transactions[tkey] =
				append(a.transactions[tkey][:idx],
					a.transactions[tkey][(idx+1):]...)
			break
		}
	}
}

func (a *Account) GetRegisterbyID(id uint) *FinancialRegister {
	for _, tlist := range a.transactions {

		for _, freg := range tlist {
			if freg.id == id {
				return freg
			}
		}
	}

	return nil
}

func (a *Account) GetRegistersbyDatePeriod(start, end time.Time) []*FinancialRegister {
	skey := uint(start.Year()*100) + uint(start.Month())
	ekey := uint(end.Year()*100) + uint(end.Month())

	ret := make([]*FinancialRegister, 0)

	for idx, tlist := range a.transactions {
		if idx < skey {
			continue
		}

		if idx > ekey {
			break
		}

		for _, freg := range tlist {
			if (freg.time.Unix() > start.Unix()) &&
				(freg.time.Unix() < end.Unix()) {
				ret = append(ret, freg)
			}

		}
	}

	return ret
}

func (a *Account) Create() {
	a.transactions = make(map[uint][]*FinancialRegister)
	a.creationDate = time.Now()
}
