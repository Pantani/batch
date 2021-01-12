package api

import (
	"time"

	_ "github.com/Pantani/batch/docs" // swagger docs.
	"github.com/Pantani/batch/internal/db/database"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// CreateAPI create the API provider with the application routes.
// It returns the gin engine.
func CreateAPI(db database.IDatabase, mode string, minValue int, duration time.Duration) *gin.Engine {
	engine := initEngine(mode)
	setupBatchAPI(engine, db, minValue, duration)
	setupSwaggerAPI(engine)
	return engine
}

// setupBatchAPI setup the batch API routes.
func setupBatchAPI(router gin.IRouter, db database.IDatabase, minValue int, duration time.Duration) {
	getBatch(router, db)
	addTransaction(router, db, minValue, duration)
}

// setupSwaggerAPI setup the swagger API routes.
func setupSwaggerAPI(router gin.IRouter) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// initEngine init the gin API engine.
// It returns the gin engine.
func initEngine(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	return engine
}
