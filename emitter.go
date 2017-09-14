package events

import (
	"regexp"
	"sync"
)

// Emitter struct
type Emitter struct {
	inited    bool
	mutex     sync.Mutex
	listeners map[string][]Listener
	// async     bool
}

// Listener struct
type Listener struct {
	ch     chan Event
	isOnce bool
}

// Emit synchronously calls each of the listeners registered for the event named eventName,passing arguments to each
func (emitter *Emitter) Emit(eventName string, args ...interface{}) chan interface{} {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	var wg sync.WaitGroup
	done := make(chan interface{})
	event := Event{
		Name: eventName,
		Args: args,
	}

	names := emitter.getMatched(eventName)
	for _, name := range names {
		_listeners := emitter.listeners[name]
		for _, listener := range _listeners {
			wg.Add(1)
			go func(event *Event, listener Listener) {
				emitter.doEmit(event, done, listener)
				if listener.isOnce {
					emitter.removeListener(name, listener.ch)
				}
				wg.Done()
			}(&event, listener)
		}
	}
	wg.Wait()
	close(done)
	return done
}

// On adds listener for the event named eventName
func (emitter *Emitter) On(eventName string) <-chan Event {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	l := Listener{
		ch:     make(chan Event),
		isOnce: false,
	}
	if _, ok := emitter.listeners[eventName]; ok {
		emitter.listeners[eventName] = append(emitter.listeners[eventName], l)
	} else {
		emitter.listeners[eventName] = []Listener{l}
	}
	return l.ch
}

// // New returns a Emitter instance
// func New(options ...func(*Emitter)) (*Emitter, error) {
// 	e := &Emitter{}
// 	for _, opt := range options {
// 		opt(e)
// 	}
// 	return e, nil
// }

// // Async set emmiter sync
// func Async() func(*Emitter) {
// 	return func(emitter *Emitter) {
// 		emitter.async = true
// 	}
// }

// Once adds a one time listner for the event named eventName
func (emitter *Emitter) Once(eventName string) <-chan Event {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	l := Listener{
		ch:     make(chan Event),
		isOnce: true,
	}
	if _, ok := emitter.listeners[eventName]; ok {
		emitter.listeners[eventName] = append(emitter.listeners[eventName], l)
	} else {
		emitter.listeners[eventName] = []Listener{l}
	}
	return l.ch
}

// Listeners returns all listeners for the event named eventName
func (emitter *Emitter) Listeners(eventName string) []<-chan Event {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	var r []<-chan Event
	names := emitter.getMatched(eventName)
	for _, name := range names {
		ls := emitter.listeners[name]
		for _, l := range ls {
			r = append(r, l.ch)
		}
	}
	return r
}

// CloseAll closes all listener channel
func (emitter *Emitter) CloseAll(eventName string) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	names := emitter.getMatched(eventName)
	for _, value := range names {
		listeners := emitter.listeners[value]
		for _, l := range listeners {
			close(l.ch)
		}
	}
	emitter.listeners = make(map[string][]Listener)
}

// RemoveListener closes listensers specieled
func (emitter *Emitter) RemoveListener(eventName string, l <-chan Event) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()
	emitter.init()
	emitter.removeListener(eventName, l)
}

func (emitter *Emitter) removeListener(eventName string, l <-chan Event) {
	ls := emitter.listeners[eventName]
	for i, v := range ls {
		if v.ch == l {
			close(v.ch)
			emitter.listeners[eventName] = append(ls[:i], ls[i+1:]...)
			break
		}
	}
}

func (emitter *Emitter) init() {
	if !emitter.inited {
		emitter.listeners = make(map[string][]Listener)
		emitter.inited = true
	}
}

func (emitter *Emitter) doEmit(event *Event, done chan interface{}, listener Listener) {
	select {
	case <-done:
		break
	case listener.ch <- *event:
		return
	}
}

func (emitter *Emitter) getMatched(name string) []string {
	var keys []string
	for key := range emitter.listeners {
		if m, _ := regexp.MatchString(name, key); m {
			keys = append(keys, key)
		} else if m, _ := regexp.MatchString(key, name); m {
			keys = append(keys, key)
		}
	}
	return keys
}
