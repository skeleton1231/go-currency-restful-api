package item

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Delete deletes an item by its ID.
func (ic *ItemController) Delete(c *gin.Context) {

	itemID, _ := strconv.Atoi(c.Param("itemID"))

	if err := ic.srv.Items().Delete(c, itemID, metav1.DeleteOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
