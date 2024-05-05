package repository

const (

	// Error codes
	DATABASE_ERROR         string = "connection_database_error"
	CACHE_ERROR            string = "connection_cache_error"
	USER_NOT_FOUND         string = "auth_user_not_found"
	INVALID_CREDENTIALS    string = "auth_invalid_credentials"
	INVALID_DATA           string = "auth_invalid_form_data"
	ERROR_LOGGING_OUT      string = "auth_could_not_log_out"
	ERROR_RETREIVING_TOKEN string = "auth_invalid_token"
	ERROR_VERIFYING_TOKEN  string = "auth_invalid_token"
	MICROSERVICE_ERROR     string = "microservice_error"
)
