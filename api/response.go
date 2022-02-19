package api

/*
Response is a generic response container.

This is useful if you don't care about the response from the API, as long as
there wasn't any errors.
*/
type Response struct {
	Status Status `json:"ResponseStatus"`
}

func (o Response) Failed() *Status {
	return o.Status.Failed()
}
