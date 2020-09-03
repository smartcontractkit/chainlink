package store

import (
	"github.com/bianjieai/irita-sdk-go/modules/service"
)

type ServiceRequset struct {
	RequestResponse service.QueryServiceRequestResponse
	Provider        string
}

var serviceMemory = make(map[string]*ServiceRequset)

func GetServiceMemory() map[string]*ServiceRequset {
	return serviceMemory
}

func AddToMemory(key string, value *ServiceRequset) {
	serviceMemory[key] = value
}

func DeleteFromMemory(key string) {
	delete(serviceMemory, key)
}
