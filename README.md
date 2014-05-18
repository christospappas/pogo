# Pogo

Make your API realtime.

## Usage
##### Start example server

`go run examples/main.go`

##### Create a websocket in a browser
`var websocket = new WebSocket("ws://localhost:8080");`

##### Subscribe to a channel
`websocket.send(JSON.stringify({event: "streama:subscribe", data: {channel: '/courses/23423423'}}));`

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


## Streaming REST protocol

This is a work in progress...

### Requests

#### AUTH

###### Request

    {
        seqId: 1,
        event: "auth", 
        data: {
            token: '7879a3aa5783ba49e89fb2ee3a27480c96c4e6ce'
        }
    }

###### Response

    {
        seqId: 1,
        status: 'ok',
        data: {
            expires: 1400382523
        }
    }


#### BIND

###### Request

    {
        seqId: 1,    
        event: "bind", 
        path: '/posts/1234'
    }
    
###### Response

    {
        seqId: 1,
        status: 'ok'
    }
    
---

    {
        event: "response",
        path: '/posts/1234',
        data: {
            ...
        }
    }

#### UNBIND

###### Request

    {
        seqId: 1,
        event: "unbind", 
        path: '/posts/1234'
    }
    
###### Response

    {
        seqId: 1,
        status: 'ok'
    }

#### GET

###### Request

    {
        seqId: 1,
        event: "get", 
        path: '/posts/1234'
    }    
    
###### Response

    {
        event: "streama:response",
        path: '/posts/1234',
        data: {
            ...
        }
    }

---

    {
        seqId: 1,
        status: 'ok'
    }
    

#### POST

###### Request

    {
        seqId: 1,
        event: "post",
        path: '/posts', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response

    {
        event: "response", 
        path: '/posts',
        status: 201,
        data: {
            ...
        }
    }
    
---

    {
        seqId: 1,
        status: 'ok'
    }
    
    
#### PUT

###### Request

    {
        seqId: 1,
        event: "put", 
        path: '/posts/1234', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response
    
    {
        event: "response", 
        path: '/posts',
        data: {
            ...
        }
    }

---

    {
        seqId: 1,
        status: 'ok'
    }
    
#### DELETE

###### Request

    {
        seqId: 1,
        event: "delete", 
        path: '/posts/1234', 
        data: {
            title: 'Post title',
            ...
        }
    }
    
###### Response


    {
        seqId: 1,
        status: 'ok'
    }
    