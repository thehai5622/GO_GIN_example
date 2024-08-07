package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var tasks []Task = []Task{
	{Name: "Task Name", Description: "Task Description"},
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
		"message": "Hello world",
	})
}

func CreateHandler(c *gin.Context) {
	var newTask Task
	var error = c.BindJSON(&newTask)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed JSON when creating a Task",
		})
		return
	}

	tasks = append(tasks, newTask)
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    newTask,
	})
}

func ReadsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":   tasks,
		"length": len(tasks),
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
