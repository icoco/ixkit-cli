package octopuskit

import (
	"net/http"
)

type OnSuccess func(resp *http.Response, body []byte)
type OnFailure func(err error)

type Event struct {
}
