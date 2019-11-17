package vm


import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//ReqUser req model
type ReqUser struct{
	Name string `form:"name" json:"name"`
	Token string `form:"token" json:"token"`
	Repo string `form:"repo" json:"repo"`
}

type CommitInfo struct{
	
}


// MakeSuccess return success response
func MakeSuccess(c *gin.Context, code int, content interface{}) {
	c.JSON(http.StatusOK, gin.H{"statusCode": code, "data": content})
}

// MakeFail return fail response
func MakeFail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{"statusCode": code, "message": message})
}
