package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"resto-app-server/internal/repo"
)

func Init(repository *repo.Repo) *gin.Engine {
	r := gin.Default()

	r.GET("/restaurants/:uuid", func(c *gin.Context) {
		uuidStr := c.Param("uuid")
		parsedUuid, err := uuid.Parse(uuidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurant, err := repository.GetOneRestaurant(&parsedUuid)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, restaurant)
	})

	r.POST("/restaurants/create", func(c *gin.Context) {
		newRestaurant := repo.Restaurant{}

		err := c.ShouldBind(&newRestaurant)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		newUUID, err := repository.CreateNewRestaurant(&newRestaurant)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, newUUID)
	})

	r.POST("/tables/create", func(c *gin.Context) {
		newTable := repo.Table{}

		err := c.ShouldBind(&newTable)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		newUUID, err := repository.CreateNewTable(&newTable)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, newUUID)
	})

	r.POST("/reservations/create", func(c *gin.Context) {
		newReservation := repo.Reservation{}

		err := c.ShouldBind(&newReservation)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		newUUID, err := repository.CreateNewReservation(&newReservation)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, newUUID)
	})

	return r
}
