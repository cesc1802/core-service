package pubsub

import "github.com/cesc1802/core-service/events"

type InMemoryPubSubOpt struct {
	name   string
	prefix string
}
type InMemoryPubSub struct {
	stream events.Stream
	*InMemoryPubSubOpt
}

func New(name, prefix string) *InMemoryPubSub {
	stream, _ := events.NewStream()
	return &InMemoryPubSub{
		InMemoryPubSubOpt: &InMemoryPubSubOpt{
			name:   name,
			prefix: prefix,
		},
		stream: stream,
	}
}

func (imps *InMemoryPubSub) Get() interface{} {
	if imps.stream != nil {
		return imps.stream
	}
	return nil
}

func (imps *InMemoryPubSub) GetPrefix() string {
	return imps.prefix
}

func (imps *InMemoryPubSub) Start() error {
	if err := imps.Configure(); err != nil {
		return err
	}

	return nil
}

func (imps *InMemoryPubSub) Stop() error {
	return nil
}

func (imps *InMemoryPubSub) Name() string {
	return imps.name
}
func (imps *InMemoryPubSub) Configure() error {
	if imps.stream == nil {
		stream, err := events.NewStream()

		if err != nil {
			return err
		}
		imps.stream = stream
	}
	return nil
}
