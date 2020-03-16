package auth

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Util", func() {

	It("generates strings", func() {
		a, err1 := randomHexBytes(1)
		b, err20 := randomHexBytes(20)
		c, err50 := randomHexBytes(50)
		d, err128 := randomHexBytes(128)
		Expect(err1).NotTo(HaveOccurred())
		Expect(err20).NotTo(HaveOccurred())
		Expect(err50).NotTo(HaveOccurred())
		Expect(err128).NotTo(HaveOccurred())
		Expect(len(a)).To(Equal(2))
		Expect(len(b)).To(Equal(40))
		Expect(len(c)).To(Equal(100))
		Expect(len(d)).To(Equal(256))
	})

	It("generates strings that are random", func() {
		a, erra := randomHexBytes(128)
		b, errb := randomHexBytes(128)
		c, errc := randomHexBytes(128)
		d, errd := randomHexBytes(128)
		Expect(erra).NotTo(HaveOccurred())
		Expect(errb).NotTo(HaveOccurred())
		Expect(errc).NotTo(HaveOccurred())
		Expect(errd).NotTo(HaveOccurred())
		Expect(a).NotTo(Equal(b))
		Expect(a).NotTo(Equal(c))
		Expect(a).NotTo(Equal(d))
		Expect(b).NotTo(Equal(c))
		Expect(b).NotTo(Equal(d))
		Expect(c).NotTo(Equal(d))
	})
})
