package crecli

import (
	"context"

	"github.com/pkg/errors"
)

type DataSource struct {
	objects []string
	handles map[string]func(context.Context, string) ([]byte, error)
	params  map[string]string
}

func NewDataSource() *DataSource {
	return &DataSource{
		objects: make([]string, 0),
		handles: make(map[string]func(context.Context, string) ([]byte, error)),
		params:  make(map[string]string),
	}
}

func (ds *DataSource) AddObject(
	object string,
	handle func(context.Context, string) ([]byte, error),
	param string,
) {
	ds.objects = append(ds.objects, object)
	ds.handles[object] = handle
	ds.params[object] = param
}

func (ds *DataSource) GetData(ctx context.Context, object string) ([]byte, error) {
	if handle, exists := ds.handles[object]; exists {
		param, ok := ds.params[object]
		if !ok {
			return nil, errors.Errorf("no parameter found for object %s", object)
		}
		return handle(ctx, param)
	}
	return nil, errors.Errorf("no handle found for object %s", object)
}
