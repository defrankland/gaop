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

