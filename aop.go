package goaop

import (
	"errors"
	. "reflect"
)

var aspects []Aspect

type Advice struct {
	Method Method
	Type   string
}

type Aspect struct {
	Name   string
	advice Advice
}

func (a *Aspect) Create(aspectName string) (err error) {

	if a.Index(aspectName) != -1 {
		return errors.New("cannot create aspect: name already used")
	}

	a.Name = aspectName
	aspects = append(aspects, *a)

	return
}

func (a *Aspect) Remove(aspectName string) (err error) {

	i := a.Index(aspectName)

	if i != -1 {
		aspects = append(aspects[:i], aspects[i+1:]...)
		return
	}
	return errors.New("cannot delete aspect: does not exist")
}

func (a *Aspect) Index(aspectName string) int {

	for i, aspect := range aspects {
		if aspectName == aspect.Name {
			return i
		}
	}
	return -1
}

func (a *Aspect) AddAdvice(adviceFunction interface{}, adviceType string) (err error) {
	if adviceFunction == nil {
		err = errors.New("cannot create advice: adviceFunction is invalid")
	} else if adviceType == "" {
		err = errors.New("cannot create advice: adviceType is invalid")
	} else {
		a.advice.Method = Method{Func: ValueOf(adviceFunction), Type: TypeOf(adviceFunction)}
		a.advice.Type = adviceType
	}
	return
}

func (a *Aspect) AddPointcut(methodName string, adviceType string, i interface{}) (fn func(args []Value) []Value, err error) {

	if adviceType == "before" {
		fn = func(args []Value) []Value {
			a.advice.Method.Func.Call(nil)
			returnValues := ValueOf(i).MethodByName(methodName).Call(args)
			return returnValues
		}
	} else if adviceType == "after" {
		fn = func(args []Value) []Value {
			returnValues := ValueOf(i).MethodByName(methodName).Call(args)
			a.advice.Method.Func.Call(nil)
			return returnValues
		}
	} else if adviceType == "after-returning" {
		fn = func(args []Value) []Value {
			returnValues := ValueOf(i).MethodByName(methodName).Call(args)

			for idx := 0; idx < len(returnValues); idx++ {
				if returnValues[idx].Type() == TypeOf((*error)(nil)).Elem() && !returnValues[idx].IsNil() {
					return returnValues
				}
			}

			a.advice.Method.Func.Call(nil)
			ValueOf(i).MethodByName("MyFunc").Call(nil)
			return returnValues

		}
	}

	return
}

func (a *Aspect) Join(fptr interface{}, pointcut func([]Value) []Value) {

	fn := ValueOf(fptr).Elem()
	fn.Set(MakeFunc(fn.Type(), pointcut))

}
