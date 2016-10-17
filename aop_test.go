package goaop_test

import (
	"errors"
	goaop "goAop"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	msgBeforeAdvice = "test before advice"
	msgAfterAdvice  = "test after advice"
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

				By("checking the index of the pointcut and before messages")
				idxBeforeMsg := strings.Index(out, msgBeforeAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxBeforeMsg).To(BeNumerically("<", idxFuncMsg))
				Expect(idxBeforeMsg).ToNot(BeNumerically("==", -1))
			})
		})
		Context("when a pointcut is registered with after advice type", func() {
			It("is called before the after advice", func() {
				t := T{}
				aspect.AddAdvice(afterAdvice, "after")
				fn, _ := aspect.AddPointcut("MyFunc", "after", &t)

				fn()

				By("checking the index of the pointcut and after messages")
				idxAfterMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterMsg).ToNot(BeNumerically("==", -1))
			})
		})
		Context("when a pointcut is registered with after-returning advice type", func() {
			It("is called if the function does not return an error", func() {
				t := T{}
				aspect.AddAdvice(afterAdvice, "after-returning")
				fn, _ := aspect.AddPointcut("MyFunc", "after-returning", &t)

				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterReturningMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterReturningMsg).ToNot(BeNumerically("==", -1))
			})
			It("is not called if the function returns an error", func() {
				t := T{}
				aspect.AddAdvice(afterAdvice, "after-returning")
				fn, _ := aspect.AddPointcut("MyFuncErr", "after-returning", &t)

				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
			})
		})
	})
})

func beforeAdvice() {
	out += msgBeforeAdvice
}

func afterAdvice() {
	out += msgAfterAdvice
}

func (t *T) MyFuncErr() error {
	out += "this is my function"

	return errors.New("test error return")
}

func (t *T) MyFunc() error {
	out += "this is my function"
	return nil
}
