package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.PodChaos",
		reflect.TypeOf((*PodChaos)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addJsonPatch", GoMethod: "AddJsonPatch"},
			_jsii_.MemberProperty{JsiiProperty: "apiGroup", GoGetter: "ApiGroup"},
			_jsii_.MemberProperty{JsiiProperty: "apiVersion", GoGetter: "ApiVersion"},
			_jsii_.MemberProperty{JsiiProperty: "chart", GoGetter: "Chart"},
			_jsii_.MemberProperty{JsiiProperty: "kind", GoGetter: "Kind"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_PodChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodChaosProps",
		reflect.TypeOf((*PodChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodChaosSpec",
		reflect.TypeOf((*PodChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.PodChaosSpecAction",
		reflect.TypeOf((*PodChaosSpecAction)(nil)).Elem(),
		map[string]interface{}{
			"POD_KILL": PodChaosSpecAction_POD_KILL,
			"POD_FAILURE": PodChaosSpecAction_POD_FAILURE,
			"CONTAINER_KILL": PodChaosSpecAction_CONTAINER_KILL,
		},
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.PodChaosSpecMode",
		reflect.TypeOf((*PodChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": PodChaosSpecMode_ONE,
			"ALL": PodChaosSpecMode_ALL,
			"FIXED": PodChaosSpecMode_FIXED,
			"FIXED_PERCENT": PodChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": PodChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodChaosSpecSelector",
		reflect.TypeOf((*PodChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.PodChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*PodChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
}
