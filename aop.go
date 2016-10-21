package goaop

import (
	"errors"
	. "reflect"
)

type Advice struct {
	Method Method
	Type   string
}

type Aspect struct {
	advice Advice
}

func (a *Aspect) AddAdvice(adviceFunction interface{}, adviceType string) (err error) {
	if adviceFunction == nil {
		err = errors.New("cannot create advice: adviceFunction is invalid")
	} else if adviceType == "" {
		err = errors.New("cannot create advice: adviceType is invalid")
	} else {
		a.advice.Method.Func = ValueOf(adviceFunction)
		a.advice.Method.Type = TypeOf(adviceFunction)
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
			a.advice.Method.Func.Call(args)
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

func (a *Aspect) MakeJoin(fptr interface{}, pointcut func([]Value) []Value) {

	fn := ValueOf(fptr).Elem()
	fn.Set(MakeFunc(fn.Type(), pointcut))

}
