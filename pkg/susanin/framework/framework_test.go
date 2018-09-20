package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func dummy(w http.ResponseWriter, r *http.Request) {}

func TestHelloWorld(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic framework", func() {
		f := NewFramework()

		Expect(f.methods[mGet]).To(BeNil())
		Expect(f.methods[mPut]).To(BeNil())
		Expect(f.methods[mPost]).To(BeNil())
		Expect(f.methods[mDelete]).To(BeNil())
		Expect(f.methods[mPatch]).To(BeNil())

		f.Get("/", dummy)
		Expect(f.methods[mGet]).NotTo(BeNil())

		f.Put("/", dummy)
		Expect(f.methods[mPut]).NotTo(BeNil())

		f.Post("/", dummy)
		Expect(f.methods[mPost]).NotTo(BeNil())

		f.Delete("/", dummy)
		Expect(f.methods[mDelete]).NotTo(BeNil())

		f.Patch("/", dummy)
		Expect(f.methods[mPatch]).NotTo(BeNil())
	})
}
