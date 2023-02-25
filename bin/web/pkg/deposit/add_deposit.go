package deposit

import (
	"core/deposit"
	"errors"
	"lib/common/exceptions"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AddDepositRequestBody struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h handler) AddDeposit(c *gin.Context) {
	body := AddDepositRequestBody{}

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

	session := sessions.Default(c)

	usr, err := h.authService.ExtractUser(session.Get("token"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, exceptions.Unauthorized("unauthorized"))
		return
	}

	err = h.service.AddDeposit(c.Request.Context(), usr, deposit.AddDepositPayload{
		Amount: body.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("error processing deposit"))
		return
	}

	c.JSON(204, nil)
}
