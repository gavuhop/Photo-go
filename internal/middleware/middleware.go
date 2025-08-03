package middleware

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"photo-go/internal"
// 	"photo-go/internal/database"
// 	"photo-go/pkg/logger"
// 	"runtime"
// 	"strings"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// func LogRequest(c fiber.Ctx, traceID string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			logger.Error(errors.New("panic in LogRequest"), "TraceID: %s - Panic in LogRequest: %+v", traceID, r)
// 		}
// 	}()

// 	// Log request body parameters
// 	if c == nil {
// 		logger.Error(errors.New("context is nil"), "Context is nil")
// 		return
// 	}

// 	body := "-empty-"
// 	if c.Body() != nil {
// 		body = string(c.Body())
// 	}

// 	headers := "-empty-"
// 	if c.GetReqHeaders() != nil {
// 		headers = fmt.Sprintf("%v", c.GetReqHeaders())
// 	}
// 	method := "-empty-"
// 	if c.Method() != "" {
// 		method = c.Method()
// 	}
// 	fullPath := "-empty-"
// 	if c.OriginalURL() != "" {
// 		fullPath = c.OriginalURL()
// 	}
// 	logger.Info("Request details - fullPath: %s, method: %s, headers: %s, body: %s - TraceID: %s",
// 		fullPath, method, headers, body, traceID)
// }

// // Error types for consistent error handling across the application

// // Transaction creates a database transaction middleware that automatically
// // handles transaction lifecycle (begin, commit, rollback) based on request success/failure.
// func Transaction(db *gorm.DB) fiber.Handler {
// 	return func(c fiber.Ctx) (err error) {
// 		// Recover from panic (if any) and return 500 error
// 		traceID := uuid.New().String()
// 		defer func() {
// 			if r := recover(); r != nil {
// 				// Log full stack trace
// 				buf := make([]byte, 1024)
// 				runtime.Stack(buf, true)
// 				logger.Error(errors.New("panic in transaction"), "Panic in transaction: %+v\n%s - TraceID: %s", r, string(buf), traceID)
// 				// Return 500 error without rollback
// 				_ = fiber.NewError(fiber.StatusInternalServerError, internal.StatusMessageInternalError)
// 			}
// 		}()
// 		// Attach traceID to context
// 		if c != nil {
// 			c.Locals("traceID", traceID)
// 		}
// 		// Log request details
// 		LogRequest(c, traceID)

// 		// Check if database is nil
// 		if db == nil {
// 			logger.Error(errors.New("database connection not available"), "Database connection not available")
// 			return fiber.NewError(fiber.StatusInternalServerError, internal.StatusMessageDatabaseConnectionNotAvailable)
// 		}

// 		// Start a transaction
// 		tx := db.Begin()
// 		if tx.Error != nil {
// 			logger.Error(tx.Error, "Could not start transaction")
// 			return fiber.NewError(fiber.StatusInternalServerError, internal.StatusMessageCouldNotStartTransaction)
// 		}

// 		// Attach tx to context
// 		c.Locals("tx", tx)

// 		// Execute next handlers
// 		if err = c.Next(); err != nil {
// 			logger.Error(err, "Unhandled error.")
// 			// Rollback transaction if there is an error
// 			database.CloseDB(tx, false, traceID)
// 			// Only override with 500 if it's not a Fiber error
// 			if fiberErr, ok := err.(*fiber.Error); ok {
// 				return fiberErr
// 			}
// 			err = fiber.NewError(fiber.StatusInternalServerError, internal.StatusMessageInternalError)
// 			return err
// 		}

// 		apiPath := c.Path()
// 		if strings.Contains(apiPath, "stream") {
// 			// Skip transaction commit for streaming
// 			// It will be committed in the service layer
// 			return nil
// 		}

// 		// Commit otherwise
// 		// Only log if there is an error and do not return error
// 		database.CloseDB(tx, true, traceID)
// 		return nil
// 	}
// }

// // ErrorHandler provides centralized error handling and response formatting
// // for all application errors.
// func ErrorHandler() fiber.Handler {
// 	return func(c fiber.Ctx) error {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				buf := make([]byte, 1024)
// 				runtime.Stack(buf, true)
// 				logger.Error(errors.New("panic in error handler"), "Panic in error handler: %+v\n%s", r, string(buf))
// 				// Return 500 error without rollback
// 				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 					Code:    internal.StatusCodeInternalError,
// 					Message: "Internal server error (panic)",
// 				})
// 				database.CloseDBFiber(&c, false)

// 			}
// 		}()
// 		// Execute next handlers
// 		err := c.Next()
// 		if err == nil {
// 			return nil
// 		}

// 		// Handle specific error types with appropriate HTTP status codes
// 		return handleError(c, err)
// 	}
// }

// // handleError processes different error types and returns appropriate JSON responses
// func handleError(c fiber.Ctx, err error) error {
// 	// Recover from panic and return 500 error
// 	defer func() {
// 		if r := recover(); r != nil {
// 			// Log full stack trace
// 			buf := make([]byte, 1024)
// 			runtime.Stack(buf, true)
// 			logger.Error(errors.New("panic in handleError"), "Panic in handleError: %+v\n%s", r, string(buf))
// 			// Set response status and JSON
// 			_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				Code:    internal.StatusCodeInternalError,
// 				Message: "Internal server error (panic)",
// 			})
// 			database.CloseDBFiber(&c, false)
// 		}
// 	}()

// 	// Handle GORM-specific errors
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			Code:    internal.StatusCodeNotFound,
// 			Message: internal.StatusMessageNotFound,
// 		})
// 	}

// 	// Handle context-related errors
// 	if errors.Is(err, context.DeadlineExceeded) {
// 		return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
// 			Code:    internal.StatusCodeRequestTimeout,
// 			Message: internal.StatusMessageRequestTimeout,
// 		})
// 	}

// 	if errors.Is(err, context.Canceled) {
// 		return c.Status(fiber.StatusGone).JSON(fiber.Map{
// 			Code:    internal.StatusCodeRequestCanceled,
// 			Message: internal.StatusMessageRequestCanceled,
// 		})
// 	}

// 	// Handle Fiber errors with custom status codes
// 	if fiberErr, ok := err.(*fiber.Error); ok {
// 		return handleFiberError(c, fiberErr)
// 	}

// 	// Default error response for unknown errors
// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 		Code:    internal.StatusCodeInternalError,
// 		Message: internal.StatusMessageInternalError,
// 	})
// }

// // handleFiberError processes Fiber-specific errors and maps them to appropriate responses
// func handleFiberError(c fiber.Ctx, fiberErr *fiber.Error) error {
// 	statusCode := fiberErr.Code
// 	errorCode := internal.StatusCodeInternalError
// 	message := internal.StatusMessageInternalError

// 	// Map known error messages to appropriate codes
// 	switch fiberErr.Message {
// 	case internal.StatusMessageDatabaseConnectionNotAvailable:
// 		errorCode = internal.StatusCodeInternalError
// 		message = internal.StatusMessageDatabaseConnectionNotAvailable
// 	case internal.StatusMessageCouldNotStartTransaction:
// 		errorCode = internal.StatusCodeInternalError
// 		message = internal.StatusMessageCouldNotStartTransaction
// 	case internal.StatusMessageInternalError:
// 		errorCode = internal.StatusCodeInternalError
// 		message = internal.StatusMessageInternalError
// 	default:
// 		// Use the original Fiber error message for unknown errors
// 		message = fiberErr.Message
// 	}

// 	return c.Status(statusCode).JSON(fiber.Map{
// 		Code:    errorCode,
// 		Message: message,
// 	})
// }
