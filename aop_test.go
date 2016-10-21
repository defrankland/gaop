package goaop_test

import (
	"errors"
	goaop "goAop"
	"math/rand"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	msgBeforeAdvice = "test before advice"
	msgAfterAdvice  = "test after advice"
)

var (
	out string
)

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
			It("returns no error if aspect doesnt match the name of an existing aspect", func() {
				err := aspect.Create("myAspect1")
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns error if already exists", func() {

				err := aspect.Create("myAspect1")
				err = aspect.Create("myAspect1")

				Expect(err).To(MatchError("cannot create aspect: name already used"))
			})

			It("can be removed if it exists", func() {

				err := aspect.Create("myAspect1")
				err = aspect.Remove("myAspect1")

				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error if it doesnt exist", func() {
				err := aspect.Remove("myAspect1")

				Expect(err).To(MatchError("cannot delete aspect: does not exist"))
			})

		})

		Context("when adding advice to an aspect", func() {
			It("returns an error when function is nil", func() {
				err := aspect.AddAdvice(nil, goaop.ADVICE_BEFORE)
				Expect(err).To(MatchError("cannot create advice: adviceFunction is invalid"))
			})

			It("returns an error when adviceType is missing", func() {
				notAnAdviceType := goaop.ADVICE_BEFORE + goaop.ADVICE_AFTER + goaop.ADVICE_AFTER_RETURNING
				err := aspect.AddAdvice(beforeAdvice, notAnAdviceType)
				Expect(err).To(MatchError("cannot create advice: adviceType is invalid"))
			})

			It("succeeds when function and advice type are valid", func() {
				err := aspect.AddAdvice(beforeAdvice, goaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())
			})

			It("allows multiple advices to be added", func() {

				err := aspect.AddAdvice(beforeAdvice, goaop.ADVICE_BEFORE)
				err = aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER)
				Expect(err).ToNot(HaveOccurred())
			})

			It("allows multiple advices to be added under the same advice type", func() {

				err := aspect.AddAdvice(beforeAdvice, goaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(afterAdvice, goaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(beforeAdvice, goaop.ADVICE_AFTER_RETURNING)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error when the same advice is added under the same type more than once", func() {

				err := aspect.AddAdvice(afterAdvice, goaop.ADVICE_BEFORE)

				err = aspect.AddAdvice(afterAdvice, goaop.ADVICE_BEFORE)
				Expect(err).To(MatchError("cannot create advice: adviceFunction already added for adviceType"))
			})

			It("can be removed after it is added", func() {

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_BEFORE)

				err := aspect.RemoveAdvice(afterAdvice, goaop.ADVICE_BEFORE)

				Expect(err).ToNot(HaveOccurred())
			})

			It("returns error if you try to remove an advice that doesn't exist", func() {

				err := aspect.RemoveAdvice(afterAdvice, goaop.ADVICE_BEFORE)

				Expect(err).To(MatchError("cannot delete advice: does not exist"))
			})

			It("doesnt delete advice if advice func & advice type match doesnt exist", func() {

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_BEFORE)
				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)

				err := aspect.RemoveAdvice(afterAdvice, goaop.ADVICE_AFTER)

				Expect(err).To(MatchError("cannot delete advice: does not exist"))
			})
		})
	})

	Describe("aop pointcuts", func() {
		Context("when a pointcut is registered with before advice type", func() {
			It("is called after the before advice", func() {
				t := T{}
				fn := t.MyFunc

				aspect.AddAdvice(beforeAdvice, goaop.ADVICE_BEFORE)
				fnV, _ := aspect.AddPointcut("MyFunc", goaop.ADVICE_BEFORE, &t)
				aspect.Join(&fn, fnV)

				fn()

				By("checking the index of the pointcut and before messages")
				idxBeforeMsg := strings.Index(out, msgBeforeAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxBeforeMsg).To(BeNumerically("<", idxFuncMsg))
				Expect(idxBeforeMsg).ToNot(BeNumerically("==", -1))
			})

			It("takes and returns correct values", func() {

				t := T{}
				fn := t.AddInts

				aspect.AddAdvice(beforeAdvice, goaop.ADVICE_BEFORE)
				fnV, _ := aspect.AddPointcut("AddInts", goaop.ADVICE_BEFORE, &t)
				aspect.Join(&fn, fnV)

				x := rand.Int()
				y := rand.Int()

				sum := fn(x, y)

				Expect(sum).To(Equal(x + y))
			})
		})

		Context("when a pointcut is registered with after advice type", func() {
			It("is called before the after advice", func() {
				t := T{}
				fn := t.MyFunc

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER)
				fnV, _ := aspect.AddPointcut("MyFunc", goaop.ADVICE_AFTER, &t)
				aspect.Join(&fn, fnV)

				fn()

				By("checking the index of the pointcut and after messages")
				idxAfterMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterMsg).ToNot(BeNumerically("==", -1))
			})
			It("takes and returns correct values", func() {
				t := T{}
				fn := t.AddInts

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER)
				fnV, _ := aspect.AddPointcut("AddInts", goaop.ADVICE_AFTER, &t)
				aspect.Join(&fn, fnV)

				x := rand.Int()
				y := rand.Int()

				sum := fn(x, y)

				Expect(sum).To(Equal(x + y))
			})
		})

		Context("when a pointcut is registered with after-returning advice type", func() {
			It("is called if the function does not return an error", func() {
				t := T{}
				fn := t.MyFunc

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				fnV, _ := aspect.AddPointcut("MyFunc", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)

				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterReturningMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterReturningMsg).ToNot(BeNumerically("==", -1))
			})
			It("is not called if the function returns an error", func() {
				t := T{}
				fn := t.MyFuncErr

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				fnV, _ := aspect.AddPointcut("MyFuncErr", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)

				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
			})
			It("is called if the function returns an non-error type", func() {
				t := T{}
				fn := t.MyFuncInt

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				fnV, _ := aspect.AddPointcut("MyFuncInt", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)

				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).ToNot(BeNumerically("==", -1))
			})
			It("is not called for any error if the function returns several errors", func() {
				t := T{}
				fn := t.MyFuncMultiReturns

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				fnV, _ := aspect.AddPointcut("MyFuncMultiReturns", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)
				fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
			})
			It("returns the functions return values", func() {
				t := T{}
				fn := t.MyFuncMultiReturns

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)

				fnV, _ := aspect.AddPointcut("MyFuncMultiReturns", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)

				_, i, s, _, _ := fn()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
				Expect(i).To(Equal(5))
				Expect(s).To(Equal("a string"))
			})
			It("takes the provided arguments", func() {
				t := T{}
				fn := t.AddInts

				aspect.AddAdvice(afterAdvice, goaop.ADVICE_AFTER_RETURNING)
				fnV, _ := aspect.AddPointcut("AddInts", goaop.ADVICE_AFTER_RETURNING, &t)
				aspect.Join(&fn, fnV)

				x := rand.Int()
				y := rand.Int()

				sum := fn(x, y)

				Expect(sum).To(Equal(x + y))
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

func (t *T) MyFuncMultiReturns() (error, int, string, error, error) {
	out += "this is my function"
	return nil, 5, "a string", errors.New("an error"), nil
}

func (t *T) MyFuncInt() int {
	out += "this is my function"
	return 5
}

func (t *T) MyFunc() error {
	out += "this is my function"
	return nil
}

func (t *T) AddInts(x int, y int) int {
	return x + y
}
