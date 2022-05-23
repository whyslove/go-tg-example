// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	goUUID "github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/seniorGolang/json"
)

const (
	maxParallelBatch = 100
	// Version defines the version of the JSON RPC implementation
	Version = "2.0"
	// contentTypeJson defines the content type to be served
	contentTypeJson = "application/json"
	// ParseError defines invalid JSON was received by the server
	// An error occurred on the server while parsing the JSON text
	ParseError = -32700
	// InvalidRequestError defines the JSON sent is not a valid Request object
	InvalidRequestError = -32600
	// MethodNotFoundError defines the method does not exist / is not available
	MethodNotFoundError = -32601
	// InvalidParamsError defines invalid method parameter(s)
	InvalidParamsError = -32602
	// InternalError defines a server error
	InternalError = -32603
)

type idJsonRPC = json.RawMessage
type ErrorDecoder func(errData json.RawMessage) error

type baseJsonRPC struct {
	ID      idJsonRPC       `json:"id"`
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method,omitempty"`
	Error   json.RawMessage `json:"error,omitempty"`
	Params  interface{}     `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`

	retHandler func(baseJsonRPC)
}

type errorJsonRPC struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err errorJsonRPC) Error() string {
	return err.Message
}

type ClientJsonRPC struct {
	url     string
	name    string
	log     zerolog.Logger
	headers []string

	errorDecoder ErrorDecoder
}

type Batch []baseJsonRPC

func (batch *Batch) Append(request baseJsonRPC) {
	*batch = append(*batch, request)
}
func (batch *Batch) ToArray() []baseJsonRPC {
	return *batch
}

func New(name string, log zerolog.Logger, url string, opts ...Option) (cli *ClientJsonRPC) {
	cli = &ClientJsonRPC{
		errorDecoder: defaultErrorDecoder,
		log:          log,
		name:         name,
		url:          url,
	}

	for _, opt := range opts {
		opt(cli)
	}
	return
}

func (cli *ClientJsonRPC) Mathematical() *ClientMathematical {
	return &ClientMathematical{ClientJsonRPC: cli}
}

func defaultErrorDecoder(errData json.RawMessage) (err error) {

	var jsonrpcError errorJsonRPC
	if err = json.Unmarshal(errData, &jsonrpcError); err != nil {
		return
	}
	return jsonrpcError
}

func (cli *ClientJsonRPC) Batch(ctx context.Context, requests ...baseJsonRPC) (err error) {
	return cli.jsonrpcCall(ctx, cli.log, requests...)
}

func (cli *ClientJsonRPC) BatchFunc(ctx context.Context, batchFunc func(requests *Batch)) (err error) {
	var requests Batch
	batchFunc(&requests)
	return cli.jsonrpcCall(ctx, cli.log, requests...)
}

func (cli *ClientJsonRPC) jsonrpcCall(ctx context.Context, log zerolog.Logger, requests ...baseJsonRPC) (err error) {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	resp := fiber.AcquireResponse()
	agent.SetResponse(resp)
	defer fiber.ReleaseResponse(resp)

	req.SetRequestURI(cli.url)
	agent.ContentType(contentTypeJson)
	req.Header.SetMethod(fiber.MethodPost)
	if err = agent.Parse(); err != nil {
		return
	}

	requestID, _ := ctx.Value(headerRequestID).(string)
	if requestID == "" {
		requestID = goUUID.New().String()
	}
	req.Header.Set(headerRequestID, requestID)
	for _, header := range cli.headers {
		if value, ok := ctx.Value(header).(string); ok {
			req.Header.Set(header, value)
		}
	}
	if err = json.NewEncoder(req.BodyWriter()).Encode(requests); err != nil {
		return
	}
	if err = agent.Do(req, resp); err != nil {
		return
	}
	responseMap := make(map[string]func(baseJsonRPC))
	for _, request := range requests {
		if request.ID != nil {
			responseMap[string(request.ID)] = request.retHandler
		}
	}
	var responses []baseJsonRPC
	if err = json.Unmarshal(resp.Body(), &responses); err != nil {
		cli.log.Error().Err(err).Str("response", string(resp.Body())).Msg("unmarshal response error")
		return
	}
	for _, response := range responses {
		if handler, found := responseMap[string(response.ID)]; found {
			handler(response)
		}
	}
	return
}
