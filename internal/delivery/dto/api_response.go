package dto

type ApiResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type UserDetailResponse struct {
	ID    uint64 `json:"id"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
