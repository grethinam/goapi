package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"log"

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
	router.GET("/employ/:id", func(c *gin.Context) {
		var (
			employ Employee
			result gin.H
		)
		id := c.Param("id")
		db := dbConnect()		
		row := db.QueryRow("SELECT id, first_name, last_name, department, email  FROM employees WHERE id=?", id)
		err := row.Scan(&employ.Id, &employ.Fname, &employ.Sname, &employ.Dname, &employ.Email)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": employ,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
		defer db.Close()
	})
	
	// POST new person details
	router.POST("/employ", func(c *gin.Context) {
		db := dbConnect()
		var buffer bytes.Buffer
		fname := c.PostForm("fname")
		sname := c.PostForm("sname")
		dname := c.PostForm("dname")
		email := c.PostForm("email")
		insForm, err := db.Prepare("INSERT INTO employees(first_name, last_name, department, email) VALUES(?,?,?,?);")
	    checkErr(err)
        insForm.Exec(fname, sname, dname, email)
        log.Println("INSERT: First Name: " + fname + " | LAST_NAME: " + sname+ " | DEPARTMENT: " + dname+ " | EMAIL: " + email)

		// Fastest way to append strings
		buffer.WriteString(fname)
		buffer.WriteString(" ")
		buffer.WriteString(sname)
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	defer db.Close()
	})

	// DELETE a person details
	router.DELETE("/employ", func(c *gin.Context) {
	db := dbConnect()
	id := c.Query("id")
	delForm, err := db.Prepare("DELETE FROM employees WHERE id=?;")
	checkErr(err)
    delForm.Exec(emp)
    log.Println("DELETE")
	
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})
    defer db.Close()
	})
	
	router.Run(":3000")	
	
}