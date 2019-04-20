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
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttClientID          = "SnipsHandler"
	mqttBrokerDefaultPort = 1883
	mqttBrokerAddrFormat  = "tcp://%s:%d"
	mqttBrokerKeepAlive   = 5
	mqttBrokerSubQOS      = 1
)

// Options contains the settings for initialisation of the client.
type Options struct {
	Name      string
	Host      string
	Port      uint16
	QOS       uint8
	KeepAlive uint8
}

// sanitiseOptions checks values of central options, filling them with defaults if empty/null
func (o *Options) sanitiseOptions() {
	if len(o.Name) == 0 {
		o.Name = mqttClientID
	}
	if len(o.Host) == 0 {
		o.Host = "127.0.0.1"
	}
	if o.Port == 0 {
		o.Port = mqttBrokerDefaultPort
	}
	if o.QOS == 0 {
		o.QOS = mqttBrokerSubQOS
	}
	if o.KeepAlive == 0 {
		o.KeepAlive = mqttBrokerKeepAlive
	}
}

// convertOptions converts snips.Options to mqtt.ClientOptions
func (o *Options) convertOptions() *mqtt.ClientOptions {
	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf(mqttBrokerAddrFormat, o.Host, o.Port))
	mqttOpts.SetClientID(o.Name)
	mqttOpts.SetCleanSession(true)
	mqttOpts.SetAutoReconnect(true)
	mqttOpts.SetKeepAlive(time.Second * time.Duration(o.KeepAlive))
	return mqttOpts
}
