package add

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/Racuwcka/service-courier/internal/http-server/handlers/types/transport"
	"github.com/Racuwcka/service-courier/internal/lib/api/response"
)

type Provider interface {
	Create(ctx context.Context, name string) (uint32, error)
}

type Handler struct {
	provider Provider
}

func New(p Provider) *Handler {
	return &Handler{
		provider: p,
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := transport.AddRequest{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, response.Error("bad validate"))
			return
		}

		courierId, err := h.provider.Create(context.Background(), req.Name)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("create courier failed"))
			return
		}

		responseOK(w, r, courierId)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, courierId uint32) {
	render.JSON(w, r, transport.AddResponse{
		Response:  response.OK(),
		CourierId: courierId,
	})
}
