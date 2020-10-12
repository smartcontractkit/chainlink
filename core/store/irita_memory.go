package store

import (
	iservice "github.com/irisnet/service-sdk-go/service"
)

type ServiceRequset struct {
	RequestResponse iservice.QueryServiceRequestResponse
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
