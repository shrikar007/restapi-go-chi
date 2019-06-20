package dberror
import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`
	Errortext         string `json:"errortext"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(er error) render.Renderer {
	return &ErrorResponse{
		Errortext:           er.Error(),
		HTTPStatusCode:     422,
	}
}
