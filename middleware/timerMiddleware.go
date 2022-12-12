package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// type TimerMiddleWareHandler struct {
// }

// func (t *TimerMiddleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Response) {

// }

func TimerMiddleWareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; elapsed time %s \n", time.Since(start))
	})
}
