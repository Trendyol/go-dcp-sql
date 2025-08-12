package couchbase

import (
	"time"

	"github.com/Trendyol/go-dcp/tracing"
)

type Event struct {
	tracing.ListenerTrace
	CollectionName string
	EventTime      time.Time
	Key            []byte
	Value          []byte
	Cas            uint64
	VbID           uint16
	IsDeleted      bool
	IsExpired      bool
	IsMutated      bool
}

func NewDeleteEvent(listenerTrace tracing.ListenerTrace,
	key []byte, value []byte,
	collectionName string, eventTime time.Time, cas uint64, vbID uint16,
) Event {
	return Event{
		ListenerTrace:  listenerTrace,
		Key:            key,
		Value:          value,
		IsDeleted:      true,
		CollectionName: collectionName,
		EventTime:      eventTime,
		Cas:            cas,
		VbID:           vbID,
	}
}

func NewExpireEvent(listenerTrace tracing.ListenerTrace,
	key []byte, value []byte,
	collectionName string, eventTime time.Time, cas uint64, vbID uint16,
) Event {
	return Event{
		ListenerTrace:  listenerTrace,
		Key:            key,
		Value:          value,
		IsExpired:      true,
		CollectionName: collectionName,
		EventTime:      eventTime,
		Cas:            cas,
		VbID:           vbID,
	}
}

func NewMutateEvent(listenerTrace tracing.ListenerTrace,
	key []byte, value []byte,
	collectionName string, eventTime time.Time, cas uint64, vbID uint16,
) Event {
	return Event{
		ListenerTrace:  listenerTrace,
		Key:            key,
		Value:          value,
		IsMutated:      true,
		CollectionName: collectionName,
		EventTime:      eventTime,
		Cas:            cas,
		VbID:           vbID,
	}
}
