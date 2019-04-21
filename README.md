# go-snips

The go-snips library provides a simple, minimal interface allowing developers to register a handler function for each intent parsed by the [Snips](https://snips.ai/) voice assistant.

Documentation for setting up the Snips voice assistant can be found [here](https://docs.snips.ai/getting-started).

### Download

    go get github.com/sturd/go-snips


### Usage

```go
import "github.com/sturd/go-snips"

func main() {
    snipsOpts := snips.Options{
        Name: "MySnipsHandler",
        Host: "[SNIPS_MQTT_ADDRESS]",
    }
    snips, err := snips.NewClient(snipsOpts)
    if err != nil {
        panic(err)
    }
    snips.SubscribeIntentHandler("mySnipsIntent", intentHandler)
}

func intentHandler(Client mqtt.Client, Intent snips.IntentMessage) {
    // Handle the parsed intent
}
```