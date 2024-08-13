package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"go_gin_example/auth"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/joho/godotenv/autoload"
)

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	UpdatedAt   string `json:"updated_at"`
}

type LonginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	app.POST("/login", loginHandler)
	app.POST("/tasks", CreateHandler)
	app.GET("/tasks", ReadsHandler)
	app.GET("/tasks/:id", ReadHandler)
	app.PUT("/tasks/:id", UpdateHandler)
	app.DELETE("/tasks/:id", DeleteHandler)
	app.Run()
}

func loginHandler(c *gin.Context) {
	var loginInput LonginInput
	var err = c.BindJSON(&loginInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Malformed JSON when login in",
		})
		return
	}

	token, err := auth.EncodeJWT(loginInput.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err,
		})
		return
	}

	data := make(map[string]any)
	data["token"] = token

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
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
	})
}

func ReadsHandler(c *gin.Context) {
	authorizationHeader := c.Request.Header["Authorization"]

	if len(authorizationHeader) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Authorization header is missing",
		})
		return
	}

	claims, err := auth.DecodeJWT(authorizationHeader[0])
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": err.Error(),
		})
		return
	}
	fmt.Println(claims.(jwt.MapClaims)[username])

	db := mySqlConnect()
	defer db.Close()

	rows, err := db.Query("SELECT `id`, `name`, `description`, `create_at`, `updated_at` FROM `task` ORDER BY `create_at` DESC")
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
	var id = c.Param("id")

	db := mySqlConnect()
	defer db.Close()

	row := db.QueryRow("SELECT `id`, `name`, `description`, `create_at`, `updated_at` FROM `task` WHERE `id` = ?", id)

	var task Task
	err := row.Scan(&task.Id, &task.Name, &task.Description, &task.CreateAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusNotFound,
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": task,
	})
}

func UpdateHandler(c *gin.Context) {
	var id = c.Param("id")

	db := mySqlConnect()
	defer db.Close()

	var updatedTask Task

	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when updating the Task",
		})
		return
	}

	stmt, err := db.Prepare("UPDATE `task` SET name=?, description=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedTask.Name, updatedTask.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Task Updated!",
	})
}

func DeleteHandler(c *gin.Context) {
	var id = c.Param("id")

	db := mySqlConnect()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM `task` WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete task",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Task Delete!",
	})
}
