package events

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"reflect"
	"sync"
)

type EventListener interface {
	ShouldSync() bool
	Handle(event any) error
}

var eventListeners = make(map[string][]EventListener)
var Wg sync.WaitGroup

var antsPool, _ = ants.NewPool(300) // 使用ants库创建一个池，用于异步执行事件处理函数

func On(event any, listeners ...EventListener) {
	name := getEventName(event)

	eventListeners[name] = append(eventListeners[name], listeners...)
}

func Fire(event any, payload ...any) (err error) {
	name := getEventName(event)

	if len(payload) == 0 {
		payload = append(payload, event)
	}

	if listeners, ok := eventListeners[name]; ok {
		for _, listener := range listeners {
			if listener == nil {
				continue
			}

			if listener.ShouldSync() {
				if err = listener.Handle(payload[0]); err != nil {
					return
				}
			} else {
				Wg.Add(1)
				listener := listener
				_ = antsPool.Submit(func() {
					_ = listener.Handle(payload[0])
					Wg.Done()
				})
			}
		}
	}

	return
}

func getEventName(event any) string {
	if v, ok := event.(string); ok {
		return v
	}

	t := reflect.TypeOf(event)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return fmt.Sprintf("%s/%s", t.PkgPath(), t.Name())
}
