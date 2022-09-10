package meanstoanend_test

import (
	"bytes"
	"testing"

	meanstoanend "github.com/nint8835/protohackers/2-means-to-an-end"
)

type exampleInput struct {
	data            []byte
	expectedMessage meanstoanend.Message
	expectedResp    []byte
}

type CloseableBuffer struct {
	*bytes.Buffer
	Closed bool
}

func (buffer *CloseableBuffer) Close() error {
	buffer.Closed = true
	return nil
}

var exampleInputs = []exampleInput{
	{
		data:            []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65},
		expectedMessage: meanstoanend.Message{Type: meanstoanend.MessageTypeInsert, Arg1: 12345, Arg2: 101},
		expectedResp:    []byte{},
	},
	{
		data:            []byte{0x49, 0x00, 0x00, 0x30, 0x3a, 0x00, 0x00, 0x00, 0x66},
		expectedMessage: meanstoanend.Message{Type: meanstoanend.MessageTypeInsert, Arg1: 12346, Arg2: 102},
		expectedResp:    []byte{},
	},
	{
		data:            []byte{0x49, 0x00, 0x00, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x64},
		expectedMessage: meanstoanend.Message{Type: meanstoanend.MessageTypeInsert, Arg1: 12347, Arg2: 100},
		expectedResp:    []byte{},
	},
	{
		data:            []byte{0x49, 0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x05},
		expectedMessage: meanstoanend.Message{Type: meanstoanend.MessageTypeInsert, Arg1: 40960, Arg2: 5},
		expectedResp:    []byte{},
	},
	{
		data:            []byte{0x51, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x40, 0x00},
		expectedMessage: meanstoanend.Message{Type: meanstoanend.MessageTypeQuery, Arg1: 12288, Arg2: 16384},
		expectedResp:    []byte{0x00, 0x00, 0x00, 0x65},
	},
}

func TestConnection_ReadMessage(t *testing.T) {
	for index, example := range exampleInputs {
		t.Logf("testing example %d", index+1)

		connection := meanstoanend.Connection{
			Connection: &CloseableBuffer{bytes.NewBuffer(example.data), false},
		}

		message, err := connection.ReadMessage()
		if err != nil {
			t.Errorf("got unexpected error: %s", err)
			return
		}

		if message.Type != example.expectedMessage.Type {
			t.Errorf("got unexpected message type: %v", message.Type)
			return
		}
		if message.Arg1 != example.expectedMessage.Arg1 {
			t.Errorf("got unexpected arg1: %v", message.Arg1)
			return
		}
		if message.Arg2 != example.expectedMessage.Arg2 {
			t.Errorf("got unexpected arg2: %v", message.Arg2)
			return
		}
	}
}

func TestConnection_HandleMessage(t *testing.T) {
	connection := meanstoanend.Connection{
		Connection: &CloseableBuffer{bytes.NewBuffer([]byte{}), false},
	}

	for index, example := range exampleInputs {
		t.Logf("testing example %d", index+1)

		resp := connection.HandleMessage(example.expectedMessage)
		if !bytes.Equal(resp, example.expectedResp) {
			t.Errorf("got unexpected response: %v", resp)
		}
	}
}
