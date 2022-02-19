package api

/*
Failure is an interface that returns an api.Status (which implements Error)
if there was an error in an API call.
*/
type Failure interface {
	Failed() *Status
}
