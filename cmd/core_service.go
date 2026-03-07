package main

import (
	"net/http"
	"strconv"

	"zakup/internal/models"
	"zakup/validation_service"

	"github.com/gin-gonic/gin"
)

//создать таблицы
//подключить бд
//структура юзера
//фио, тип продукта, почта

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

func getApplicationStatus(c *gin.Context) { //тут нужно вывести статус заявки согласно id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	for _, app := range applications {
		if app.ID == id {
			//fmt.Println("Статус заявки №", id, " ", app.Status)
			c.JSON(http.StatusOK, gin.H{"status": app.Status})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
}

func updateApplication(c *gin.Context) {
	//изменить заявку по id. Тут можно добавить права доступа
}

func deleteApplication(c *gin.Context) {
	//удалить заявку по id
}

func postApplications(c *gin.Context) { //должна быть валидация
	// 1) Считываем JSON в input (из validation-пакета)
	var in validation_service.CreateApplicationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// 2) Валидируем бизнес-правила
	if err := validation_service.ValidateCreateApplication(in); err != nil {
		// хотим красиво отдать список ошибок полей
		if ve, ok := err.(validation_service.ValidationError); ok {
			c.JSON(http.StatusBadRequest, ve)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3) Создаём доменную модель
	app := models.NewApplication(in.ProductName, in.Department, in.Amount)
	app.ID = nextID
	nextID++

	applications = append(applications, app)

	c.JSON(http.StatusCreated, app)
}

func main() {
	router := gin.Default()

	router.GET("/applications", getApplications)                 //вывести все заявки
	router.POST("/applications", postApplications)               //отправить заявку
	router.GET("/applications/:id/status", getApplicationStatus) //получить статус заявки
	router.DELETE("/applications/:id", deleteApplication)        //удалить заявку
	router.PUT("/applications/:id", updateApplication)           //обновить заявку

	router.Run(":8080")
}
