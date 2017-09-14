# go-events
[![Build Status](https://travis-ci.org/Dongss/go-events.svg?branch=master)](https://travis-ci.org/Dongss/go-events.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/Dongss/go-events/badge.svg?branch=master)](https://coveralls.io/github/Dongss/go-events?branch=master)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/Dongss/go-events)

Events emitter by Go channel

[Documents](https://godoc.org/github.com/Dongss/go-events)

## Getting started

```
e := &Emitter{}

go e.Emit("testonece", "test me")
event := <-e.Once("testonece")
println(event.Args[0])
```

## Usage

### Emitt

`e.Emit("testonece", "test me")`

### Once

```
go e.Emit("testonece", "test me")
event := <-e.Once("testonece")
```

### On

```
e := &Emitter{}
go func() {
    <-e.Emit("testc", "test me 1", map[string]int{"dsds": 33})
    <-e.Emit("testc", "test me 2")
    <-e.Emit("testc", "test me 3")
    e.CloseAll("testc")
}()

i := 0
for event := range e.On("testc") {
    println(event.Args[0])
}
```

### Wildcard

```
e := &Emitter{}
go func() {
    <-e.Emit("testw1", "test me 1", map[string]int{"dsds": 33})
    <-e.Emit("testw2", "test me 2")
    <-e.Emit("testw3", "test me 3")
    e.CloseAll("testw*")
}()

i := 0
for event := range e.On("testw*") {
    println(event.Name)
}
```

### List listeners

```
e := &Emitter{}
e.On("name")
ls := e.Listeners("ls*")
```

### Remove listeners

```
e := &Emitter{}
l := e.On("testrml")
e.RemoveListener("testrml", l)
```

### Timeout and Force drop

```
e := &Emitter{}
ch := make(chan bool)
go func() {
    done := e.Emit("testc", "test me 1", map[string]int{"dsds": 33})
    select {
    case <-done:
        ch <- false
    case <-time.After(10000 * time.Millisecond):
        close(done)
        ch <- true
    }
}()
m := <-ch
```