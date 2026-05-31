package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/lilongjie1137/HelloWorld/common/errcode"
)

func TestSuccess(t *testing.T) {
	b := Success(map[string]int{"n": 1})
	if b.Code != 0 || b.Msg != "ok" || b.Data == nil {
		t.Errorf("unexpected success body: %+v", b)
	}
}

func TestFail_BusinessError(t *testing.T) {
	b := Fail(errcode.ErrSpuNotFound)
	if b.Code != errcode.ErrSpuNotFound.Code || b.Msg != errcode.ErrSpuNotFound.Message {
		t.Errorf("unexpected fail body: %+v", b)
	}
}

func TestFail_GenericError(t *testing.T) {
	b := Fail(errors.New("boom"))
	if b.Code != errcode.ErrInternal.Code {
		t.Errorf("generic error should map to internal, got %d", b.Code)
	}
}

func TestHTTPStatus(t *testing.T) {
	cases := map[int]int{
		errcode.OK.Code:               http.StatusOK,
		errcode.ErrInvalidParam.Code:  http.StatusBadRequest,
		errcode.ErrUnauthorized.Code:  http.StatusUnauthorized,
		errcode.ErrForbidden.Code:     http.StatusForbidden,
		errcode.ErrOrderNotFound.Code: http.StatusNotFound,
		errcode.ErrPrintJobDup.Code:   http.StatusConflict,
		errcode.ErrInternal.Code:      http.StatusInternalServerError,
	}
	for code, want := range cases {
		if got := HTTPStatus(code); got != want {
			t.Errorf("HTTPStatus(%d) = %d, want %d", code, got, want)
		}
	}
}
