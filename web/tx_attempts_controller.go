package web

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// TxAttemptsController lists TxAttempts requests.
type TxAttemptsController struct {
	App services.Application
}

// Index returns paginated transaction attempts
func (tac *TxAttemptsController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	var attempts []models.TxAttempt
	query, count, err := getTxAttempts(tac.App.GetStore(), size, offset)
	if err != nil {
		c.AbortWithError(500, err)
	} else if err := query.Find(&attempts); err == orm.ErrorNotFound {
		c.Data(404, MediaType, emptyJSON)
	} else if err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting paged TxAttempts: %+v", err))
	} else if buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, attempts); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, MediaType, buffer)
	}
}

func getTxAttempts(store *store.Store, size int, offset int) (storm.Query, int, error) {
	count, err := store.Count(&models.TxAttempt{})
	query := store.Select().OrderBy("SentAt").Reverse().Limit(size).Skip(offset)
	return query, count, err
}
