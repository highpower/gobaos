package mqttbridge

import (
	"context"
	"fmt"
	"gobaos/internal/pattern"
	"gobaos/pkg/baos"
	"gobaos/pkg/knxtypes"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DatapointsSettings interface {
	Pattern() string
	UpdatablePattern() string
	ReadOnlyContains(id baos.DatapointID) bool
	UpdatableContains(id baos.DatapointID) bool
}

type Settings interface {
	MqttSettings
	DatapointsSettings
}

type Handler struct {
	client    mqtt.Client
	worker    *baos.Worker
	log       Logger
	readOnly  map[baos.DatapointID]readOnlyItem
	updatable map[string]updatableItem
}

type BridgeOption interface {
	apply(b *Handler)
}

type withLoggerOption struct {
	log Logger
}

type readOnlyItem struct {
	topic string
	code  knxtypes.TypeCode
}

type updatableItem struct {
	id   baos.DatapointID
	code knxtypes.TypeCode
}

func (h *Handler) onConnect(client mqtt.Client) {
	h.log.Debugf("handler.onConnect: entering")
	defer h.log.Debugf("handler.onConnect: exiting")
	topics := make(map[string]byte, len(h.updatable))
	for topic := range h.updatable {
		topics[topic] = 0
	}
	if token := client.SubscribeMultiple(topics, h.onMessage); token.Wait() && token.Error() != nil {
		h.log.Errorf("bridge.onConnect: failed to subscribe to topics: %s", token.Error().Error())
	}
}

func (h *Handler) onMessage(_ mqtt.Client, message mqtt.Message) {
	if item, ok := h.updatable[message.Topic()]; ok {
		dpt, err := knxtypes.NewDPT(item.code)
		if err != nil {
			h.log.Errorf("bridge.onMessage: can not create dpt[%s]", err.Error())
			return
		}
		value := string(message.Payload())
		if err := dpt.Set(value); err != nil {
			h.log.Errorf("bridge.onMessage: can not parse message[%s]", err.Error())
			return
		}
		ctx := context.Background()
		if err := h.worker.Set(ctx, item.id, dpt); err != nil {
			h.log.Errorf("bridge.onMessage: can not set dpt value[%s]", err.Error())
		}
		h.log.Infof("handler.onMessage: dpt with id %d set to '%s'", item.id, value)
	}
}

func (h *Handler) OnIndication(id baos.DatapointID, data []byte) error {
	if item, ok := h.readOnly[id]; ok {
		dpt, err := knxtypes.NewDPT(item.code)
		if err != nil {
			return fmt.Errorf("bridge.onIndication: can not create dpt[%w]", err)
		}
		if err := dpt.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("bridge.onIndication: can not unmarshal data[%w]", err)
		}
		value := dpt.Get()
		h.client.Publish(item.topic, 0, false, value)
		h.log.Infof("handler.onIndication: message with value '%s' of dpt %d sent to '%s'",
			id, value, item.topic)
	}
	return nil
}

func (o withLoggerOption) apply(b *Handler) {
	b.log = o.log
}

func NewHandler(worker *baos.Worker, s Settings, opts ...BridgeOption) (*Handler, error) {
	ctx := context.Background()
	desc, err := worker.AllDescriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("bridge.new: can not get descriptions: %w", err)
	}
	result := Handler{
		worker:    worker,
		log:       nopLogger{},
		readOnly:  map[baos.DatapointID]readOnlyItem{},
		updatable: map[string]updatableItem{},
	}
	for _, opt := range opts {
		opt.apply(&result)
	}
	if err := result.setupDescriptions(s, desc); err != nil {
		return nil, fmt.Errorf("bridge.new: can not setup descriptions: %w", err)
	}
	client, err := connect(s, result.onConnect)
	if err != nil {
		return nil, fmt.Errorf("bridge.new: can not connect to mqtt: %w", err)
	}
	result.client = client
	return &result, nil
}

func (h *Handler) setupDescriptions(s DatapointsSettings, desc []baos.DatapointDescription) error {
	h.log.Debugf("setupDescriptions: entering")
	defer h.log.Debugf("setupDescriptions: exiting")
	p, err := pattern.Parse(s.Pattern())
	if err != nil {
		return fmt.Errorf("can not parse pattern: %w", err)
	}
	up, err := pattern.Parse(s.UpdatablePattern())
	if err != nil {
		return fmt.Errorf("can not parse updatable pattern: %w", err)
	}
	for _, desc := range desc {
		switch {
		case s.ReadOnlyContains(desc.ID):
			h.readOnly[desc.ID] = readOnlyItem{topic: p.Apply(strconv.Itoa(int(desc.ID))), code: desc.Code}
		case s.UpdatableContains(desc.ID):
			topic := up.Apply(strconv.Itoa(int(desc.ID)))
			h.updatable[topic] = updatableItem{id: desc.ID, code: desc.Code}
		}
	}
	return nil
}

func WithLogger(log Logger) BridgeOption {
	return withLoggerOption{log: log}
}
