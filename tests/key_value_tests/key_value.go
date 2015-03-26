package key_value_tests

import (
	"github.com/onsi/composition-ginkgo-example/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var SharedContext helpers.SharedContext

var _ = Describe("Basic Key-Value store", func() {
	var keyA, keyB string
	BeforeEach(func() {
		keyA = SharedContext.PrefixedKey("A")
		keyB = SharedContext.PrefixedKey("B")
	})

	Describe("storing off keys", func() {
		It("should store the key", func() {
			Ω(SharedContext.Client.Set(keyA, "value A")).Should(Succeed())
			Ω(SharedContext.Client.Get(keyA)).Should(Equal("value A"))
		})
	})

	Context("with keys in the database", func() {
		BeforeEach(func() {
			Ω(SharedContext.Client.Set(keyA, "value A")).Should(Succeed())
			Ω(SharedContext.Client.Set(keyB, "value B")).Should(Succeed())
		})

		Describe("getting a single key", func() {
			It("should return the key in question", func() {
				Ω(SharedContext.Client.Get(keyA)).Should(Equal("value A"))
			})

			Context("when the key does not exist", func() {
				It("should return the empty string", func() {
					Ω(SharedContext.Client.Get(SharedContext.PrefixedKey("bloop"))).Should(BeEmpty())
				})
			})
		})

		Describe("deleting a single key", func() {
			It("should no longer have the key", func() {
				Ω(SharedContext.Client.Delete(keyA)).Should(Succeed())
				Ω(SharedContext.Client.Get(keyA)).Should(BeEmpty())
			})
		})
	})
})
