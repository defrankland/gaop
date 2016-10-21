package gaop

import (
	"errors"
	. "reflect"
)

const (
	ADVICE_BEFORE AopAdviceType = iota
	ADVICE_AFTER
	ADVICE_AFTER_RETURNING
)

var aspects []Aspect

type AopAdviceType int

type Aspect struct {
	Name    string
	advices []Advice
}

type Advice struct {
	Method Method
	Type   AopAdviceType
}

func (a *Aspect) Create(aspectName string) (err error) {

	if index(aspectName) != -1 {
		return errors.New("cannot create aspect: name already used")
	}

	a.Name = aspectName
	aspects = append(aspects, *a)

	return
}

func (a *Aspect) Remove(aspectName string) (err error) {

	i := index(aspectName)

	if i != -1 {
		aspects = append(aspects[:i], aspects[i+1:]...)
		return
	}
	return errors.New("cannot delete aspect: does not exist")
}

func index(aspectName string) int {

	for i, aspect := range aspects {
		if aspectName == aspect.Name {
			return i
		}
	}
	return -1
}

func GetAspectByName(name string) (a Aspect, err error) {

	i := index(name)

	if i != -1 {
		a = aspects[i]

		return
	}

	return a, errors.New("cannot find aspect: does not exist")
}

func (a *Aspect) AddAdvice(adviceFunction interface{}, adviceType AopAdviceType) (err error) {

	if adviceFunction == nil {

		err = errors.New("cannot create advice: adviceFunction is invalid")

	} else if adviceType != ADVICE_BEFORE && adviceType != ADVICE_AFTER && adviceType != ADVICE_AFTER_RETURNING {

		err = errors.New("cannot create advice: adviceType is invalid")

	} else {

		for _, advice := range a.advices {
			if ValueOf(adviceFunction) == advice.Method.Func && adviceType == advice.Type {
				return errors.New("cannot create advice: adviceFunction already added for adviceType")
			}
		}

		a.advices = append(a.advices, Advice{
			Method: Method{Func: ValueOf(adviceFunction), Type: TypeOf(adviceFunction)},
			Type:   adviceType,
		})
	}

	return
}

func (a *Aspect) RemoveAdvice(adviceFunction interface{}, adviceType AopAdviceType) (err error) {

	i := a.GetAdviceIndex(adviceFunction, adviceType)

	if i != -1 {
		a.advices = append(a.advices[:i], a.advices[i+1:]...)
		return
	}

	return errors.New("cannot delete advice: does not exist")
}

func (a *Aspect) GetAdviceIndex(adviceFunc interface{}, adviceType AopAdviceType) (index int) {

	for i, advice := range a.advices {
		if ValueOf(adviceFunc) == advice.Method.Func && adviceType == advice.Type {
			return i
		}
	}

	return -1
}

func (a *Aspect) AddPointcut(methodName string, adviceType AopAdviceType, i interface{}) (fn func(args []Value) []Value, err error) {

	if adviceType == ADVICE_BEFORE {

		fn = func(args []Value) []Value {
			for j, advice := range a.advices {
				if advice.Type == ADVICE_BEFORE {
					a.advices[j].Method.Func.Call(nil)
				}
			}
			returnValues := ValueOf(i).MethodByName(methodName).Call(args)
			return returnValues
		}

	} else if adviceType == ADVICE_AFTER {

		fn = func(args []Value) []Value {
			returnValues := ValueOf(i).MethodByName(methodName).Call(args)

			for j, advice := range a.advices {
				if advice.Type == ADVICE_AFTER {
					a.advices[j].Method.Func.Call(nil)
				}
			}
			return returnValues
		}

	} else if adviceType == ADVICE_AFTER_RETURNING {

		fn = func(args []Value) []Value {

			returnValues := ValueOf(i).MethodByName(methodName).Call(args)

			for idx := 0; idx < len(returnValues); idx++ {
				if returnValues[idx].Type() == TypeOf((*error)(nil)).Elem() && !returnValues[idx].IsNil() {
					return returnValues
				}
			}

			for j, advice := range a.advices {
				if advice.Type == ADVICE_AFTER_RETURNING {
					a.advices[j].Method.Func.Call(nil)
				}
			}

			return returnValues
		}

	}

	return
}

func (a *Aspect) Join(fptr interface{}, pointcut func([]Value) []Value) {

	fn := ValueOf(fptr).Elem()
	fn.Set(MakeFunc(fn.Type(), pointcut))

}
