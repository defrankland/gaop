package gaop_test

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/defrankland/gaop"

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

var _ = Describe("aop", func() {

	var (
		aspect gaop.Aspect
	)

	BeforeEach(func() {
		aspect = gaop.Aspect{}
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

			It("can be found by name", func() {

				err := aspect.Create("myAspect1")
				Expect(err).ToNot(HaveOccurred())

				_, err = gaop.GetAspectByName("myAspect1")
				Expect(err).ToNot(HaveOccurred())

			})

			It("gives error when not found by name", func() {

				_, err := gaop.GetAspectByName("myAspect1")
				Expect(err).ToNot(HaveOccurred())

			})
		})

		Context("when adding advice to an aspect", func() {
			It("returns an error when function is nil", func() {
				err := aspect.AddAdvice(nil, gaop.ADVICE_BEFORE)
				Expect(err).To(MatchError("cannot create advice: adviceFunction is invalid"))
			})

			It("returns an error when adviceType is missing", func() {
				notAnAdviceType := gaop.ADVICE_BEFORE + gaop.ADVICE_AFTER + gaop.ADVICE_AFTER_RETURNING
				err := aspect.AddAdvice(beforeAdvice, notAnAdviceType)
				Expect(err).To(MatchError("cannot create advice: adviceType is invalid"))
			})

			It("succeeds when function and advice type are valid", func() {
				err := aspect.AddAdvice(beforeAdvice, gaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())
			})

			It("allows multiple advices to be added", func() {

				err := aspect.AddAdvice(beforeAdvice, gaop.ADVICE_BEFORE)
				err = aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER)
				Expect(err).ToNot(HaveOccurred())
			})

			It("allows multiple advices to be added under the same advice type", func() {

				err := aspect.AddAdvice(beforeAdvice, gaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(afterAdvice, gaop.ADVICE_BEFORE)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(beforeAdvice, gaop.ADVICE_AFTER_RETURNING)
				Expect(err).ToNot(HaveOccurred())

				err = aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error when the same advice is added under the same type more than once", func() {

				err := aspect.AddAdvice(afterAdvice, gaop.ADVICE_BEFORE)

				err = aspect.AddAdvice(afterAdvice, gaop.ADVICE_BEFORE)
				Expect(err).To(MatchError("cannot create advice: adviceFunction already added for adviceType"))
			})

			It("can be removed after it is added", func() {

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_BEFORE)

				err := aspect.RemoveAdvice(afterAdvice, gaop.ADVICE_BEFORE)

				Expect(err).ToNot(HaveOccurred())
			})

			It("returns error if you try to remove an advice that doesn't exist", func() {

				err := aspect.RemoveAdvice(afterAdvice, gaop.ADVICE_BEFORE)

				Expect(err).To(MatchError("cannot delete advice: does not exist"))
			})

			It("doesnt delete advice if advice func & advice type match doesnt exist", func() {

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_BEFORE)
				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)

				err := aspect.RemoveAdvice(afterAdvice, gaop.ADVICE_AFTER)

				Expect(err).To(MatchError("cannot delete advice: does not exist"))
			})
		})
	})

	Describe("aop pointcuts", func() {
		Context("when a pointcut is registered with before advice type", func() {
			It("is called after the before advice", func() {
				t := New()

				aspect.AddAdvice(beforeAdvice, gaop.ADVICE_BEFORE)
				aspect.AddPointcut("MyFuncImpl", gaop.ADVICE_BEFORE, t, &t.MyFunc)
				t.MyFunc()

				By("checking the index of the pointcut and before messages")
				idxBeforeMsg := strings.Index(out, msgBeforeAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxBeforeMsg).To(BeNumerically("<", idxFuncMsg))
				Expect(idxBeforeMsg).ToNot(BeNumerically("==", -1))
			})

			It("takes and returns correct values", func() {

				t := New()

				aspect.AddAdvice(beforeAdvice, gaop.ADVICE_BEFORE)
				aspect.AddPointcut("AddIntsImpl", gaop.ADVICE_BEFORE, t, &t.AddInts)

				x := rand.Int()
				y := rand.Int()

				sum := t.AddInts(x, y)

				Expect(sum).To(Equal(x + y))
			})
		})

		Context("when a pointcut is registered with after advice type", func() {
			It("is called before the after advice", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER)
				aspect.AddPointcut("MyFuncImpl", gaop.ADVICE_AFTER, t, &t.MyFunc)

				t.MyFunc()

				By("checking the index of the pointcut and after messages")
				idxAfterMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterMsg).ToNot(BeNumerically("==", -1))
			})

			It("takes and returns correct values", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER)
				aspect.AddPointcut("AddIntsImpl", gaop.ADVICE_AFTER, t, &t.AddInts)

				x := rand.Int()
				y := rand.Int()

				sum := t.AddInts(x, y)

				Expect(sum).To(Equal(x + y))
			})
		})

		Context("when a pointcut is registered with after-returning advice type", func() {
			It("is called if the function does not return an error", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("MyFuncImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.MyFunc)

				t.MyFunc()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxAfterReturningMsg).To(BeNumerically(">", idxFuncMsg))
				Expect(idxAfterReturningMsg).ToNot(BeNumerically("==", -1))
			})

			It("is not called if the function returns an error", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("MyFuncErrImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.MyFuncErr)

				t.MyFuncErr()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
			})

			It("is called if the function returns an non-error type", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("MyFuncIntImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.MyFuncInt)

				t.MyFuncInt()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).ToNot(BeNumerically("==", -1))
			})

			It("is not called for any error if the function returns several errors", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("MyFuncMultiReturnsImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.MyFuncMultiReturns)

				t.MyFuncMultiReturns()

				idxAfterReturningMsg := strings.Index(out, msgAfterAdvice)
				idxFuncMsg := strings.Index(out, "this is my function")
				Expect(idxFuncMsg).To(BeNumerically("==", 0))
				Expect(idxAfterReturningMsg).To(BeNumerically("==", -1))
			})

			It("returns the functions return values", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("MyFuncMultiReturnsImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.MyFuncMultiReturns)

				err1, i, s, err2, err3 := t.MyFuncMultiReturns()

				Expect(i).To(Equal(5))
				Expect(s).To(Equal("a string"))
				Expect(err1).To(BeNil())
				Expect(err2.Error()).To(Equal("error 2"))
				Expect(err3.Error()).To(Equal("error 3"))

			})

			It("takes the provided arguments", func() {
				t := New()

				aspect.AddAdvice(afterAdvice, gaop.ADVICE_AFTER_RETURNING)
				aspect.AddPointcut("AddIntsImpl", gaop.ADVICE_AFTER_RETURNING, t, &t.AddInts)

				x := rand.Int()
				y := rand.Int()

				sum := t.AddInts(x, y)

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

type T struct {
	funcs
}

type funcs struct {
	MyFunc             func() error
	MyFuncMultiReturns func() (error, int, string, error, error)
	MyFuncInt          func() int
	AddInts            func(x int, y int) int
	MyFuncErr          func() error
}

func New() *T {

	t := T{}
	t.MyFunc = t.MyFuncImpl
	t.MyFuncMultiReturns = t.MyFuncMultiReturnsImpl
	t.MyFuncInt = t.MyFuncIntImpl
	t.AddInts = t.AddIntsImpl
	t.MyFuncErr = t.MyFuncErrImpl

	return &t
}

func (t *T) MyFuncImpl() error {
	out += "this is my function"
	return nil
}

func (t *T) MyFuncErrImpl() error {
	out += "this is my function"

	return errors.New("test error return")
}

func (t *T) MyFuncMultiReturnsImpl() (error, int, string, error, error) {
	out += "this is my function"
	return nil, 5, "a string", errors.New("error 2"), errors.New("error 3")
}

func (t *T) MyFuncIntImpl() int {
	out += "this is my function"
	return 5
}

func (t *T) AddIntsImpl(x int, y int) int {
	return x + y
}
