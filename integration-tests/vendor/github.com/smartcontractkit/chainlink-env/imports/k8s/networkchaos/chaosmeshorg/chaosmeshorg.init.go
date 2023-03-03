package chaosmeshorg

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"chaos-meshorg.NetworkChaos",
		reflect.TypeOf((*NetworkChaos)(nil)).Elem(),
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
			j := jsiiProxy_NetworkChaos{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosProps",
		reflect.TypeOf((*NetworkChaosProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpec",
		reflect.TypeOf((*NetworkChaosSpec)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.NetworkChaosSpecAction",
		reflect.TypeOf((*NetworkChaosSpecAction)(nil)).Elem(),
		map[string]interface{}{
			"NETEM": NetworkChaosSpecAction_NETEM,
			"DELAY": NetworkChaosSpecAction_DELAY,
			"LOSS": NetworkChaosSpecAction_LOSS,
			"DUPLICATE": NetworkChaosSpecAction_DUPLICATE,
			"CORRUPT": NetworkChaosSpecAction_CORRUPT,
			"PARTITION": NetworkChaosSpecAction_PARTITION,
			"BANDWIDTH": NetworkChaosSpecAction_BANDWIDTH,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecBandwidth",
		reflect.TypeOf((*NetworkChaosSpecBandwidth)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecCorrupt",
		reflect.TypeOf((*NetworkChaosSpecCorrupt)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecDelay",
		reflect.TypeOf((*NetworkChaosSpecDelay)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecDelayReorder",
		reflect.TypeOf((*NetworkChaosSpecDelayReorder)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.NetworkChaosSpecDirection",
		reflect.TypeOf((*NetworkChaosSpecDirection)(nil)).Elem(),
		map[string]interface{}{
			"TO": NetworkChaosSpecDirection_TO,
			"FROM": NetworkChaosSpecDirection_FROM,
			"BOTH": NetworkChaosSpecDirection_BOTH,
			"VALUE_": NetworkChaosSpecDirection_VALUE_,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecDuplicate",
		reflect.TypeOf((*NetworkChaosSpecDuplicate)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecLoss",
		reflect.TypeOf((*NetworkChaosSpecLoss)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.NetworkChaosSpecMode",
		reflect.TypeOf((*NetworkChaosSpecMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": NetworkChaosSpecMode_ONE,
			"ALL": NetworkChaosSpecMode_ALL,
			"FIXED": NetworkChaosSpecMode_FIXED,
			"FIXED_PERCENT": NetworkChaosSpecMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": NetworkChaosSpecMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecSelector",
		reflect.TypeOf((*NetworkChaosSpecSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecSelectorExpressionSelectors",
		reflect.TypeOf((*NetworkChaosSpecSelectorExpressionSelectors)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecTarget",
		reflect.TypeOf((*NetworkChaosSpecTarget)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"chaos-meshorg.NetworkChaosSpecTargetMode",
		reflect.TypeOf((*NetworkChaosSpecTargetMode)(nil)).Elem(),
		map[string]interface{}{
			"ONE": NetworkChaosSpecTargetMode_ONE,
			"ALL": NetworkChaosSpecTargetMode_ALL,
			"FIXED": NetworkChaosSpecTargetMode_FIXED,
			"FIXED_PERCENT": NetworkChaosSpecTargetMode_FIXED_PERCENT,
			"RANDOM_MAX_PERCENT": NetworkChaosSpecTargetMode_RANDOM_MAX_PERCENT,
		},
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecTargetSelector",
		reflect.TypeOf((*NetworkChaosSpecTargetSelector)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"chaos-meshorg.NetworkChaosSpecTargetSelectorExpressionSelectors",
		reflect.TypeOf((*NetworkChaosSpecTargetSelectorExpressionSelectors)(nil)).Elem(),
	)
}
