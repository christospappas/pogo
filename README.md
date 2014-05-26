# Pogo

Make your API realtime.

## Usage
##### Start example server

`go run examples/main.go`

##### Create a websocket in a browser
`var websocket = new WebSocket("ws://localhost:8080");`

##### Subscribe to a channel
`websocket.send(JSON.stringify({event: "channel:subscribe", channel: '/courses/23423423'}));`

##### Unsubscribe from a channel
`websocket.send(JSON.stringify({event: "channel:unsubscribe", channel: '/courses/23423423'}));`

##### Send an event to a channel
`websocket.send(JSON.stringify({event: "someEvent", channel: '/courses/23423423', data: {test: ''}}));`

##### Close websocket
`websocket.close()`

## Ideas

- MsgPack for binary transport (make it configurable)
- Namespaced/seperate handlers
- Compress json format to single character keys
- Autobind on any REST action, unless explicitly disabled
- Optional TTL
- Side benefit of providing API CDN
- Security filters such as ip blacklisting/rate limiting
- Web Interface for configuring REST endpoints
- Libraries to integrate push into apps
- Clustering support


## Streaming API

This is a work in progress...

### Requests

#### AUTH

###### Request

    {
        id: 1,
        event: "auth", 
        data: {
            token: '7879a3aa5783ba49e89fb2ee3a27480c96c4e6ce'
        }
    }

###### Response

    {
        id: 1,
        status: 'ok',
        data: {
            expires: 1400382523
        }
    }


#### BIND

###### Request

    {
        id: 1,    
        event: "bind", 
        uri: '/posts/1234'
    }
    
###### Response

    {
        id: 1,
        status: 'ok'
    }
    
---

    {
        event: "update",
        uri: '/posts/1234',
        data: {
            ...
        }
    }

#### UNBIND

###### Request

    {
        id: 1,
        event: "unbind", 
        uri: '/posts/1234'
    }
    
###### Response

    {
        id: 1,
        status: 'ok'
    }

#### LIST

###### Request

    {
        id: 1,
        event: "list", 
        uri: '/posts'
    }    
    
###### Response

    {
        id: 1,
        status: 'ok'
        data: {
            ...
        }
    }


#### RETRIEVE

###### Request

    {
        id: 1,
        event: "retrieve", 
        uri: '/posts/1234'
    }    
    
###### Response

    {
        id: 1,
        status: 'ok'
        data: {
            ...
        }
    }

#### CREATE

###### Request

    {
        id: 1,
        event: "create",
        uri: '/posts', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response

    {
        id: 1,
        status: 'ok',
        data: {
            ...
        }
    }
    

#### UPDATE

###### Request

    {
        id: 1,
        event: "update", 
        uri: '/posts/1234', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response
    
    {
        id: 1,
        status: 'ok'
    }
    
#### DELETE

###### Request

    {
        id: 1,
        event: "delete", 
        uri: '/posts/1234', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response


    {
        id: 1,
        status: 'ok'
    }
    
