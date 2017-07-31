package klessgo

// Context is probably good for something later on...
type Context struct {
	Info map[string]string
}

// Request is a generic inbound message with some headers and body
type Request struct {
	Headers map[string]string
	Body    []byte
}

// Response is a generic outbound message with some headers and body
type Response struct {
	Headers map[string]string
	Body    []byte
}

// KlessHandler interface that needs to be implemented by event handlers
type KlessHandler interface {
	Handler(c *Context, resp *Response, req *Request)
}
