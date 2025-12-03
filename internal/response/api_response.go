package response

import "github.com/gin-gonic/gin"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

// MetaData contains pagination metadata
type MetaData struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PageSize    int   `json:"page_size"`
}

// Success sends a success response
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, err error, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err.Error(),
		Message: message,
	})
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, statusCode int, data interface{}, total int64, limit, offset int) {
	// Calculate current page (1-indexed)
	currentPage := (offset / limit) + 1
	if limit == 0 {
		currentPage = 1
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: MetaData{
			CurrentPage: currentPage,
			TotalPages:  totalPages,
			TotalItems:  total,
			PageSize:    limit,
		},
	})
}
