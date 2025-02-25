package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/pkg/log"
)

// Get get an user by the user identifier.
func (u *UserController) Get(c *gin.Context) {
	log.L(c).Info("get user function called.")

	user, err := u.srv.Users().Get(c, c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
