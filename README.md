# gaop

Aspect-oriented Programming library for golang. 

## Current State 

Supports advice types 
- Before
- After
- After returning

Currently it is a slight bummer to create a pointcut. The method to do so takes both the pointcut methodName as a string and the pointcut method itself (as an interface{}):

```
AddPointcut(methodName string, adviceType AopAdviceType, i, pointcut interface{}) (err error)
``` 

# Usage:

Intended usage is to create a type where the methods map to fields like this: 

```
type Doggo struct {
	Bark func()
	aspect gaop.Aspect
}

func (d *Doggo) BarkImpl() {
	fmt.Println("BARK!!")
}
```

Next create some advice: 

```
func openMouth() {
	fmt.Println("Opening Mouth...")
}
```


Map the method to the field and add the advice: 

```
func NewDoggo() *Doggo {

	doggo := Doggo{}
	doggo.Bark = doggo.BarkImpl

	doggo.aspect.AddAdvice(openMouth, gaop.ADVICE_BEFORE)
	doggo.aspect.AddPointcut("BarkImpl", gaop.ADVICE_BEFORE, &doggo, &doggo.Bark)

	return &doggo
}
```

Now run:

```
func main() {

	doggo := NewDoggo()

	doggo.Bark()
}
```

Note that Bark() is called, not BarkImpl(). 

Gives the output:
```
⇒  go run main.go
Opening Mouth...
BARK!!
```
 