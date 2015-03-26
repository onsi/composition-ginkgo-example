package prefix_tests

import (
	"github.com/onsi/composition-ginkgo-example/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var SharedContext helpers.SharedContext

var _ = Describe("Working with prefixes", func() {
	var keyA, keyB string
	BeforeEach(func() {
		keyA = SharedContext.PrefixedKey("A")
		keyB = SharedContext.PrefixedKey("B")

		Ω(SharedContext.Client.Set(keyA, "value A")).Should(Succeed())
		Ω(SharedContext.Client.Set(keyB, "value B")).Should(Succeed())
	})

	Describe("getting keys by prefix", func() {
		Context("when there are keys at the requested prefix", func() {
			It("should return the set of keys", func() {
				Ω(SharedContext.Client.GetPrefix(SharedContext.Prefix)).Should(ConsistOf("value A", "value B"))
			})
		})

		Context("when there are no keys at the requested prefix", func() {
			It("should return empty", func() {
				Ω(SharedContext.Client.GetPrefix("no-such-prefix")).Should(BeEmpty())
			})
		})
	})

	Context("when a key has been deleted", func() {
		BeforeEach(func() {
			Ω(SharedContext.Client.Delete(keyA)).Should(Succeed())
		})

		It("should not be returned when fetching a matching prefix", func() {
			Ω(SharedContext.Client.GetPrefix(SharedContext.Prefix)).Should(ConsistOf("value B"))
		})
	})

	Describe("deleting keys by prefix", func() {
		It("should remove the keys", func() {
			Ω(SharedContext.Client.DeletePrefix(SharedContext.Prefix)).Should(Succeed())
			Ω(SharedContext.Client.GetPrefix(SharedContext.Prefix)).Should(BeEmpty())
		})
	})
})
