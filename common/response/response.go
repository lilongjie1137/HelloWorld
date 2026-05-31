// Package response 定义统一的 HTTP 响应包装。
//
// 约定：所有接口返回 {"code":int,"msg":string,"data":any}，
// code=0 表示成功，非 0 为业务错误码（见 common/errcode）。
package response

import (
	"errors"
	"net/http"

	"github.com/lilongjie1137/HelloWorld/common/errcode"
)

// Body 统一响应体。
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 构造成功响应。
func Success(data any) Body {
	return Body{Code: errcode.OK.Code, Msg: errcode.OK.Message, Data: data}
}

// Fail 由 error 构造失败响应；非业务错误统一归为内部错误。
func Fail(err error) Body {
	var be *errcode.Error
	if errors.As(err, &be) {
		return Body{Code: be.Code, Msg: be.Message}
	}
	return Body{Code: errcode.ErrInternal.Code, Msg: errcode.ErrInternal.Message}
}

// HTTPStatus 将业务错误码映射到合适的 HTTP 状态码（便于网关/监控）。
func HTTPStatus(code int) int {
	switch code {
	case errcode.OK.Code:
		return http.StatusOK
	case errcode.ErrInvalidParam.Code, errcode.ErrSpecRequired.Code, errcode.ErrModifierInvalid.Code:
		return http.StatusBadRequest
	case errcode.ErrUnauthorized.Code:
		return http.StatusUnauthorized
	case errcode.ErrForbidden.Code, errcode.ErrTenantScope.Code, errcode.ErrStoreScope.Code, errcode.ErrDiscountDenied.Code:
		return http.StatusForbidden
	case errcode.ErrNotFound.Code, errcode.ErrSpuNotFound.Code, errcode.ErrOrderNotFound.Code, errcode.ErrBillNotFound.Code:
		return http.StatusNotFound
	case errcode.ErrConflict.Code, errcode.ErrPrintJobDup.Code:
		return http.StatusConflict
	case errcode.ErrTooManyReq.Code:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
