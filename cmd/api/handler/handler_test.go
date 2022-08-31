// build+integration
package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var API_URL = "http://localhost:3000"
var _ = Describe("Handler", Label("integration"), func() {
	It("Should return OK 200", func() {
		req, _ := http.NewRequest(http.MethodGet, API_URL+"/healthcheck", nil)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		res, _ := io.ReadAll(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(string(res)).To(Equal("\"OK\""))
	})
})
