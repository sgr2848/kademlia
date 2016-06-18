package dht

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetMsg(t *testing.T) {
	netMsgInit()
	var conn bytes.Buffer

	node := newNode(&NetworkNode{})
	id, _ := newID()
	node.ID = id
	node.Port = 3000
	node.IP = net.ParseIP("127.0.0.1")

	msg := &message{}
	msg.Type = messageTypeQueryFindNode
	msg.Receiver = node.NetworkNode
	msg.Data = &queryDataFindNode{
		Target: id,
	}

	serialized, err := serializeMessage(msg)
	if err != nil {
		panic(err)
	}

	conn.Write(serialized)

	deserialized, err := deserializeMessage(&conn)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, msg, deserialized)
}