package main

import (
//	"bytes"
	"database/sql"
//	"fmt"
	"net/http"
//	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
    Fname, Sname, Dname, Email string
	Id int
}

func dbConnect() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "supersecret"
    dbHost := "mysql.go"
	dbPort := "3306"
	dbName := "company"
    db, err := sql.Open(dbDriver, dbUser +":"+ dbPass +"@tcp("+ dbHost +":"+ dbPort +")/"+ dbName +"?charset=utf8")
	checkErr(err)
    return db
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
	router := gin.Default()
	
	// GET all persons
	router.GET("/employee", func(c *gin.Context) {

	employee := Employee{}
	employees := []Employee{}

	db := dbConnect()
	rows, err := db.Query("select * from employees")
	checkErr(err)
	for rows.Next() {
		var first_name, last_name, department, email string
		var id int
		err = rows.Scan(&id, &first_name, &last_name, &department, &email)
		checkErr(err)
		employee.Id = id
		employee.Fname = first_name
		employee.Sname = last_name
		employee.Dname = department
		employee.Email = email
		employees = append(employees, employee)
		
	}
	defer db.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": employees,
			"count":  len(employees),
		})
	})
	
	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {
		var (
			employee Employee
			result gin.H
		)
		id := c.Param("id")
		db := dbConnect()		
		selDB, err := db.Query("SELECT * FROM employees WHERE id=?", id)
		checkErr(err)
		err = selDB.Scan(&employee.Id, &employee.Fname, &employee.Sname, &employee.Dname, &employee.Email)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": employee,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
		defer db.Close()
	})
	router.Run(":3000")	
	
}