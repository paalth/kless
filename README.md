# kless
`kless` - a native &amp; naive event-handler framework for Kubernetes

# Description 
`kless` is a simple serverless framework that supports push & pull events with event handlers written in mainly typed languages.

It is a native Kubernetes framework as it attempts to use as many existing Kubernetes features as possible, and it is a naive framework as it currently goes about it's business in the simplest possible way that could conceivably work.

# Concepts
Currently event handlers can be written in the following languages:
- Go
- Java

However, event handlers can be written in any language that has an implemented "**event handler builder**". Event Handler Builders are responsible for compiling the event handler source code and turning them into containers that can be deployed onto Kubernetes. Event Handler Builders are themselves just containers registered with `kless` so in the future any other language can be supported. 

`kless` event handlers have "**frontends**". Frontends are sidecar containers that either pull events from event sources or are targets that receive events pushed to the handler. Frontends are associated with configuration information that applies to the event source. 

Frontends also have a "**frontend type**". Frontend Types are containers that implement the frontend, and they are registered with `kless` so in the future any other event source can be supported.

Currently the following push methods are supported:

- HTTP/HTTPS

Currently the following pull methods are supported:

- NATS
- Kafka
- RabbitMQ

`kless` event handlers implement a language-specific interface and the signature of the event handler consists of a **context**, a **request** and a **response**. 

The **context** is a map of strings that contain information about the event handler.

The **request** and **response** are both simple objects that have two fields. The first field is a map of string that contains the request/response headers, and the second field is a byte array containing the request/response body.

# Usage

To create an event handler from the CLI you'd run a command of the following format:

```kless create handler -e <event handler name> -b <event handler builder> -f <frontend> <event handler source code>```

For example to create an event handler implemented in Go from a local file with events pushed to it over HTTP you'd run:

```kubeless create handler -e go-http-handler1 -b go-1.7.4 -f http-local -s EventHandler.go```

Similary, if you'd want to create an event handler implemented in Java from a local file that would pull events from Kafka:

```kubeless create handler -e java-kafka-handler1 -b java-8u111 -f kafka-local -s EventHandler.java```

# Examples

The following is an example of the simples possible event handler written in Go:

```
package eventhandler

import (
	"fmt"

	kl "github.com/paalth/kless/pkg/interface/klessgo"
)

type EventHandler struct {
}

func (t EventHandler) Handler(c *kl.Context, resp *kl.Response, req *kl.Request) {
	fmt.Printf("Inside funcHandler\n")
}
```

The following is an example of the simples possible event handler written in Java:

```
package io.kless;

class EventHandler1 implements EventHandlerInterface {

    public Response eventHandler(Context context, Request req) {
        System.out.println("Inside event handler...");

        return null;
    }

}
```

The following is an example of a simplistic event handler written in Go that dumps incoming HTTP PUT or POST requests into a PostgreSQL table:

```
package eventhandler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/paalth/kless/pkg/interface/klessgo"
)

const (
	HOST     = "10.245.1.2"
	PORT     = 5432
	USERNAME = "postgres"
	PASSWORD = "postgres"
	DB_NAME  = "postgres"
)

// EventHandler is an example Kless event handler
type EventHandler struct {
}

// Handler stores the information in the incoming request in a database
func (t EventHandler) Handler(c *klessgo.Context, resp *klessgo.Response, req *klessgo.Request) {
	fmt.Printf("Inside Event Handler\n")

	if req.Headers["kless-method"] == "PUT" || req.Headers["kless-method"] == "POST" {
		fmt.Println("Writing request to PostgreSQL DB")

		insertIntoPostgreSQLDB(HOST, PORT, USERNAME, PASSWORD, DB_NAME, req.Headers, string(req.Body), c.Info["name"], c.Info["namespace"], c.Info["version"])

		resp.Body = []byte("Request written to DB\n")
	} else {
		resp.Body = []byte("Request received\n")
	}

	fmt.Println("Event handler processing complete")
}

func insertIntoPostgreSQLDB(host string, port int, username string, password string, dbname string, headers map[string]string, body string, eventHandler string, namespace string, version string) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)

	db, err := sql.Open("postgres", dbinfo)
	if nil != err {
		panic(err)
	}
	defer db.Close()

	creationDate := time.Now().Format(time.RFC3339Nano)

	fmt.Println("Inserting values into KLESS_EVENT_BODY")

	var id int
	err = db.QueryRow("INSERT INTO KLESS_EVENT_BODY(EVENT_HANDLER,NAMESPACE,VERSION,BODY,CREATION_DATE) VALUES($1,$2,$3,$4,$5) returning ID;", eventHandler, namespace, version, body, creationDate).Scan(&id)
	if nil != err {
		log.Fatal(err)
	}

	fmt.Println("Inserting values into KLESS_EVENT_HEADER")

	for k, v := range headers {
		var hid int

		err = db.QueryRow("INSERT INTO KLESS_EVENT_HEADER(HEADER_NAME,HEADER_VALUE,CREATION_DATE,BODY_ID) VALUES($1,$2,$3,$4) returning ID;", k, v, creationDate, id).Scan(&hid)
		if nil != err {
			log.Fatal(err)
		}
	}
}

```

# Status

kless is currently a work-in-progress/POC not intended for serious usage...
