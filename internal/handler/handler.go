package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
	"resto-app-server/docs"
	"resto-app-server/internal/repo"
)

type Handler struct {
	repo *repo.Repo
}

func New(repo *repo.Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetRestaurant godoc
// @Id GetRestaurant
// @Summary get-restaurant example
// @Description get restaurant full info
// @Tags restaurant
// @Accept json
// @Produce json
// @Param uuid path string true "restaurant uuid"
// @Success 200 {object} repo.RestaurantWithTables
// @Router /restaurants/{uuid} [get]
func (handler *Handler) GetRestaurant(c *gin.Context) {
	uuidStr := c.Param("uuid")
	parsedUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	restaurant, err := handler.repo.GetOneRestaurant(&parsedUuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// CreateRestaurant godoc
// @Id CreateRestaurant
// @Summary create-restaurant example
// @Description create restaurant full info
// @Tags restaurant
// @Accept json
// @Produce json
// @Param restaurant body repo.Restaurant true "restaurant body"
// @Success 200 {string} newUuid
// @Router /restaurants/create [post]
func (handler *Handler) CreateRestaurant(c *gin.Context) {
	newRestaurant := repo.Restaurant{}

	err := c.ShouldBind(&newRestaurant)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newUUID, err := handler.repo.CreateNewRestaurant(&newRestaurant)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, newUUID)
}

// CreateTable godoc
// @Id CreateTable
// @Summary create-table example
// @Description create table full info
// @Tags table
// @Accept json
// @Produce json
// @Param table body repo.Table true "table body"
// @Success 200 {string} newUuid
// @Router /tables/create [post]
func (handler *Handler) CreateTable(c *gin.Context) {
	newTable := repo.Table{}

	err := c.ShouldBind(&newTable)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newUUID, err := handler.repo.CreateNewTable(&newTable)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, newUUID)
}

// CreateReservation godoc
// @Id CreateReservation
// @Summary create-reservation example
// @Description create reservation full info
// @Tags reservation
// @Accept json
// @Produce json
// @Param reservation body repo.Reservation true "reservation body"
// @Success 200 {string} newUuid
// @Router /reservations/create [post]
func (handler *Handler) CreateReservation(c *gin.Context) {
	newReservation := repo.Reservation{}

	err := c.ShouldBind(&newReservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newUUID, err := handler.repo.CreateNewReservation(&newReservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, newUUID)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (handler *Handler) Init() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.Use(CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/restaurants/:uuid", handler.GetRestaurant)
	r.POST("/restaurants/create", handler.CreateRestaurant)
	r.POST("/tables/create", handler.CreateTable)
	r.POST("/reservations/create", handler.CreateReservation)

	return r
}
