package crud

// User Create response from Users microservice
type CreateUserResponse struct {
	Error   bool   `json:"error"`
	Code    string `json:"code"`
	Message string `json:"msg"`
}
