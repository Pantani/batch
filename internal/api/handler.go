package api

import (
	"net/http"
	"time"

	"github.com/Pantani/batch/internal/db/database"
	"github.com/Pantani/batch/internal/model"

	"github.com/gin-gonic/gin"
)

// @Summary Get Batch
// @ID batch
// @Description Get the pending current batch or by id
// @Accept json
// @Produce json
// @Tags Transactions
// @Param id query int false "the batch id" default(0)
// @Success 200 {object} model.Batch
// @Failure 500 {object} errResponse
// @Router /batch [get]
func getBatch(router gin.IRouter, db database.IDatabase) {
	router.GET("/batch", func(c *gin.Context) {
		id := c.Query("id")
		var b model.Batch
		var err error
		if len(id) > 0 {
			b, err = db.GetBatch(id)
		} else {
			b, err = db.GetPendingBatch()
		}
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				errorResponse(err),
			)
			return
		}
		c.JSON(http.StatusOK, &b)
	})
}

// @Summary Add Transaction
// @ID tx
// @Description Add a new transactions to batch
// @Accept json
// @Produce json
// @Tags Transactions
// @Param transaction body model.Transaction true "The transaction details"
// @Success 200 {object} model.Batch
// @Failure 500 {object} errResponse
// @Router /transaction [post]
func addTransaction(router gin.IRouter, db database.IDatabase, minValue int, duration time.Duration) {
	router.POST("/transaction", func(c *gin.Context) {
		var tx model.Transaction
		if err := c.Bind(&tx); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		tx.CreatedAt = time.Now()

		ok, err := db.IsEmpty()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// create a first batch.
		b := model.FirstBatch(minValue, duration)
		if !ok {
			// get the last pending batch.
			b, err = db.GetPendingBatch()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}

		// add the new transaction to the batch and save it.
		b.Transactions = append(b.Transactions, tx)
		err = db.SaveBatch(b)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// check the batch send availability.
		if b.CheckAvailability() {
			// TODO send transaction

			// create a new batch.
			newBatch := model.NewBatch(b.ID, minValue, duration)
			err = db.SaveBatch(newBatch)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}

		c.JSON(http.StatusOK, &b)
	})
}
