package internal

var (
	StatusCodeInternalError   = "INTERNAL_SERVER_ERROR"
	StatusCodeRequestCanceled = "REQUEST_CANCELED"
	StatusCodeRequestTimeout  = "REQUEST_TIMEOUT"
	StatusCodeNotFound        = "NOT_FOUND"

	StatusMessageInternalError                  = "Internal server error"
	StatusMessageFailedToCommitTransaction      = "Failed to commit transaction"
	StatusMessageDatabaseConnectionNotAvailable = "Database connection not available"
	StatusMessageDatabaseTransactionNotFound    = "Database transaction not found"
	StatusMessageCouldNotStartTransaction       = "Could not start transaction"
	StatusMessageRequestCanceled                = "Request was canceled"
	StatusMessageRequestTimeout                 = "Request timeout"
	StatusMessageNotFound                       = "Resource not found"
)
