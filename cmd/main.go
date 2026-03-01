package main

import (
	"fmt"
	"net/http"
	"strconv"

	"zakup/internal/models"

	"github.com/gin-gonic/gin"
)

var nextID = 1

func newApp(product, dept string, amount float64) *models.Application {
	app := models.NewApplication(product, dept, amount)
	app.ID = nextID
	nextID++
	return app
}

var applications = []*models.Application{
	newApp("notebook", "management department", 20),
	newApp("laptop", "development department", 6),
	newApp("pens", "management department", 20),
}

func getApplications(c *gin.Context) { //вывести все заявки
	c.JSON(http.StatusOK, applications)
}
func postApplications(c *gin.Context) {
	//тут должна быть валидация данных или передача в сервис валидации
	c.IndentedJSON(http.StatusCreated, "New application was added")
}
func getApplicationStatus(c *gin.Context) { //тут нужно вывести статус заявки согласно id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	for _, app := range applications {
		if app.ID == id {
			fmt.Println("Статус заявки №", id, " ", app.Status)
			c.JSON(http.StatusOK, gin.H{"status": app.Status})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
}
func deleteApplication(c *gin.Context) {
	//удалить заявку по id
}
func updateApplication(c *gin.Context) {
	//изменить заявку по id. Тут можно добавить права доступа
}

func main() {
	router := gin.Default()

	router.GET("/applications", getApplications)
	router.POST("/applications", postApplications)
	router.GET("/applications/:id/status", getApplicationStatus)
	router.DELETE("/applications/:id", deleteApplication)
	router.PUT("/applications/:id", updateApplication)

	router.Run(":8080")
}
