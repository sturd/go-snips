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
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IntentHandler func(mqtt.Client, IntentMessage)

type Intent struct {
	Name       string  `json:"intentName"`
	Confidence float64 `json:"confidenceScore"`
}

type IntentMessage struct {
	Id     string `json:"id"`
	Input  string `json:"input"`
	Intent Intent `json:"intent"`
	Slots  []Slot `json:"slots"`
}

func (i *IntentMessage) LoadFromJSON(Data []byte) {
	Err := json.Unmarshal(Data, i)
	if Err != nil {
		panic(Err)
	}
}
