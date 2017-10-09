package main

/*
 *  Defines most common account types
 *  Copyright (C) 2017 Arthur Mendes
 */
import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func (a *Account) GetCreationDate() time.Time {
	return a.creationDate
}

func (a *Account) GetValue(month, year uint) (float32, error) {
	if month >= 12 {
		month = 1
		year++
	} else {
		month++
	}
	tend := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0,
		time.Now().Location())

	vtotal := float32(0.0)
	err := CreateDatabase()
	if err != nil {
		return 0.0, err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return 0.0, err
	}

	res, err := db.Query("SELECT val, fromaccount, toaccount "+
		"FROM registers WHERE time < ?", tend.Unix())

	if err != nil {
		return 0.0, err
	}

	for res.Next() {
		var val float64
		var fromacc, toacc int

		err = res.Scan(&val, &fromacc, &toacc)
		if err != nil {
			return 0.0, err
		}

		var vcredit float32 = 0.0
		if fromacc == int(a.id) {
			vcredit = float32(-val)
		}

		if toacc == int(a.id) {
			vcredit = float32(val)
		}

		vtotal += vcredit
	}

	res.Close()
	db.Close()
	return vtotal, nil
}

func (a *Account) AddRegister(f *FinancialRegister) error {
	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	fromid, toid := 0, 0
	if f.from != nil {
		fromid = int(f.from.GetID())
	}

	if f.to != nil {
		toid = int(f.to.GetID())
	}

	res, err := db.Exec("INSERT INTO registers (name, time, val, "+
		"fromaccount, toaccount) VALUES (?, ?, ?, ?, ?)",
		f.name, f.time.Unix(), f.value, fromid, toid)

	if err != nil {
		return err
	}

	lid, _ := res.LastInsertId()
	f.id = uint(lid)
	db.Close()
	return nil
}

func (a *Account) RemoveRegister(f *FinancialRegister) error {
	if f.id <= 0 {
		return &AccountError{"Invalid financial register ID", 1001}
	}

	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM registers WHERE id = ? AND name = ?",
		f.id, f.name)

	if err != nil {
		return err
	}

	f.id = 0 // invalidate ID
	db.Close()
	return nil
}

func (a *Account) GetRegisterbyID(id uint) (*FinancialRegister, error) {
	err := CreateDatabase()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return nil, err
	}

	res, err := db.Query("SELECT id, name, time, val, fromaccount, "+
		"toaccount FROM registers WHERE id = ? ", id)

	if err != nil {
		return nil, err
	}

	if !res.Next() {
		return nil, &AccountError{"No results", 1000}
	}

	var rid int
	var name string
	var timestamp int64
	var val float64
	var fromaccid, toaccid uint

	err = res.Scan(&rid, &name, &timestamp, &val, &fromaccid, &toaccid)
	if err != nil {
		return nil, err
	}

	res.Close()
	db.Close()
	var fromacc, toacc *Account
	fromacc = &Account{}
	err = fromacc.GetbyID(fromaccid)
	if err != nil {
		fromacc = nil
	}

	toacc = &Account{}
	err = toacc.GetbyID(toaccid)
	if err != nil {
		toacc = nil
	}

	fr := &FinancialRegister{id: uint(rid), name: name,
		time:  time.Unix(timestamp, 0),
		value: float32(val), from: fromacc, to: toacc}

	return fr, nil
}

func (a *Account) GetRegistersbyDatePeriod(start, end time.Time) ([]*FinancialRegister, error) {
	startts := start.Unix()
	endts := end.Unix()

	err := CreateDatabase()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return nil, err
	}

	res, err := db.Query("SELECT id, name, time, val, fromaccount, toaccount "+
		"FROM registers WHERE time > ? AND time < ?", startts, endts)

	if err != nil {
		return nil, err
	}

	registers := make([]*FinancialRegister, 0)

	var id int
	var name string
	var timestamp int64
	var val float64
	var fromaccid, toaccid int

	for res.Next() {
		err = res.Scan(&id, &name, &timestamp, &val, &fromaccid, &toaccid)
		if err != nil {
			return nil, err
		}

		var fromacc, toacc *Account
		fromacc = &Account{}
		err = fromacc.GetbyID(uint(fromaccid))
		if err != nil {
			fromacc = nil
		}

		toacc = &Account{}
		err = toacc.GetbyID(uint(toaccid))
		if err != nil {
			toacc = nil
		}

		registers = append(registers, &FinancialRegister{id: uint(id),
			name: name, time: time.Unix(timestamp, 0),
			value: float32(val), from: fromacc, to: toacc})
	}

	db.Close()
	return registers, nil

}

/* Add account in the database */
func (a *Account) Create() error {
	a.transactions = make(map[uint][]*FinancialRegister)
	a.creationDate = time.Now()

	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	res, err := db.Exec("INSERT INTO accounts (name, ctime) VALUES (?, ?)",
		a.name, a.creationDate.Unix())

	if err != nil {
		panic(err)
	}

	lastid, _ := res.LastInsertId()
	a.id = uint(lastid)
	db.Close()

	return nil
}

/* Update account info in the database */
func (a *Account) Update() error {
	a.transactions = make(map[uint][]*FinancialRegister)
	a.creationDate = time.Now()

	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE accounts SET name = ? WHERE id = ?",
		a.name, a.id)

	if err != nil {
		panic(err)
	}

	db.Close()

	return nil
}

/* Get account in the db by id */
func (a *Account) GetbyID(id uint) error {

	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	res, err := db.Query("SELECT id, name, ctime FROM accounts WHERE id = ?",
		id)
	if err != nil {
		panic(err)
	}

	var sid int
	var sname string
	var sctime int

	if !res.Next() {
		return &AccountError{"No results", 1000}
	}

	err = res.Scan(&sid, &sname, &sctime)
	if err != nil {
		return err
	}

	a.id = uint(sid)
	a.name = sname
	a.creationDate = time.Unix(int64(sctime), 0)

	res.Close()
	db.Close()
	return nil
}

/* Get account in the db by name */
func (a *Account) GetbyName(name string) error {
	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return err
	}

	res, err := db.Query("SELECT id, name, ctime FROM accounts "+
		"WHERE name = ?", name)
	if err != nil {
		panic(err)
	}

	var sid int
	var sname string
	var sctime int

	if !res.Next() {
		return &AccountError{"No results", 1000}
	}
	err = res.Scan(&sid, &sname, &sctime)
	if err != nil {
		return err
	}

	a.id = uint(sid)
	a.name = sname
	a.creationDate = time.Unix(int64(sctime), 0)

	res.Close()
	db.Close()
	return nil
}

func GetAllAccounts() ([]*Account, error) {
	err := CreateDatabase()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", GetDatabasePath())
	if err != nil {
		return nil, err
	}

	res, err := db.Query("SELECT id, name, ctime FROM accounts")
	if err != nil {
		panic(err)
	}

	accounts := make([]*Account, 0)
	for res.Next() {
		var sid int
		var sname string
		var sctime int

		err = res.Scan(&sid, &sname, &sctime)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &Account{id: uint(sid),
			name:         sname,
			creationDate: time.Unix(int64(sctime), 0)})
	}

	return accounts, nil
}
