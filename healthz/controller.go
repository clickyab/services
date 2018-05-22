package healthz

import (
	"context"
	"net/http"
	"reflect"

	"time"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
	"github.com/sirupsen/logrus"
)

type route struct {
}

type healthRes struct {
	Time     string                 `json:"time"`
	Errors   map[string]string      `json:"errors"`
	Messages map[string]interface{} `json:"messages"`
}

func (r route) check(ctx context.Context, w http.ResponseWriter, rq *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	var (
		messages = make(map[string]interface{}, len(all))
		errs     = make(map[string]string, len(all))
	)

	for i := range all {
		msg, err := all[i].Healthy(ctx)
		name := reflect.TypeOf(all[i]).Elem().Name()

		if err != nil {
			logrus.Error(err)
			errs[name] = err.Error()
		}

		if msg != "" {
			messages[name] = msg
		}
	}

	code := http.StatusOK
	if len(errs) > 0 {
		code = http.StatusInternalServerError
	}

	w.Header().Set("time", time.Now().String())
	framework.JSON(
		w,
		code,
		healthRes{
			Time:     time.Now().String(),
			Errors:   errs,
			Messages: messages,
		},
	)

	return
}

func (r route) Routes(mux framework.Mux) {
	mux.RootMux().GET("/healthz", xhandler.HandlerFuncC(r.check))
}

func init() {
	router.Register(&route{})
}
