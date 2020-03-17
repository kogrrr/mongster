package auth

import (
	"fmt"

	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Endpoint", func() {
	It("adds required routes", func() {
		a := &Auth{}
		router := mux.NewRouter()
		a.AddRoutes(router)

		var login *mux.Route
		var callback *mux.Route

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tmpl, err := route.GetPathTemplate()
			Expect(err).NotTo(HaveOccurred())
			switch tmpl {
			case "/auth/login":
				fmt.Println("found login")
				login = route
			case "/auth/callback":
				fmt.Println("found callback")
				callback = route
			}
			return nil
		})

		Expect(login).NotTo(BeNil())
		Expect(callback).NotTo(BeNil())
	})
})
