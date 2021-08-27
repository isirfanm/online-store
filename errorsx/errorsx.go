package errorsx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Error interface
type Error struct {
	RespCode int    `json:"-"`
	Message  string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorResponse handle error
func ErrorResponse(c *gin.Context, err error) {
	if err == nil {
		return
	}

	log.Error(err)

	e, ok := err.(*Error)
	if !ok {
		c.AbortWithStatusJSON(
			500,
			gin.H{
				"message": err.Error(),
			})
		return
	}

	c.AbortWithStatusJSON(e.RespCode, e)
}

// NewGeneralErrorMsg create error
func NewGeneralErrorMsg(msg string) *Error {
	return &Error{
		RespCode: 500,
		Message:  msg,
	}
}
