package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	config "UserStore/pkg/config"
	"UserStore/pkg/models"

	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const mockEnv = "MOCKING"

func SetupOrDie(dbconfig *config.Dbconfig) *gorm.DB {
	var err error
	var ok bool
	db := connectRoot(dbconfig)
	//Check if Service Database Exist on Postgresql
	ok, err = checkServiceDBExist(db, dbconfig)
	if err != nil {
		fmt.Println(err)
	}
	//Create Database if not exist
	if !ok {
		if err := createServiceDB(db, dbconfig); err != nil {
			capture(err)
		}
	}
	//Connect to Database
	if db, err = ConnectServiceDB(dbconfig); err != nil {
		capture(err)
	}
	//create Table schema for user using gorm
	if err = db.AutoMigrate(&models.User{}); err != nil {
		capture(err)
	}
	return db
}

func connectRoot(config *config.Dbconfig) *gorm.DB {

	if strings.ToLower(os.Getenv(mockEnv)) == "true" { // If we run in a Mocking environement, we do not use a real Database
		return nil
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=postgres password=%s sslmode=disable", config.Uri, config.Port, config.UserName, config.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	capture(err)
	return db
}

func checkServiceDBExist(dbconn *gorm.DB, config *config.Dbconfig) (bool, error) {

	tx := dbconn.Exec(fmt.Sprintf(`SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = '%s';`, strings.ToLower(config.Name)))
	// tx := dbconn.Exec("SELECT datname FROM 'pg_catalog.pg_database';")
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return false, tx.Error
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func createServiceDB(dbconn *gorm.DB, config *config.Dbconfig) *pgconn.PgError {
	db := dbconn.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Name))
	if db.Error != nil {
		fmt.Println("Unable to create DB test_db, attempting to connect assuming it exists...")
		return db.Error.(*pgconn.PgError)
	}
	//defer db.Close()
	fmt.Println("Service DB succesfully created")
	return nil
}

func ConnectServiceDB(config *config.Dbconfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.Uri, config.Port, config.UserName, config.Name, config.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Unable to connect to " + config.Name)
		return nil, err
	}
	return db, nil
}

func capture(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
