package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	UpdatedAt   string `json:"updated_at"`
}

var tasks []Task = []Task{
	{Name: "Task Name", Description: "Task Description"},
}

var (
	dbname   = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func mySqlConnect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10)
	return db
}

func main() {
	app := gin.Default()
	app.GET("/", helloHandler)
	app.POST("/tasks", CreateHandler)
	app.GET("/tasks", ReadsHandler)
	app.GET("/tasks/:id", ReadHandler)
	app.PUT("/tasks/:id", UpdateHandler)
	app.DELETE("/tasks/:id", DeleteHandler)
	app.Run()
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"app_env": os.Getenv("APP_ENV"),
	})
}

func CreateHandler(c *gin.Context) {
	db := mySqlConnect()
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO `task`(`name`, `description`) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	var newTask Task
	var error = c.BindJSON(&newTask)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when creating a Task",
		})
		return
	}

	_, err = stmtIns.Exec(newTask.Name, newTask.Description)
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OK",
		"data":    newTask,
	})
}

func ReadsHandler(c *gin.Context) {
	db := mySqlConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description, create_at, updated_at FROM `task`")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Internal Server Error",
		})
		return
	}
	defer rows.Close()

	columnsName, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columnsName))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var data []Task

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Initialize a new Task object
		var task Task
		for i, col := range values {
			switch columnsName[i] {
			case "id":
				task.Id = string(col)
			case "name":
				task.Name = string(col)
			case "description":
				task.Description = string(col)
			case "create_at":
				task.CreateAt = string(col)
			case "updated_at":
				task.UpdatedAt = string(col)
			}
		}

		data = append(data, task)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"data":   data,
		"length": len(data),
	})
}

func ReadHandler(c *gin.Context) {
	var id, error = strconv.Atoi(c.Param("id"))

	if error != nil {
		fmt.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid ID",
		})
		return
	}

	for i, v := range tasks {
		if i == id {
			c.JSON(http.StatusOK, gin.H{
				"data": v,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Task ID not Found",
	})
}

func UpdateHandler(c *gin.Context) {
	var id, error1 = strconv.Atoi(c.Param("id"))

	if error1 != nil {
		fmt.Println(error1)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid ID",
		})
		return
	}

	var oldTask = tasks[id]
	var newTask Task
	var error2 = c.BindJSON(&newTask)
	if error2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON",
		})
		return
	}
	oldTask.Name = newTask.Name
	oldTask.Description = newTask.Description
	tasks[id] = newTask

	c.JSON(http.StatusOK, gin.H{
		"message": "Task Updated",
		"data":    oldTask,
	})
}

func DeleteHandler(c *gin.Context) {
	var id, error1 = strconv.Atoi(c.Param("id"))

	if error1 != nil {
		fmt.Println(error1)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid ID",
		})
		return
	}

	firstHalf := tasks[:id]
	secondHalf := tasks[id+1:]
	tasks = append(firstHalf, secondHalf...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Task Delete",
	})
}
