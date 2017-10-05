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

var dbpath string = "clinancial.db"

func SetDatabasePath(s string) {
	dbpath = s
}

func CreateDatabase() error {
	db, eerr := sql.Open("sqlite3", dbpath)
	if eerr != nil {
		return eerr
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS accounts (" +
		"id INTEGER PRIMARY KEY, name TEXT, ctime INTEGER)")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS registers (" +
		"id INTEGER PRIMARY KEY, sid INTEGER, name string, " +
		"time INTEGER, val REAL, fromaccount INTEGER, " +
		"toaccount INTEGER) ")
	if err != nil {
		return err
	}
	stmt.Exec()
	db.Close()
	return nil
}

func DropDatabase() error {
	db, eerr := sql.Open("sqlite3", dbpath)
	if eerr != nil {
		return eerr
	}

	stmt, err := db.Prepare("DROP TABLE IF EXISTS accounts")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = db.Prepare("DROP TABLE IF EXISTS registers")
	if err != nil {
		return err
	}
	stmt.Exec()
	db.Close()
	return nil
}

/* Add account in the database */
func (a *Account) Create() error {
	a.transactions = make(map[uint][]*FinancialRegister)
	a.creationDate = time.Now()

	err := CreateDatabase()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbpath)
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

	db, err := sql.Open("sqlite3", dbpath)
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

	db, err := sql.Open("sqlite3", dbpath)
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

	db, err := sql.Open("sqlite3", dbpath)
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

	db, err := sql.Open("sqlite3", dbpath)
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
