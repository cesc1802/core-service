package response

type response struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

type errorField struct {
	Field        string `json:"field,omitempty"`
	Tag          string `json:"tag,omitempty"`
	ErrorMessage string `json:"errorMsg,omitempty"`
}

func NewSuccessResponse(data, paging, filter interface{}) *response {
	return &response{Data: data, Paging: paging, Filter: filter}
}

func SimpleSuccessResponse(data interface{}) *response {
	return NewSuccessResponse(data, nil, nil)
}
