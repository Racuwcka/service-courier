package transport

import "github.com/Racuwcka/service-courier/internal/lib/api/response"

type AddRequest struct {
	Name string `validate:"required" json:"name"`
}

type AddResponse struct {
	response.Response
	CourierId uint32 `json:"courier_id"`
}
