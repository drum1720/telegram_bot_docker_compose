package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"strings"
)

type Db struct {
	Db *gorm.DB
}

type Word struct {
	Word   string
	Status string
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
	}

	if !db.Db.Migrator().HasTable(&Word{}) {
		db.Db.Migrator().CreateTable(&Word{})
		go createAndFillWordsTable()
		log.Println("create table Words")
	}

	return err
}

func createAndFillWordsTable() {
	addWordsToTable("russian.txt")
	log.Println("ya vse!")
	addWordsToTable("russian1.txt")
	log.Println("ya vse!")
	addWordsToTable("russian2.txt")
	log.Println("ya vse!")
	addWordsToTable("russian3.txt")
	log.Println("ya vse!")
	addWordsToTable("russian4.txt")
	log.Println("ya vse!")
	addWordsToTable("russian5.txt")
	log.Println("ya vse!")
}

func addWordsToTable(fullNameFileTxt string) {
	var db Db
	err := db.connect()
	if err != nil {
		return
	}

	buffFile, err := ioutil.ReadFile(fullNameFileTxt)
	if err != nil {
		log.Println(err)
	}

	text := (string)(buffFile)
	words := strings.Split(text, "\n")
	for i := 0; i < len(words); i++ {
		word := Word{
			Word:   words[i],
			Status: "New"}

		db.Db.Create(&word)
	}
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

func (word *Word) UpdateStatus(db Db) error {
	result := db.Db.Model(&Word{}).Where("word = ?", word.Word).Update("status", word.Status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (word *Word) WhereOneResponse(db Db) error {
	db.Db.Where("status = ?", "New").Find(&word)
	if db.Db.Error != nil {
		return db.Db.Error
	}
	return nil
}
