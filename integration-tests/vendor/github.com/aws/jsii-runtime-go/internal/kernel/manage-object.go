package kernel

import (
	"reflect"

	"github.com/aws/jsii-runtime-go/internal/api"
)

const objectFQN = "Object"

func (c *Client) ManageObject(v reflect.Value) (ref api.ObjectRef, err error) {
	// Ensuring we use a pointer, so we can see pointer-receiver methods, too.
	var vt reflect.Type
	if v.Kind() == reflect.Interface || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Interface) {
		vt = reflect.Indirect(reflect.ValueOf(v.Interface())).Addr().Type()
	} else {
		vt = reflect.Indirect(v).Addr().Type()
	}
	interfaces, overrides := c.Types().DiscoverImplementation(vt)

	var resp CreateResponse
	resp, err = c.Create(CreateProps{
		FQN:        objectFQN,
		Interfaces: interfaces,
		Overrides:  overrides,
	})

	if err == nil {
		if err = c.objects.Register(v, api.ObjectRef{InstanceID: resp.InstanceID, Interfaces: interfaces}); err == nil {
			ref.InstanceID = resp.InstanceID
		}
	}

	return
}
