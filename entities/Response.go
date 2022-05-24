package entities

type Response struct {
	Items []S3Object `json:"items"`
}

func CreateNewResponse() Response {
	return Response{}
}

func (response *Response) AddItem(s3Object S3Object) {
	response.Items = append(response.Items, s3Object)
}
