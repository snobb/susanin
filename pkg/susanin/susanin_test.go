package susanin

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func dummy(w http.ResponseWriter, r *http.Request)    {}
func dynamic(w http.ResponseWriter, r *http.Request)  {}
func dynamic1(w http.ResponseWriter, r *http.Request) {}
func static(w http.ResponseWriter, r *http.Request)   {}
func static1(w http.ResponseWriter, r *http.Request)  {}
func static2(w http.ResponseWriter, r *http.Request)  {}
func byName(w http.ResponseWriter, r *http.Request)   {}
func splat(w http.ResponseWriter, r *http.Request)    {}

func Test(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Susanin Handle method", func() {
		g.It("Should create a search chain for given path", func() {
			s := NewSusanin()
			s.Handle("/test/this/:uri", dummy)

			cur := s.root
			Expect(cur.nextConst).NotTo(BeNil())
			cur, ok := cur.nextConst["test"]
			Expect(ok).To(BeTrue())
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("test"))

			cur, ok = cur.nextConst["this"]
			Expect(ok).To(BeTrue())
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("this"))

			cur = cur.nextVar
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("uri"))

			s.Handle("/hello/*", dummy)
			cur, ok = s.root.nextConst["hello"]
			Expect(ok).To(BeTrue())
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("hello"))

			Expect(cur.nextSplat).NotTo(BeNil())
			cur = cur.nextSplat
			Expect(ok).To(BeTrue())
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("*"))
			Expect(cur.handler).NotTo(BeNil())

			s.Handle("/test/this/:uri/test", dummy)
			cur = s.root.nextConst["test"].nextConst["this"].nextVar
			cur, ok = cur.nextConst["test"]
			Expect(ok).To(BeTrue())
			Expect(cur).NotTo(BeNil())
			Expect(cur.name).To(Equal("test"))
		})

		g.It("Should return an error if splat is in the middle of the path", func() {
			s := NewSusanin()
			err := s.Handle("/test/*/hello", dummy)
			Expect(err).NotTo(BeNil())
		})

		g.It("Should return an error if splat is in the middle of the path", func() {
			s := NewSusanin()
			err := s.Handle("/test/hello/*/*", dummy)
			Expect(err).NotTo(BeNil())
		})

		g.It("Should return an error if different variable patterns set at the same level", func() {
			s := NewSusanin()
			err := s.Handle("/test/:param1/hello", dummy)
			Expect(err).To(BeNil())

			err = s.Handle("/test/:param2/hello", dummy)
			Expect(err).NotTo(BeNil())
		})
	})

	g.Describe("Susanin Lookup method", func() {
		var s *Susanin

		g.Before(func() {
			s = NewSusanin()
			s.Handle("/hello/:name", dynamic)
			s.Handle("/hello/:name/by-name", dynamic1)
			s.Handle("/hello/*", splat)
			s.Handle("/hello/test", static)
			s.Handle("/hello/test/all", static1)
			s.Handle("/test/all", static2)
		})

		g.It("Should find a static handler", func() {
			handler, values, err := s.Lookup("/hello/test")
			Expect(err).To(BeNil())
			Expect(values).To(BeEmpty())
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(static).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})

		g.It("Should find a static1 handler", func() {
			handler, values, err := s.Lookup("/hello/test/all")
			Expect(err).To(BeNil())
			Expect(values).To(BeEmpty())
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(static1).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})

		g.It("Should find a static2 handler", func() {
			handler, values, err := s.Lookup("/test/all")
			Expect(err).To(BeNil())
			Expect(values).To(BeEmpty())
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(static2).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})

		g.It("Should find a dynamic handler", func() {
			handler, values, err := s.Lookup("/hello/alex")
			Expect(err).To(BeNil())
			Expect(values).To(HaveKey("name"))
			Expect(handler).NotTo(BeNil())
			Expect(values["name"]).To(Equal("alex"))
			f1 := runtime.FuncForPC(reflect.ValueOf(dynamic).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})

		g.It("Should find a dynamic1 handler", func() {
			handler, values, err := s.Lookup("/hello/alex/by-name")
			Expect(err).To(BeNil())
			Expect(values).To(HaveKey("name"))
			Expect(handler).NotTo(BeNil())
			Expect(values["name"]).To(Equal("alex"))
			f1 := runtime.FuncForPC(reflect.ValueOf(dynamic1).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})

		g.It("Should find a splat handler", func() {
			handler, values, err := s.Lookup("/hello/alex/nonexistant")
			Expect(err).To(BeNil())
			Expect(values).To(HaveKeyWithValue("name", "alex"))
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(splat).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})
	})

	g.Describe("Test splat fallback", func() {
		var s *Susanin
		g.Before(func() {
			s = NewSusanin()
			s.Handle("/short", static)
			s.Handle("/*", splat)
		})

		g.It("Should fallback to splat if no longer matches with the specific path", func() {
			handler, values, err := s.Lookup("/short/")
			Expect(err).To(BeNil())
			Expect(values).To(BeEmpty())
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(static).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))

			handler, values, err = s.Lookup("/short/aaa/bbb")
			Expect(err).To(BeNil())
			Expect(values).To(BeEmpty())
			Expect(handler).NotTo(BeNil())
			f1 = runtime.FuncForPC(reflect.ValueOf(splat).Pointer()).Name()
			f2 = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})
	})

	g.Describe("Test splat fallback after matching variable", func() {
		var s *Susanin
		g.Before(func() {
			s = NewSusanin()
			s.Handle("/hello/:fname/:lname", static)
			s.Handle("/hello/*", splat)
		})

		g.It("Should fallback to splat if longer matches with the specific path", func() {
			handler, values, err := s.Lookup("/hello/john/doe/")
			Expect(err).To(BeNil())
			Expect(values).To(HaveKeyWithValue("fname", "john"))
			Expect(values).To(HaveKeyWithValue("lname", "doe"))
			Expect(handler).NotTo(BeNil())
			f1 := runtime.FuncForPC(reflect.ValueOf(static).Pointer()).Name()
			f2 := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))

			handler, values, err = s.Lookup("/hello/john")
			Expect(err).To(BeNil())
			// currently it still fills the values during the longest match search
			Expect(values).To(HaveKeyWithValue("fname", "john"))
			Expect(handler).NotTo(BeNil())
			f1 = runtime.FuncForPC(reflect.ValueOf(splat).Pointer()).Name()
			f2 = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			Expect(f1).To(Equal(f2))
		})
	})
}
