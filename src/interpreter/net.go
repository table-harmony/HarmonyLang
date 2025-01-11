package interpreter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RequestType represents the type of an HTTP request
type RequestType struct{}

func (RequestType) String() string { return "request" }
func (r RequestType) Equals(other Type) bool {
	_, ok := other.(RequestType)
	return ok
}
func (r RequestType) DefaultValue() Value { return NewRequest() }

// ResponseType represents the type of an HTTP response
type ResponseType struct{}

func (ResponseType) String() string { return "response" }
func (r ResponseType) Equals(other Type) bool {
	_, ok := other.(ResponseType)
	return ok
}
func (r ResponseType) DefaultValue() Value { return NewResponse() }

// Request represents an HTTP request with methods
type Request struct {
	Method      string
	Path        string
	Query       map[string]string
	Headers     map[string]string
	Body        string
	Params      map[string]string
	ContentType string
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
		Params:  make(map[string]string),
	}
}

// Request implements Value interface
func (Request) Type() Type     { return RequestType{} }
func (r Request) Clone() Value { return r }
func (r Request) String() string {
	return fmt.Sprintf("Request{method: %s, path: %s}", r.Method, r.Path)
}

// Response represents an HTTP response with methods
type Response struct {
	StatusCode  int
	Headers     map[string]string
	Body        strings.Builder
	ContentType string
	Writer      http.ResponseWriter
	Methods     map[string]Function
}

func NewResponse() *Response {
	res := &Response{
		StatusCode: 200,
		Headers:    make(map[string]string),
		Methods:    make(map[string]Function),
	}
	res.init_methods()
	return res
}

func (res *Response) init_methods() {
	res.Methods["status"] = NewNativeFunction(
		func(args ...Value) Value {
			code := int(args[0].(Number).Value())
			res.StatusCode = code
			return res
		},
		[]Type{PrimitiveType{NumberType}},
		ResponseType{},
	)

	res.Methods["setHeader"] = NewNativeFunction(
		func(args ...Value) Value {
			name := args[0].(String).Value()
			value := args[1].(String).Value()
			res.Headers[name] = value
			return res
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{StringType}},
		ResponseType{},
	)

	res.Methods["write"] = NewNativeFunction(
		func(args ...Value) Value {
			content := args[0].(String).Value()
			res.Body.WriteString(content)
			return res
		},
		[]Type{PrimitiveType{StringType}},
		ResponseType{},
	)

	res.Methods["json"] = NewNativeFunction(
		func(args ...Value) Value {
			data := args[0]
			res.ContentType = "application/json"
			jsonBytes, err := json.Marshal(convert_to_native(data))
			if err != nil {
				panic(err)
			}
			res.Body.Write(jsonBytes)
			return res
		},
		[]Type{PrimitiveType{AnyType}},
		ResponseType{},
	)

	res.Methods["html"] = NewNativeFunction(
		func(args ...Value) Value {
			content := args[0].(String).Value()
			res.ContentType = "text/html"
			res.Body.WriteString(content)
			return res
		},
		[]Type{PrimitiveType{StringType}},
		ResponseType{},
	)
}

// Response implements Value interface
func (Response) Type() Type     { return ResponseType{} }
func (r Response) Clone() Value { return r }
func (r Response) String() string {
	return fmt.Sprintf("Response{status: %d}", r.StatusCode)
}

type ServerType struct {
}

func NewServerType() ServerType {
	return ServerType{}
}

// ServerType implements the Type interface
func (ServerType) String() string { return "server" }
func (s ServerType) Equals(other Type) bool {
	_, ok := other.(ServerType)
	return ok
}
func (s ServerType) DefaultValue() Value { return NewNil() }

// Server represents an HTTP server with improved routing
type Server struct {
	routes  map[string]map[string]Function
	methods map[string]Function
}

func NewServer() *Server {
	server := &Server{
		routes:  make(map[string]map[string]Function),
		methods: make(map[string]Function),
	}
	server.init_methods()
	return server
}

// Server implements Value interface
func (Server) Type() Type     { return ServerType{} }
func (s Server) Clone() Value { return s }
func (s Server) String() string {
	return fmt.Sprintf("Server{routes: %d}", len(s.routes))
}

func (s *Server) init_methods() {
	// GET method
	s.methods["get"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			handler := args[1].(Function)

			if s.routes[path] == nil {
				s.routes[path] = make(map[string]Function)
			}
			s.routes[path]["GET"] = handler
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{AnyType}},
		PrimitiveType{NilType},
	)

	// POST method
	s.methods["post"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			handler := args[1].(Function)

			if s.routes[path] == nil {
				s.routes[path] = make(map[string]Function)
			}
			s.routes[path]["POST"] = handler
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{AnyType}},
		PrimitiveType{NilType},
	)
}

func init_net_module() Module {
	module := NewModule()

	// Create a new server instance
	module.exports["create_server"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewServer()
		},
		[]Type{},
		ServerType{},
	)

	// Serve function with improved request/response handling
	module.exports["serve"] = NewNativeFunction(
		func(args ...Value) Value {
			server := args[0].(Server)
			port := args[1].(Number)

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Create request object
				req := NewRequest()
				req.Method = r.Method
				req.Path = r.URL.Path
				req.Headers = parse_headers(r.Header)
				req.Query = parse_query(r.URL.Query())
				body, _ := io.ReadAll(r.Body)
				req.Body = string(body)
				req.ContentType = r.Header.Get("Content-Type")

				// Create response object
				res := NewResponse()
				res.Writer = w

				// Find and execute route handler
				if routeHandlers := server.routes[req.Path]; routeHandlers != nil {
					if handler := routeHandlers[req.Method]; handler != nil {
						handler.Call(req, res)

						// Write response
						for key, value := range res.Headers {
							w.Header().Set(key, value)
						}
						if res.ContentType != "" {
							w.Header().Set("Content-Type", res.ContentType)
						}
						w.WriteHeader(res.StatusCode)
						w.Write([]byte(res.Body.String()))
						return
					}
				}

				// Route not found
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not Found"))
			})

			fmt.Printf("Server listening on port %d\n", int(port.Value()))
			err := http.ListenAndServe(fmt.Sprintf(":%d", int(port.Value())), nil)
			if err != nil {
				panic(err)
			}

			return NewNil()
		},
		[]Type{ServerType{}, PrimitiveType{NumberType}},
		PrimitiveType{NilType},
	)

	return *module
}

func parse_query(values map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, vals := range values {
		if len(vals) > 0 {
			result[key] = vals[0]
		}
	}
	return result
}

func parse_headers(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, vals := range headers {
		if len(vals) > 0 {
			result[key] = vals[0]
		}
	}
	return result
}
