package auth

import (
	"errors"
	"net/http"
	"web/pkg/common/exceptions"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h handler) Login(c *gin.Context) {
	body := LoginRequestBody{}

	// getting request's body
	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]exceptions.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = exceptions.ApiError{Param: fe.Field(), Message: exceptions.MsgForTag(fe)}
			}

			c.JSON(http.StatusBadRequest, exceptions.BadValidation(out))
			return
		}
	}

	usr, err := h.service.LoginUser(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.BadRequest("Invalid credentials"))
		return
	}

	token, err := h.service.GenerateToken(usr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("error while authenticating"))
		return
	}

	session := sessions.Default(c)
	session.Set("token", token)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("error while authenticating"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
