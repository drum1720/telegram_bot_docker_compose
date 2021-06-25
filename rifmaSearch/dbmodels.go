package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Db struct {
	Db *gorm.DB
}

type Rifma struct {
	Request string `gorm:"primaryKey; unique"`
	Rifma   string
}

func (db *Db) connect() error {
	var err error
	var settings Settings
	settings.updateData()
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", settings.PgHost, settings.PgUser, settings.PgPass, settings.PgDbName, settings.PgPort)
	db.Db, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}
	if !db.Db.Migrator().HasTable(&Rifma{}) {
		db.Db.Migrator().CreateTable(&Rifma{})
		log.Println("create table")
	} else {
		log.Println("table is done")
	}
	return err
}

func (rifma Rifma) AddToTable(db Db) error {
	result := db.Db.Create(&rifma)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (rifma *Rifma) WhereOneResponse(db Db) error {
	db.Db.Where("request = ?", rifma.Request).Find(&rifma)
	if db.Db.Error != nil {
		return db.Db.Error
	}
	return nil
}
