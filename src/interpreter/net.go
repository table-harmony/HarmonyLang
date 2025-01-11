package interpreter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/table-harmony/HarmonyLang/src/helpers"
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
	Method  string
	Path    string
	Query   map[string]string
	Headers map[string]string
	Body    string
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
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
	StatusCode int
	Headers    map[string]string
	Body       strings.Builder
	Writer     http.ResponseWriter
	Methods    map[string]Function
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
			res.Headers["Content-Type"] = "application/json"
			jsonBytes, err := json.Marshal(convert_to_native(data))
			if err != nil {
				panic(err)
			}
			res.Body.Write(jsonBytes)
			return res
		},
		[]Type{NewMapType(PrimitiveType{AnyType}, PrimitiveType{AnyType})},
		ResponseType{},
	)

	res.Methods["xml"] = NewNativeFunction(
		func(args ...Value) Value {
			data := args[0]
			res.Headers["Content-Type"] = "application/xml"
			xmlMap, err := helpers.MapToXml("", convert_to_native(data).(map[string]interface{}))
			if err != nil {
				panic(err)
			}
			res.Body.Write(xmlMap)
			return res
		},
		[]Type{NewMapType(PrimitiveType{AnyType}, PrimitiveType{AnyType})},
		ResponseType{},
	)

	res.Methods["html"] = NewNativeFunction(
		func(args ...Value) Value {
			content := args[0].(String).Value()
			res.Headers["Content-Type"] = "text/html"
			res.Body.WriteString(content)
			return res
		},
		[]Type{PrimitiveType{StringType}},
		ResponseType{},
	)

	res.Methods["redirect"] = NewNativeFunction(
		func(args ...Value) Value {
			url := args[0].(String).Value()
			statusCode := int(args[1].(Number).Value())

			res.Headers["Location"] = url
			res.StatusCode = statusCode
			return res
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{NumberType}},
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

	// Serves an http route
	module.exports["serve"] = NewNativeFunction(
		func(args ...Value) Value {
			server := args[0].(Server)
			port := args[1].(Number)

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				req := parse_request(r)

				res := NewResponse()
				res.Writer = w

				if routeHandlers := server.routes[req.Path]; routeHandlers != nil {
					if handler := routeHandlers[req.Method]; handler != nil {
						_, err := handler.Call(req, res)
						if err != nil {
							panic(err)
						}

						for key, value := range res.Headers {
							w.Header().Set(key, value)
						}
						w.WriteHeader(res.StatusCode)
						w.Write([]byte(res.Body.String()))
						return
					}
				}

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

func init_http_module() Module {
	module := NewModule()

	sendRequest := func(method, url string, reqMap Map) Value {
		req := NewRequest()
		req.Path = url
		req.Method = method

		for _, entry := range *reqMap.entries {
			if entry.key.(String).value == "headers" {
				headers := convert_to_native(entry.value).(map[string]interface{})
				for k, v := range headers {
					req.Headers[k] = v.(string)
				}
			}

			if entry.key.(String).value == "query" {
				query := convert_to_native(entry.value).(map[string]interface{})
				for k, v := range query {
					req.Query[k] = v.(string)
				}
			}

			if entry.key.(String).value == "body" {
				req.Body = entry.value.(String).Value()
			}
		}

		var queryParams string
		if strings.Contains(url, "?") {
			queryParams = "&"
		} else {
			queryParams = "?"
		}

		for key, value := range req.Query {
			queryParams += fmt.Sprintf("%s=%s&", key, value)
		}
		fullURL := url + strings.TrimRight(queryParams, "&")

		httpReq, err := http.NewRequest(req.Method, fullURL, strings.NewReader(req.Body))
		if err != nil {
			panic(err)
		}
		for key, value := range req.Headers {
			httpReq.Header.Set(key, value)
		}

		client := &http.Client{}
		res, err := client.Do(httpReq)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		headerEntries := []MapEntry{}
		for key, values := range res.Header {
			if len(values) > 0 {
				headerEntries = append(headerEntries, MapEntry{NewString(key), NewString(values[0])})
			}
		}
		headersMap := NewMap(headerEntries, PrimitiveType{StringType}, PrimitiveType{StringType})

		entries := []MapEntry{
			{NewString("statusCode"), NewNumber(float64(res.StatusCode))},
			{NewString("body"), NewString(string(body))},
			{NewString("headers"), headersMap},
		}

		for key, values := range res.Header {
			if len(values) > 0 {
				entries = append(entries, MapEntry{NewString(key), NewString(values[0])})
			}
		}

		result := NewMap(entries, PrimitiveType{StringType}, PrimitiveType{AnyType})
		return result
	}

	methods := []string{"get", "post", "put", "patch", "delete"}
	for _, method := range methods {
		module.exports[method] = NewNativeFunction(
			func(args ...Value) Value {
				url := args[0].(String).Value()
				reqMap := args[1].(Map)
				return sendRequest(strings.ToUpper(method), url, reqMap)
			},
			[]Type{PrimitiveType{StringType}, PrimitiveType{AnyType}},
			NewMapType(PrimitiveType{StringType}, PrimitiveType{AnyType}),
		)
	}

	return *module
}

func parse_request(r *http.Request) *Request {
	req := NewRequest()

	req.Method = r.Method
	req.Path = r.URL.Path

	headers := make(map[string]string)
	for key, vals := range r.Header {
		if len(vals) > 0 {
			headers[key] = vals[0]
		}
	}
	req.Headers = headers

	query := make(map[string]string)
	for key, vals := range r.URL.Query() {
		if len(vals) > 0 {
			query[key] = vals[0]
		}
	}
	req.Query = query

	body, _ := io.ReadAll(r.Body)
	req.Body = string(body)

	return req
}
