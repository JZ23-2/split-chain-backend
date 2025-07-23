package dtos

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type FailedResponse struct {
	Error string `json:"error"`
}
