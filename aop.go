package goaop

import (
	"errors"
	"reflect"
)

type Advice struct {
	method     func()
	adviceType string
}

type Aspect struct {
	advice Advice
}

func (a *Aspect) AddAdvice(adviceFunction func(), adviceType string) (err error) {
	if adviceFunction == nil {
		err = errors.New("cannot create advice: adviceFunction is invalid")
	} else if adviceType == "" {
		err = errors.New("cannot create advice: adviceType is invalid")
	} else {
		a.advice.method = adviceFunction
		a.advice.adviceType = adviceType
	}
	return
}

func (a *Aspect) AddPointcut(methodName string, adviceType string, i interface{}) (fn func(), err error) {
	if adviceType == "before" {
		fn = func() {
			a.advice.method()
			reflect.ValueOf(i).MethodByName("MyFunc").Call(nil)
		}
	}
	return
}
