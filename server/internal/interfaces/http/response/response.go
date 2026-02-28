package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BusinessCode 业务状态码定义
const (
	// 通用成功
	CodeSuccess = 0

	// 通用错误 (1-99)
	CodeInvalidParam    = 1
	CodeInternalError   = 2
	CodeUnauthorized    = 3
	CodeForbidden       = 4
	CodeNotFound        = 5
	CodeMethodNotAllowed  = 6
	CodeTimeout         = 7

	// 认证错误 (100-199)
	CodeAuthFailed          = 100
	CodeTokenInvalid        = 101
	CodeTokenExpired        = 102
	CodeTokenMissing        = 103
	CodeInvalidCredentials  = 104

	// 文章错误 (200-299)
	CodePostNotFound        = 200
	CodePostAlreadyExists   = 201
	CodeSlugExists          = 202
	CodeInvalidTitle        = 203
	CodeInvalidContent      = 204
	CodeInvalidSlug         = 205
	CodeVersionConflict     = 206
	CodeInvalidStatus       = 207
	CodeInvalidTag          = 208

	// BFF 模块错误 (300-399)
	CodeModuleNotFound      = 300
	CodeModuleExecuteError  = 301

	// 文件上传错误 (400-499)
	CodeUploadFailed        = 400
	CodeInvalidFileType     = 401
	CodeFileTooLarge        = 402
	CodeFileNotFound        = 403
)

// CodeMessageMap 错误码映射表
var CodeMessageMap = map[int]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "invalid parameter",
	CodeInternalError:      "internal error",
	CodeUnauthorized:       "unauthorized",
	CodeForbidden:          "forbidden",
	CodeNotFound:           "resource not found",
	CodeMethodNotAllowed:   "method not allowed",
	CodeTimeout:            "request timeout",

	CodeAuthFailed:         "authentication failed",
	CodeTokenInvalid:       "token invalid",
	CodeTokenExpired:       "token expired",
	CodeTokenMissing:       "token missing",
	CodeInvalidCredentials: "invalid username or password",

	CodePostNotFound:       "post not found",
	CodePostAlreadyExists:  "post already exists",
	CodeSlugExists:         "slug already exists",
	CodeInvalidTitle:       "invalid title",
	CodeInvalidContent:     "invalid content",
	CodeInvalidSlug:        "invalid slug",
	CodeVersionConflict:    "version conflict",
	CodeInvalidStatus:      "invalid status",
	CodeInvalidTag:         "invalid tag",

	CodeModuleNotFound:     "module not found",
	CodeModuleExecuteError: "module execute error",

	CodeUploadFailed:       "upload failed",
	CodeInvalidFileType:    "invalid file type",
	CodeFileTooLarge:       "file too large",
	CodeFileNotFound:       "file not found",
}

// GetMessage 获取错误码对应的错误信息
func GetMessage(code int) string {
	if msg, ok := CodeMessageMap[code]; ok {
		return msg
	}
	return "unknown error"
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: GetMessage(code),
	})
}

// ErrorWithMessage 返回自定义错误信息的响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 返回带数据的错误响应
func ErrorWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
	})
}

// Wrapper 包装错误为统一响应
func Wrapper(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
	})
}
