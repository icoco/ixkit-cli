package octopuskit

import ()

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

type Service struct {
	Route      string
	Method     string
	Parameters map[string]interface{}
	Headers    map[string]interface{}
	Output     string
}

func InitService() *Service {
	s := new(Service)
	s.Parameters = make(map[string]interface{})
	s.Headers = make(map[string]interface{})
	return s
}

func (s *Service) Exec(success OnSuccess, failure OnFailure) {

	m := ManagerInstance()
	m.invoke(*s, success, failure)
}
