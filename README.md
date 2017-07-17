# gaop

Aspect-oriented Programming library for golang. 

## Current State 

Supports advice types 
- Before
- After
- After returning

Currently it is a slight bummer to create a pointcut. The method to do so takes both the pointcut methodName as a string and the pointcut method itself (as an interface{}):

```go
AddPointcut(methodName string, adviceType AopAdviceType, i, pointcut interface{}) (err error)
``` 

# Usage:

See tests for intended modes of operation.

Intended usage is to create a type where the methods map to fields like this: 

```go
type Doggo struct {
	Bark   func(int, bool) error
	aspect gaop.Aspect
}

func (d *Doggo) BarkImpl(numBarks int, gotChokedUp bool) error {
	fmt.Println("BARK!!")
}
```

Next create some advice (note the type is `ADVICE_BEFORE`): 

```go
func openMouth() {
	fmt.Println("Opening Mouth...")
}
```


Map the method to the field and add the advice: 

```go
func NewDoggo() *Doggo {

	doggo := Doggo{}
	doggo.Bark = doggo.BarkImpl

	doggo.aspect.AddAdvice(openMouth, gaop.ADVICE_BEFORE)
	doggo.aspect.AddPointcut("BarkImpl", gaop.ADVICE_BEFORE, &doggo, &doggo.Bark)

	return &doggo
}
```

Now run:

```go
func main() {

	doggo := NewDoggo()

	err := doggo.Bark(5, false)
	if err != nil {
		fmt.Println(err)
	}
}
```

Note that `Bark()` is called, not `BarkImpl()`. 

Gives the output:
```
⇒  go run main.go
Opening Mouth...
BARK!!
```

## Advice Type Constants
- `ADVICE_BEFORE`
- `ADVICE_AFTER`
- `ADVICE_AFTER_RETURNING`