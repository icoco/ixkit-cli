package octopuskit

import (
	"fmt"
	"gopkg.in/resty.v1"
	"sync"
)

type Manager struct {
	endPoint string
	agent    string
}

var instantiated *Manager
var once sync.Once

func ManagerInstance() *Manager {
	once.Do(func() {
		instantiated = &Manager{agent: "octopuskit_go v0.1"}
	})
	return instantiated
}

func (m *Manager) Setup(endPoint string) *Manager {
	m.endPoint = endPoint
	return m
}

func (m *Manager) invoke(service Service, success OnSuccess, failure OnFailure) {
	method := service.Method

	var resp *resty.Response
	var err error
	switch method {
	case MethodGet:
		resp, err = m.get(service)
	case MethodPost:
		resp, err = m.post(service)
	}
	if err != nil {
		failure(err)
		return
	}
	success(resp.RawResponse, resp.Body())
}
func (m *Manager) assembleUrl(service Service) string {
	url := m.endPoint + service.Route
	return url
}

func safeMapStrings(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	return mapString
}

func (m *Manager) get(service Service) (*resty.Response, error) {
	url := m.assembleUrl(service)
	p := safeMapStrings(service.Parameters)
	h := safeMapStrings(service.Headers)
	if "" == h["User-Agent"] {
		h["User-Agent"] = m.agent
	}
	if "" != service.Output {
		resp, err := resty.R().SetQueryParams(p).
			SetOutput(service.Output).
			SetHeaders(h).
			Get(url)
		return resp, err
	}
	resp, err := resty.R().SetQueryParams(p).
		SetHeaders(h).
		Get(url)
	return resp, err
}

func (m *Manager) post(service Service) (*resty.Response, error) {
	url := m.assembleUrl(service)
	p := safeMapStrings(service.Parameters)
	h := safeMapStrings(service.Headers)
	if "" == h["User-Agent"] {
		h["User-Agent"] = m.agent
	}

	if "" != service.Output {
		resp, err := resty.R().SetBody(p).
			SetOutput(service.Output).
			SetHeaders(h).
			Post(url)
		return resp, err
	}
	resp, err := resty.R().SetBody(p).
		SetHeaders(h).
		Post(url)
	return resp, err
}
