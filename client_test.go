/*******************************************************************************

The MIT License (MIT)

Copyright (c) 2019 Craig Sturdy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*******************************************************************************/
package snips

import (
	"io/ioutil"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/magiconair/properties/assert"
)

type testMessage struct {
	Data   []byte
	Parsed IntentMessage
}

func (t *testMessage) Duplicate() bool   { return false }
func (t *testMessage) Qos() byte         { return 0 }
func (t *testMessage) Retained() bool    { return false }
func (t *testMessage) Topic() string     { return "test/message" }
func (t *testMessage) MessageID() uint16 { return 0 }
func (t *testMessage) Payload() []byte   { return t.Data }
func (t *testMessage) Ack()              {}

func (t *testMessage) handleIntentParsed(client mqtt.Client, message IntentMessage) {
	t.Parsed = message
}

func TestClient_getIntentName(t *testing.T) {
	ExpectedIntent := "playMusic"
	PrefixedIntent := "user_name:playMusic"

	// Prove that a non-prefixed string returns the same
	Result := getIntentName(ExpectedIntent)
	if ExpectedIntent != Result {
		t.Fatalf("Expected: %s -- Got: %s", ExpectedIntent, Result)
	}

	Result = getIntentName(PrefixedIntent)
	if ExpectedIntent != Result {
		t.Fatalf("Expected: %s -- Got: %s", ExpectedIntent, Result)
	}
}

func TestClient_SubscribeIntentHandler(t *testing.T) {
	client := buildClient(Options{})
	assert.Equal(t, len(client.intents), 0)

	err := client.SubscribeIntentHandler("testIntent", func(client mqtt.Client, message IntentMessage) {})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(client.intents), 1)

	err = client.SubscribeIntentHandler("testIntent", func(client mqtt.Client, message IntentMessage) {})
	assert.Equal(t, err.Error(), "handler \"testIntent\" already registered")
	assert.Equal(t, len(client.intents), 1)
}

func TestClient_handleIntentParsed(t *testing.T) {
	client := buildClient(Options{})

	var err error
	tm := testMessage{}
	tm.Data, err = ioutil.ReadFile("test-data/intent-parsed.json")
	assert.Equal(t, err, nil)

	if err == nil {
		err = client.SubscribeIntentHandler("testIntent", tm.handleIntentParsed)
		client.handleIntentParsed(client.mqtt, &tm)

		message := &tm.Parsed
		assert.Equal(t, message.Input, "this here is a test")
		assert.Equal(t, message.Intent.Name, "testIntent")
		assert.Equal(t, message.Intent.Confidence, 0.9744988)
		assert.Equal(t, len(message.Slots), 1)

		slot := &message.Slots[0]
		assert.Equal(t, slot.RawValue, "test")
		assert.Equal(t, slot.Value.Kind, "Custom")
		assert.Equal(t, slot.Value.Value, "test")
		assert.Equal(t, slot.Range.Start, int32(15))
		assert.Equal(t, slot.Range.End, int32(19))
		assert.Equal(t, slot.Entity, "snips/default--intentMode")
		assert.Equal(t, slot.Name, "intentMode")

		sliceValue := message.Input[slot.Range.Start:slot.Range.End]
		assert.Equal(t, sliceValue, slot.RawValue)
	}
}
