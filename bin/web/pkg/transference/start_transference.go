package transference

import (
	"core/transference"
	"errors"
	"lib/common/exceptions"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type StartTransferRequestBody struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	ToUser  string  `json:"to_user" binding:"required,uuid4"`
	Message *string `json:"message" binding:""`
}

func (h handler) StartTransference(c *gin.Context) {
	body := StartTransferRequestBody{}

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

	err = h.service.StartTransfer(c.Request.Context(), usr, transference.StartTransferencePayload{
		Amount:  body.Amount,
		ToUser:  uuid.MustParse(body.ToUser), // validated but idk
		Message: body.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError("error processing transference"))
		return
	}

	c.JSON(204, nil)
}
