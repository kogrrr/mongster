package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"

	"github.com/gargath/mongster/pkg/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Middleware", func() {
	Context("when accessing API", func() {

		var server *httptest.Server
		var client *http.Client

		BeforeEach(func() {
			router := mux.NewRouter()
			api := router.PathPrefix("/api").Subrouter()
			api.Use(auth.TokenVerifierMiddleware)
			api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "{\n\"message\": \"Test\"\n}\n")
			})
			client = &http.Client{}
			server = httptest.NewServer(router)
		})

		AfterEach(func() {
			server.Close()
		})

		It("rejects invalid tokens with 403", func() {
			req, err := http.NewRequest("GET", server.URL+"/api/", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Authorization", "Bearer fuddleduck")
			res, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(403))
		})

		It("rejects requests missing tokens with 401", func() {
			resp, err := client.Get(server.URL + "/api/")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(401))
		})
	})
})
