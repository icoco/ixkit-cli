package core

import (
	octopuskit "github.com/icoco/ixkit-cli/vendors/octopuskit"
)

type ApiDef struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Route  string `json:"route"`
	Method string `json:"method"`

	Parameters map[string]interface{} `json:"parameters"`
	Headers    map[string]interface{} `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	Response   string                 `json:"response"`

	Extra map[string]interface{} `json:"extra"`
}

func NewApiDef() *ApiDef {
	return &ApiDef{}
}

func (a *ApiDef) Put(k string, v interface{}) bool {
	if nil == a.Extra {
		a.Extra = make(map[string]interface{})
	}
	a.Extra = make(map[string]interface{})
	a.Extra[k] = v
	return true
}

type KitService struct {
	service *octopuskit.Service
}

func SetupManager() {
	m := octopuskit.ManagerInstance()
	endPoint := "http://127.0.0.1/ixkit/"
	m.Setup(endPoint)
}

func (k *KitService) create(data string, output string, success octopuskit.OnSuccess, failure octopuskit.OnFailure) *KitService {
	s := octopuskit.InitService()
	s.Route = "api/app"
	s.Method = "POST"
	s.Parameters["body"] = data
	s.Output = output
	k.service = s

	k.service.Exec(success, failure)

	return k
}

func (k *KitService) maps(data string, success octopuskit.OnSuccess, failure octopuskit.OnFailure) *KitService {
	s := octopuskit.InitService()
	s.Route = "api/service"
	s.Method = "POST"
	s.Parameters["body"] = data
	k.service = s

	k.service.Exec(success, failure)

	return k
}

func (k *KitService) load(data string, output string, success octopuskit.OnSuccess, failure octopuskit.OnFailure) *KitService {
	s := octopuskit.InitService()
	s.Route = "api/app/pull"
	s.Method = "GET"
	s.Parameters["body"] = data
	s.Output = output
	k.service = s

	k.service.Exec(success, failure)

	return k
}
