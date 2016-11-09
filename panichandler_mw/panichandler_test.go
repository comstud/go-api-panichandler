package panichandler_mw

import (
	"context"
	"net/http"
	"testing"

	"github.com/tilteng/go-api-router/api_router"
)

func dummyRoute() *api_router.Route {
	cont := api_router.NewMuxRouter()
	return cont.GET("/", func(ctx context.Context) {})
}

func TestPanicString(t *testing.T) {
	ctx := api_router.NewContextForRequest(nil, &http.Request{}, dummyRoute())
	rctx := api_router.RequestContextFromContext(ctx)

	error_string := "Oops, something bad happened"

	var got_rctx *api_router.RequestContext
	var got_val interface{}

	handler := PanicHandlerFn(func(ctx context.Context, val interface{}) {
		got_rctx = api_router.RequestContextFromContext(ctx)
		got_val = val
	})

	mw := NewMiddleware(handler)
	wrapper := mw.NewWrapper().Wrap(func(ctx context.Context) {
		panic(error_string)
	})

	wrapper(rctx.Context())

	if got_rctx != rctx {
		t.Errorf("panic handler not same: %+v", got_rctx)
	}

	if got_val.(string) != error_string {
		t.Errorf("panic handler received wrong value: %+v", got_val)
	}
}

func TestPanicObject(t *testing.T) {
	type error_obj_type struct{}

	ctx := api_router.NewContextForRequest(nil, &http.Request{}, dummyRoute())
	rctx := api_router.RequestContextFromContext(ctx)

	error_obj := &error_obj_type{}

	var got_rctx *api_router.RequestContext
	var got_val interface{}

	handler := PanicHandlerFn(func(ctx context.Context, val interface{}) {
		got_rctx = api_router.RequestContextFromContext(ctx)
		got_val = val
	})

	mw := NewMiddleware(handler)
	wrapper := mw.NewWrapper().Wrap(func(ctx context.Context) {
		panic(error_obj)
	})

	wrapper(rctx.Context())

	if got_rctx != rctx {
		t.Errorf("panic handler not same: %+v", got_rctx)
	}

	if got_val.(*error_obj_type) != error_obj {
		t.Errorf("panic handler received wrong value: %+v", got_val)
	}
}
