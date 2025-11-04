package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
)

// ErrorResponse padroniza erros simples.
type ErrorResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Error     string `json:"error"`
	Code      int    `json:"code"`
}

// ValidationDetail descreve um erro de validação de campo.
type ValidationDetail struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

// ValidationErrorResponse padroniza 422 com detalhes.
type ValidationErrorResponse struct {
	RequestID string             `json:"request_id,omitempty"`
	Error     string             `json:"error"`
	Code      int                `json:"code"`
	Details   []ValidationDetail `json:"details,omitempty"`
}

// HandleError transforma erros de domínio em HTTP.
func HandleError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := err.Error()

	switch err {
	case domain.ErrNotFound:
		code = http.StatusNotFound
		msg = "registro não encontrado"
	case domain.ErrInvalidInput:
		code = http.StatusBadRequest
		msg = "entrada inválida"
	case domain.ErrAlreadyExists:
		code = http.StatusConflict
		msg = "registro já existente"
	case domain.ErrPrecondition:
		code = http.StatusPreconditionFailed
		msg = "pré-condição não atendida"
	case domain.ErrUnauthorized:
		code = http.StatusUnauthorized
		msg = "não autorizado"
	case domain.ErrForbidden:
		code = http.StatusForbidden
		msg = "acesso negado"
	}

	reqID, _ := c.Get("request_id")
	rid, _ := reqID.(string)

	c.JSON(code, ErrorResponse{
		RequestID: rid,
		Error:     msg,
		Code:      code,
	})
	c.Abort()
}

// ValidationError retorna 422 e, quando possível, detalhes por campo/regra.
func ValidationError(c *gin.Context, err error) {
	reqID, _ := c.Get("request_id")
	rid, _ := reqID.(string)

	details := make([]ValidationDetail, 0, 4)

	// Se o erro for do validator.v10, extraímos os campos.
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			details = append(details, ValidationDetail{
				Field: fe.Field(),
				Rule:  fe.Tag(),
			})
		}
	}

	c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
		RequestID: rid,
		Error:     "erro de validação",
		Code:      http.StatusUnprocessableEntity,
		Details:   details,
	})
	c.Abort()
}
