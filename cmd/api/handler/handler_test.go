//go:build integration
// +build integration

package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var API_URL = "http://localhost:3000"
var _ = Describe("Handler", func() {

	It("Should return OK 200", func() {
		req, _ := http.NewRequest(http.MethodGet, API_URL+"/healthcheck", nil)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		res, _ := ioutil.ReadAll(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(string(res)).To(Equal("\"OK\""))
	})

})
