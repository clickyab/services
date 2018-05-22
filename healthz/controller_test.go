package healthz

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"context"
	"net/http"
	"sync"

	"net/http/httptest"

	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
)

type mysqlHealth struct {
	Err error
	Msg string
}

type redisHealth struct {
	Err error
	Msg string
}

type brokerHealth struct {
	Err error
	Msg string
}

func (h *mysqlHealth) Healthy(context.Context) (interface{}, error) {
	return h.Msg, h.Err
}

func (h *redisHealth) Healthy(context.Context) (interface{}, error) {
	return h.Msg, h.Err
}

func (h *brokerHealth) Healthy(context.Context) (interface{}, error) {
	return h.Msg, h.Err
}

type MyHandler struct {
	sync.Mutex
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	defer h.Unlock()
	rout := route{}
	ctx := context.Background()
	rout.check(ctx, w, r)

}

func TestHealth(t *testing.T) {
	// hRes := make(map[string]map[string]interface{}, 2)
	hRes := healthRes{}

	handler := &MyHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()
	m := &mysqlHealth{
		Err: nil,
		Msg: "",
	}
	r := &redisHealth{
		Err: nil,
		Msg: "",
	}
	b := &brokerHealth{
		Err: nil,
		Msg: "",
	}
	Register(m, r, b)
	convey.Convey("test with all ok", t, func() {

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		//logrus.Fatal(m.Err,r.Err,b.Err)
		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusOK)
	})

	convey.Convey("mysql error", t, func() {
		m.Err = errors.New("mysql error here")
		r.Msg = "redis message here"

		logrus.Warn(m.Err)
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(msg, &hRes); err != nil {
			t.Fatal(err)
		}

		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusInternalServerError)
		convey.So(fmt.Sprint(hRes.Errors["mysqlHealth"]), convey.ShouldEqual, "mysql error here")
		convey.So(fmt.Sprint(hRes.Messages["redisHealth"]), convey.ShouldEqual, "redis message here")
	})

	convey.Convey("mysql and redis error", t, func() {
		m.Err = errors.New("mysql error here")
		r.Err = errors.New("redis error here")
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(msg, &hRes); err != nil {
			t.Fatal(err)
		}

		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusInternalServerError)
		convey.So(fmt.Sprint(hRes.Errors["mysqlHealth"]), convey.ShouldEqual, "mysql error here")
		convey.So(fmt.Sprint(hRes.Errors["redisHealth"]), convey.ShouldEqual, "redis error here")
	})

}
