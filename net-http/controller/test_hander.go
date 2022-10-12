package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
)

type Handler struct {
	mu sync.Mutex
	n  int
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	_, err := fmt.Fprintf(w, "count is %d\n", h.n)
	if err != nil {
		glog.Errorf("write response failed, err:%v", err)
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func NewFoo() HandlerFunc {
	return Foo1
}

func Foo1(w http.ResponseWriter, r *http.Request) {
	obj := map[string]interface{}{
		"data": nil,
	}

	w.Header().Set("Content-Type", "application/json")

	foo := r.FormValue("foo")
	if foo == "" {
		obj["data"] = "foo is empty"
	} else {
		obj["data"] = foo
	}

	byt, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(byt)

	return
}
