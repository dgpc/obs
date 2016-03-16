package metrics

import (
	"fmt"
	"sync"
)

type Sink interface {
	Handle(metric string, tags Tags, value float64, metricType metricType) error
	Flush() error
	Close()
}

type nullSink struct{}

func (sink *nullSink) Handle(metric string, tags Tags, value float64, metricType metricType) error {
	return nil
}

func (sink *nullSink) Flush() error {
	return nil
}

func (sink *nullSink) Close() {
}

var NullSink Sink = &nullSink{}

type MockSink struct {
	mutex       sync.Mutex
	Invocations map[string]int
	NumFlushes  int
}

func (sink *MockSink) Handle(metric string, tags Tags, value float64, metricType metricType) error {
	sink.mutex.Lock()
	defer sink.mutex.Unlock()

	formatted := fmt.Sprintf("%v, %v, %v, %v\n", metric, tags, value, metricType)
	sink.Invocations[formatted]++
	return nil
}

func (sink *MockSink) Flush() error {
	sink.mutex.Lock()
	defer sink.mutex.Unlock()

	sink.NumFlushes--
	return nil
}
func (sink *MockSink) Close() {}