package datapoint

import (
	"context"
	"gobaos/pkg/baos"
	"gobaos/pkg/knxtypes"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type sensor struct {
	id    baos.DatapointID
	tc    knxtypes.TypeCode
	topic string
}

type Processor struct {
	topics  map[string]*sensor
	sensors map[baos.DatapointID]*sensor
	client  mqtt.Client
	worker  *baos.Worker
	log     Logger
}

func (p *Processor) Handle(id baos.DatapointID, data []byte) error {
	s, ok := p.sensors[id]
	if !ok {
		return nil
	}
	dpt := knxtypes.MustNewDPT(s.tc)
	if err := dpt.UnmarshalBinary(data); err != nil {
		return err
	}
	token := p.client.Publish(s.topic, 0, false, dpt.Get())
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (p *Processor) AddSensor(id baos.DatapointID, desc baos.DatapointDescription, pat pattern) error {
	result := sensor{id: id, tc: desc.DPT, topic: pat.apply(id.String())}
	p.sensors[id] = &result
	return nil
}

func (p *Processor) AddController(id baos.DatapointID, desc baos.DatapointDescription, pat pattern) error {
	result := sensor{id: id, tc: desc.DPT, topic: pat.apply(id.String())}
	p.sensors[id] = &result
	p.topics[pat.apply(id.String())] = &result
	return nil
}

func (p *Processor) handleMessage(_ mqtt.Client, msg mqtt.Message) {
	s, ok := p.topics[msg.Topic()]
	if !ok {
		return
	}
	dpt := knxtypes.MustNewDPT(s.tc)
	if err := dpt.Set(string(msg.Payload())); err != nil {
		p.log.Errorf("dpt.set: %s", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.worker.Set(ctx, s.id, dpt); err != nil {
		p.log.Errorf("worker.set: %s", err.Error())
	}
	p.log.Infof("processor.handleMessage: topic=%s, value=%s", msg.Topic(), dpt.Get())
}

func NewProcessor(client mqtt.Client, worker *baos.Worker, log Logger) Processor {
	result := Processor{topics: make(map[string]*sensor), sensors: make(map[baos.DatapointID]*sensor),
		client: client, worker: worker, log: log}
	return result
}
