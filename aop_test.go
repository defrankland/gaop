package goaop_test

import (
	goaop "goAop"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	msgBeforeAdvice = "test before advice"
)

var out string

type T struct{}

var _ = Describe("aop", func() {

	var (
		aspect goaop.Aspect
	)
	BeforeEach(func() {
		aspect = goaop.Aspect{}
		out = ""
	})

	Describe("aspects", func() {
		Context("when creating an aspect", func() {
			It("returns an error when function is nil", func() {
				err := aspect.AddAdvice(nil, "before")
				Expect(err).To(MatchError("cannot create advice: adviceFunction is invalid"))
			})

			It("returns an error when adviceType is missing", func() {
				err := aspect.AddAdvice(beforeAdvice, "")
				Expect(err).To(MatchError("cannot create advice: adviceType is invalid"))
			})

			It("succeeds when function and advice type are valid", func() {
				err := aspect.AddAdvice(beforeAdvice, "before")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("aop pointcuts", func() {
		Context("when a pointcut is registered with before advice type", func() {
			It("is called after the before advice", func() {
				t := T{}
				aspect.AddAdvice(beforeAdvice, "before")
				fn, _ := aspect.AddPointcut("MyFunc", "before", &t)

				fn()
				Expect(out).To(ContainSubstring(msgBeforeAdvice))
				Expect(out).To(ContainSubstring("this is my function"))
			})
		})
	})
})

func beforeAdvice() {
	out += msgBeforeAdvice
}

func (t *T) MyFunc() {
	out += "this is my function"
}
