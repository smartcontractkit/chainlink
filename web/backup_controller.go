package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
)

// BackupController streams backups over GET.
type BackupController struct {
	App *services.ChainlinkApplication
}

// Show streams a backup of the current db through a read-only transaction.
func (bc *BackupController) Show(c *gin.Context) {
	tx, err := bc.App.GetStore().GetBolt().Begin(false)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer tx.Rollback()

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/octet-stream"}
	header["Content-Disposition"] = []string{"attachment; filename=backup.bolt"}
	header["Content-Length"] = []string{strconv.Itoa(int(tx.Size()))}

	_, err = tx.WriteTo(c.Writer)
	if err != nil {
		c.AbortWithError(500, err)
	}
}
