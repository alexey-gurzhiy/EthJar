package db

import (
	log "EthJar/app/log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type dbSettings struct {
	driver string
	dbLink string
}

var dbSet = dbSettings{
	driver: "sqlite3",
	dbLink: "app/db/jar.db",
}

type Jar struct {
	Jar_id          int
	Jar_password    string
	Jar_name        string
	Jar_description string
	Jar_status      string
}

type Wallet struct {
	Wallet     string
	Jar_id     int
	PrivateKey string
}

func (j Jar) Write() int {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "app/db",
		"\n	2. Function": "Write (j Jar)",
	})
	db, err := sql.Open(dbSet.driver, dbSet.dbLink)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("sql connection OK")
	defer db.Close()

	if j.Jar_id != 0 {
		log.Trace("Updating with jar_id")
		_, err := db.Exec("update jar set jar_password = $2, jar_name = $3, jar_description = $4, jar_status = $5, where jar_id = $1", j.Jar_id, j.Jar_password, j.Jar_name, j.Jar_description, j.Jar_status)
		if err != nil {
			log.Panic(err)
		}
		return 0
	}
	log.Trace("Creating new jar_id")
	createJarToDB, err := db.Exec("insert into jar (jar_password, jar_name, jar_description, jar_status) values ($1, $2, $3, $4)", j.Jar_password, j.Jar_name, j.Jar_description, j.Jar_status)
	if err != nil {
		log.Panic(err)
	}
	jar_id64, err := createJarToDB.LastInsertId()
	if err != nil {
		log.Panic(err)
	}
	log.Info("Data stored, jar ID: ", jar_id64)
	return int(jar_id64)
}

func (w Wallet) Write() {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "app/db",
		"\n	2. Function": "Write (w Wallet)",
	})
	db, err := sql.Open(dbSet.driver, dbSet.dbLink)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("sql connection OK")
	defer db.Close()

	_, err = db.Exec("insert into wallet (wallet, jar_id, privateKey) values ($1, $2, $3)", w.Wallet, w.Jar_id, w.PrivateKey)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("Wallet stored")
}

func (j Jar) Read() Jar {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "app/db",
		"\n	2. Function": "Write (j Jar)",
	})
	db, err := sql.Open(dbSet.driver, dbSet.dbLink)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("sql connection OK")
	defer db.Close()
	row := db.QueryRow("select * from jar where jar_id = $1", j.Jar_id)
	if err != nil {
		log.Panic(err)
	}

	err = row.Scan(&j.Jar_id, &j.Jar_password, &j.Jar_name, j.Jar_description, j.Jar_status)
	if err != nil {
		log.Panic(err)
	}
	log.Info(j)
	

	return j
}

func (w Wallet) Read() Wallet {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "app/db",
		"\n	2. Function": "Read (w Wallet)",
	})
	db, err := sql.Open(dbSet.driver, dbSet.dbLink)
	if err != nil {
		log.Panic(err)
	}
	log.Trace("sql connection OK")
	defer db.Close()

	row := db.QueryRow("select wallet, privateKey from wallet where jar_id = $1", w.Jar_id)
	if err != nil {
		log.Panic(err)
	}

	err = row.Scan(&w.Wallet, &w.PrivateKey)
	if err != nil {
		log.Error(err)
	}

	log.Info(w)
	return w
}
