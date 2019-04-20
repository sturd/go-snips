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
	"fmt"

	"errors"

	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type snipsIntents map[string]IntentHandler

// Client is the object which listens to events fired by the Snips ASR, parses the intents
// and passes them to the relevant subscribed handler.
type Client struct {
	mqtt    mqtt.Client
	intents snipsIntents
}

// MQTT topics which Snips ASR publishes messages on
const (
	mqttIntentTopic  = "hermes/nlu/intentParsed"
	mqttHotwordTopic = "hermes/hotword/#"
)

// NewClient creates a new Snips client, establishing a connection to the MQTT broker.
func NewClient(options Options) (*Client, error) {

	client := buildClient(options)
	if token := client.mqtt.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := client.mqtt.Subscribe(mqttIntentTopic, mqttBrokerSubQOS, client.handleIntentParsed); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := client.mqtt.Subscribe(mqttHotwordTopic, mqttBrokerSubQOS, client.handleHotword); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

// buildClient sanitises provided options, converts them and creates the mqtt.Client
func buildClient(options Options) *Client {
	client := &Client{}
	options.sanitiseOptions()
	client.mqtt = mqtt.NewClient(options.convertOptions())
	client.intents = make(snipsIntents)
	return client
}

// Function is subscribed to mqttIntentTopic during NewClient().  The raw, JSON data is unmarshaled
// and passed to the callback which is subscribed to its particular intent name.
func (c *Client) handleIntentParsed(client mqtt.Client, msg mqtt.Message) {
	var parsedIntent IntentMessage
	err := json.Unmarshal(msg.Payload(), &parsedIntent)
	if err != nil {
		panic(err)
	}

	intentName := getIntentName(parsedIntent.Intent.Name)
	handler := c.intents[intentName]
	if handler != nil {
		handler(client, parsedIntent)
	}
}

// Function to log when hotword has been activated/deactivated
func (c *Client) handleHotword(client mqtt.Client, msg mqtt.Message) {
	fmt.Println(string(msg.Topic()))
}

// Provide a handler specific to a parsed intent from Snips' NLU
func (c *Client) SubscribeIntentHandler(name string, handler IntentHandler) error {
	if c.intents[name] != nil {
		return errors.New(fmt.Sprintf("handler \"%s\" already registered", name))
	}
	c.intents[name] = handler
	return nil
}

// If intents have been forked on the Snips console, intents
// will be prefixed by the developer/user's name.  This function
// will remove the username and delimiter, if it's there.
func getIntentName(name string) string {
	names := strings.Split(name, ":")
	return names[len(names)-1]
}
