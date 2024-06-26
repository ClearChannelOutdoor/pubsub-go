# pubsub-go/v2

This is a Go module that provides a convenience wrapper for publishing and receiving messages via Google Cloud Pub/Sub.

This library allows for messages to be received via channel, and to rapidly create multiple subscriptions using a single function call.

## Installation

```bash
go get github.com/clearchanneloutdoor/pubsub-go/v2
```

## Usage

### Create a new PubSub

```go
import psb "github.com/clearchanneloutdoor/pubsub-go/pkg"

func main() {
  client, err := psb.NewPubSub(
    context.Background(),
    pb.Options("<projectID>"))
  if err != nil {
    panic(err)
  }
}
```

#### Provide additional configuration

```go
import (
  "cloud.google.com/go/pubsub"
  "google.golang.org/api/option"
  psb "github.com/clearchanneloutdoor/pubsub-go/pkg"
)

func main() {
  opts := pb.Options("<projectID>").
    SetClientOptions(
      option.WithGRPCConnectionPool(100),
      option.WithCredentialsFile("<path to credentials file>")
  client, err := psb.NewPubSub(context.Background(), opts)
  if err != nil {
    panic(err)
  }
}
```

### Create a Topic

```go
package main

import (
  "cloud.google.com/go/pubsub"
  psb "github.com/clearchanneloutdoor/pubsub-go/pkg"
)

func main() {
  // Initialize new pubsub-go PubSub 
  client, err := psb.NewPubSub(
    context.Background(), 
    pb.Options("<project ID>"))
  if err != nil {
    panic(err)
  }
  defer func() {
    if err := client.Close(); err != nil {
      panic(err)
    }
  }()

  // create the topic 
  if err := client.CreateTopic(
    "<topic ID>", 
    pubsub.TopicConfig{
      RetentionDuration: time.Hour * 24 * time.Duration(3),
    }); err != nil {
    panic(err)
  }

  // send a message to the topic
}
```

### Create Multiple Subscriptions with Filters for a Topic

```go
import (
  "cloud.google.com/go/pubsub"
  psb "github.com/clearchanneloutdoor/pubsub-go/pkg"
)

func main() {
  // Initialize new pubsub-go PubSub

  // Define a map for a list of Subscription Names and their (optional) filter definitions
  subs := map[string]string{
    "topic-sub": "",
    "topic-sub-ca": "attributes.region = \"CA\"",
    "topic-sub-nv": "attributes.region = \"NV\"",
    "topic-sub-tx": "attributes.region = \"TX\"",
  }

  // Create subscriptions
  if err := client.CreateSubscriptions("<topic ID>", subs, pubsub.SubscriptionConfig{
    EnableMessageOrdering: true,
    RetainAckedMessages:   false,
  }); err != nil {
    panic(err)
  }
}
```

### Publish Message

The Publish function can be used to publish a string or an object (which is serialized as JSON).

```go
type Example struct {
  CoolNumber int `json:"coolNumber"`
  Message string `json:"message"`
}

func main() {
  // Initialize new PubSub client

  // Publish a string message
  if err := client.Publish("<topic ID>", "hello world"); err != nil {
    panic(err)
  }

  // Publish an object as JSON
  e := Example{
    CoolNumber: 13,
    Message: "hello world",
  }
  if err := client.Publish("<topic ID>", e); err != nil {
    panic(err)
  }
}
```

#### Publish Messages in Other Formats

Publish can also be used to send messages serialized in other formats as well, if you serialize the payload prior to calling the function and pass in the serialized `[]byte` slice:

```go
type Example struct {
  Message string `yaml:"message"`
}

func main() {
  // Initialize new PubSub client

  // Serialize as yaml
  sy, err := yaml.Marshal(Example{Message: "hello world"})
  if err != nil {
    panic(err)
  }

  // Publish serialized yaml
  if err := client.Publish("<topic ID>", sy); err != nil {
    panic(err)
  }
}
```

#### Publish Messages with PublishSettings

PublishSettings can be specified in options used when creating the PubSub client. The settings are then used to control the behavior of the publication.

```go
opts := psb.Options("<project ID>").
  SetPublishSettings(pubsub.PublishSettings{
    DelayThreshold:  (10 * time.Millisecond),
    CountThreshold:  1000,
    ByteThreshold:   1000000,
    Timeout:         (10 * time.Second),
  })
client, err := psb.NewPubSub(context.Background(), opts)

if err := client.Publish("<topic ID>", "hello world"); err != nil {
  panic(err)
}
```

### Receive Messages

```go
func main() {
  // Initialize new PubSub client

  // Create a channel for receiving messages
  messages := make(chan *pubsub.Message)
  go func() {
    // loop to continuously receive
    for {
      msg := <-messages

      // Process message
      fmt.Printf("Received a message: %s\n", string(msg.Data))

      // Acknowledge message
      msg.Ack()
    }
  }()

  // create a subscription
  if err := client.CreateSubscription("<topic ID>", "<subscription ID>", ""); err != nil {
    panic(err)
  }

  // receive messages
  if err := client.Receive("<subscription ID>", messages); err != nil {
    panic(err)
  }
}
```

#### Receive Messages with ReceiveSettings

ReceiveSettings can be specified in options used when creating the PubSub client. The settings are then used to control the behavior of the subscription.

```go
opts := psb.Options("<project ID>").
  SetReceiveSettings(pubsub.ReceiveSettings{
   MaxExtension:           (15 * time.Second),
   MaxOutstandingMessages: 1000,
   NumGoroutines:          10,
  })
client, err := psb.NewPubSub(context.Background(), opts)

if err := ps.Receive("<subscription ID>", messages); err != nil {
  panic(err)
}
```

## Running GCP PubSub Locally

### GCP SDK

If you have installed the GCP SDK, you can run the following command to start the emulator.

```bash
gcloud beta emulators pubsub start --project=example-project --host-port=localhost:8085 
```

Then, export the following environment variables to point to the emulator.

```bash
export PUBSUB_EMULATOR_HOST=localhost:8085
export PUBSUB_PROJECT_ID=example-project
```

#### Run the Example

After starting the emulator, in the same terminal window where you have exported the environment variables, run the following command to start the example.

```bash
cd v2
go run examples/pubsub.go
```

### Docker

Google publishes an emulator for GCP PubSub, so you can run it locally. This repo includes a script that will spin up a docker container with the emulator started so running a local dev environment is easier.

Huge shout out to [@anguillanneuf](https://github.com/anguillanneuf), who wrote a [blog post](https://medium.com/google-cloud/things-i-wish-i-knew-about-pub-sub-part-3-b8947b49224b) that made building this script much easier.

#### Dependencies

- Docker
- Openssl

#### Start

Open a terminal and navigate to this project's directory. once there
run the following command...

```bash
sh ./local-pubsub.sh
```

This will take you through a wizard to get all information necessary to start a Project, Topic, and Subscription.

There is an `export` command that is out put once the script has completed that you'll need to copy and paste across all your open terminal windows. The reason being is that there is no way to set the GCP PubSub endpoint directly in your application; the GCP libray looks to an environment to know which endpoint to use.

#### The -m Option

Sending a message directly to your queue using something like `curl` isn't clearly documented. Rather than sending a json payload, you send a base64 encoded string of the message data.

Execute the following command to go throw a wizard that will output a properly formatted request for you.

```bash
sh ./local-pubsub.sh -m
```
