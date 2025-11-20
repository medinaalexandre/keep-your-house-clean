package events

import (
	"context"
	"log"
	"sync"
)

type EventDispatcher interface {
	Dispatch(event Event) error
	RegisterHandler(eventType EventType, handler EventHandler)
	Start()
	Stop()
}

type Dispatcher struct {
	handlers map[EventType][]EventHandler
	events   chan Event
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewDispatcher(ctx context.Context) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)
	return &Dispatcher{
		handlers: make(map[EventType][]EventHandler),
		events:   make(chan Event, 100),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (d *Dispatcher) RegisterHandler(eventType EventType, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *Dispatcher) Dispatch(event Event) error {
	select {
	case d.events <- event:
		return nil
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
		log.Printf("Warning: event channel is full, dropping event: %s", event.Type)
		return nil
	}
}

func (d *Dispatcher) Start() {
	d.wg.Add(1)
	go d.worker()
}

func (d *Dispatcher) Stop() {
	d.cancel()
	close(d.events)
	d.wg.Wait()
}

func (d *Dispatcher) worker() {
	defer d.wg.Done()
	for {
		select {
		case event, ok := <-d.events:
			if !ok {
				return
			}
			d.processEvent(event)
		case <-d.ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) processEvent(event Event) {
	d.mu.RLock()
	handlers := d.handlers[event.Type]
	d.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(d.ctx, event); err != nil {
			log.Printf("Error processing event %s: %v", event.Type, err)
		}
	}
}

