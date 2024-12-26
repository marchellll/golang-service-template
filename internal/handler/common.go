package handler

type Meta map[string]any

type Response struct {
	Message string `json:"message"`
	Meta    Meta   `json:"meta"`
	Data    any    `json:"data"`
}

func NewResponse() *Response {
	return &Response{
		Meta: Meta{},
		Data: map[string]any{},
	}
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message

	return r
}

func (r *Response) SetMeta(value Meta) *Response {
	r.Meta = value

	return r
}

func (r *Response) AddMeta(key string, value any) *Response {
	r.Meta[key] = value

	return r
}

func (r *Response) SetData(value any) *Response {
	r.Data = value

	return r
}
