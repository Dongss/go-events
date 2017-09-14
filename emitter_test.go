package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	e := &Emitter{}
	go e.Emit("testonece", "test me")
	event := <-e.Once("testonece")
	assert.Equal(t, "test me", event.Args[0], "args should be equal")
	assert.Equal(t, "testonece", event.Name, "name should be equal")
	assert.Equal(t, 0, len(e.listeners["testonece"]), "listener once should be removed")
}

func TestOn(t *testing.T) {
	e := &Emitter{}
	go func() {
		<-e.Emit("testc", "test me 1", map[string]int{"dsds": 33})
		<-e.Emit("testc", "test me 2")
		<-e.Emit("testc", "test me 3")
		e.CloseAll("testc")
	}()

	i := 0
	for event := range e.On("testc") {
		t.Log(event)
		assert.Equal(t, "testc", event.Name)
		i++
	}
	assert.Equal(t, i, 3)
}

func TestWildcardOn(t *testing.T) {
	e := &Emitter{}
	go func() {
		<-e.Emit("testw1", "test me 1", map[string]int{"dsds": 33})
		<-e.Emit("testw2", "test me 2")
		<-e.Emit("testw3", "test me 3")
		e.CloseAll("testw*")
	}()

	i := 0
	for event := range e.On("testw*") {
		t.Log(event)
		i++
	}
	assert.Equal(t, 3, i, "3 events")
	assert.Equal(t, 0, len(e.listeners["testw1"]), "listeners should be removed")
	assert.Equal(t, 0, len(e.listeners["testw2"]), "listeners should be removed")
	assert.Equal(t, 0, len(e.listeners["testw3"]), "listeners should be removed")
}

func TestCloseAll(t *testing.T) {
	e := &Emitter{}
	go func() {
		<-e.Emit("testc", "test me 1", map[string]int{"dsds": 33})
		<-e.Emit("testc", "test me 2")
		<-e.Emit("testc", "test me 3")
		e.CloseAll("testc")
	}()

	i := 0
	for event := range e.On("testc") {
		t.Log(event)
		assert.Equal(t, "testc", event.Name)
		i++
	}
	assert.Equal(t, 3, i)
	assert.Equal(t, 0, len(e.listeners["testc"]), "listeners should be removed")
}

func TestRemoveListener(t *testing.T) {
	e := &Emitter{}
	l := e.On("testrml")
	e.RemoveListener("testrml", l)
	assert.Equal(t, 0, len(e.listeners["testrml"]))
}

func TestListeners(t *testing.T) {
	e := &Emitter{}
	e.On("ls")
	l := e.On("ls2")
	e.On("ls3")
	ls := e.Listeners("ls*")
	assert.Equal(t, 3, len(ls))

	e.RemoveListener("ls2", l)
	ls2 := e.Listeners("ls*")
	assert.Equal(t, 2, len(ls2))
}

func TestDrop(t *testing.T) { //TODO
	e := &Emitter{}
	ch := make(chan bool)
	go func() {
		done := e.Emit("testc", "test me 1", map[string]int{"dsds": 33})
		select {
		case <-done:
			ch <- false
			t.Log("emit done")
		case <-time.After(1 * time.Millisecond):
			t.Log("timeout")
			close(done)
			ch <- true
		}
	}()
	m := <-ch
	t.Log(m)
}

func TestEmittBatch(t *testing.T) { // synchronous
	// e, err := New(Async())
	e := &Emitter{}
	// if err != nil {
	// 	panic(err)
	// }

	n := 100
	myargs := genEvents(n)
	go func() {
		m := 0

		for i := 0; i < n; i++ {
			<-e.Emit("testc", myargs[i])

			t.Log(m)
			if m == n-1 {
				t.Log("going to close all", i, m)
				e.CloseAll("testc")
			} else {
				m++
			}
		}
	}()

	i := 0
	for event := range e.On("testc") {
		i++
		t.Log(event)
		assert.Equal(t, "testc", event.Name)
		assert.Equal(t, i, event.Args[0])

	}
	assert.Equal(t, n, i)
}

func genEvents(n int) []int {
	var r []int
	for i := 1; i <= n; i++ {
		r = append(r, i)
	}
	return r
}
