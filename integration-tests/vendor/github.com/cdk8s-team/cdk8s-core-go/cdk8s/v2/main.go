package cdk8s

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"cdk8s.ApiObject",
		reflect.TypeOf((*ApiObject)(nil)).Elem(),
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
			j := jsiiProxy_ApiObject{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.ApiObjectMetadata",
		reflect.TypeOf((*ApiObjectMetadata)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.ApiObjectMetadataDefinition",
		reflect.TypeOf((*ApiObjectMetadataDefinition)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "add", GoMethod: "Add"},
			_jsii_.MemberMethod{JsiiMethod: "addAnnotation", GoMethod: "AddAnnotation"},
			_jsii_.MemberMethod{JsiiMethod: "addFinalizers", GoMethod: "AddFinalizers"},
			_jsii_.MemberMethod{JsiiMethod: "addLabel", GoMethod: "AddLabel"},
			_jsii_.MemberMethod{JsiiMethod: "addOwnerReference", GoMethod: "AddOwnerReference"},
			_jsii_.MemberMethod{JsiiMethod: "getLabel", GoMethod: "GetLabel"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "namespace", GoGetter: "Namespace"},
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
		},
		func() interface{} {
			return &jsiiProxy_ApiObjectMetadataDefinition{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.ApiObjectProps",
		reflect.TypeOf((*ApiObjectProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.App",
		reflect.TypeOf((*App)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "charts", GoGetter: "Charts"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "outdir", GoGetter: "Outdir"},
			_jsii_.MemberProperty{JsiiProperty: "outputFileExtension", GoGetter: "OutputFileExtension"},
			_jsii_.MemberMethod{JsiiMethod: "synth", GoMethod: "Synth"},
			_jsii_.MemberMethod{JsiiMethod: "synthYaml", GoMethod: "SynthYaml"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberProperty{JsiiProperty: "yamlOutputType", GoGetter: "YamlOutputType"},
		},
		func() interface{} {
			j := jsiiProxy_App{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.AppProps",
		reflect.TypeOf((*AppProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Chart",
		reflect.TypeOf((*Chart)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "generateObjectName", GoMethod: "GenerateObjectName"},
			_jsii_.MemberProperty{JsiiProperty: "labels", GoGetter: "Labels"},
			_jsii_.MemberProperty{JsiiProperty: "namespace", GoGetter: "Namespace"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Chart{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.ChartProps",
		reflect.TypeOf((*ChartProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Cron",
		reflect.TypeOf((*Cron)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "expressionString", GoGetter: "ExpressionString"},
		},
		func() interface{} {
			return &jsiiProxy_Cron{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.CronOptions",
		reflect.TypeOf((*CronOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.DependencyGraph",
		reflect.TypeOf((*DependencyGraph)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "root", GoGetter: "Root"},
			_jsii_.MemberMethod{JsiiMethod: "topology", GoMethod: "Topology"},
		},
		func() interface{} {
			return &jsiiProxy_DependencyGraph{}
		},
	)
	_jsii_.RegisterClass(
		"cdk8s.DependencyVertex",
		reflect.TypeOf((*DependencyVertex)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addChild", GoMethod: "AddChild"},
			_jsii_.MemberProperty{JsiiProperty: "inbound", GoGetter: "Inbound"},
			_jsii_.MemberProperty{JsiiProperty: "outbound", GoGetter: "Outbound"},
			_jsii_.MemberMethod{JsiiMethod: "topology", GoMethod: "Topology"},
			_jsii_.MemberProperty{JsiiProperty: "value", GoGetter: "Value"},
		},
		func() interface{} {
			return &jsiiProxy_DependencyVertex{}
		},
	)
	_jsii_.RegisterClass(
		"cdk8s.Duration",
		reflect.TypeOf((*Duration)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toDays", GoMethod: "ToDays"},
			_jsii_.MemberMethod{JsiiMethod: "toHours", GoMethod: "ToHours"},
			_jsii_.MemberMethod{JsiiMethod: "toHumanString", GoMethod: "ToHumanString"},
			_jsii_.MemberMethod{JsiiMethod: "toIsoString", GoMethod: "ToIsoString"},
			_jsii_.MemberMethod{JsiiMethod: "toMilliseconds", GoMethod: "ToMilliseconds"},
			_jsii_.MemberMethod{JsiiMethod: "toMinutes", GoMethod: "ToMinutes"},
			_jsii_.MemberMethod{JsiiMethod: "toSeconds", GoMethod: "ToSeconds"},
			_jsii_.MemberMethod{JsiiMethod: "unitLabel", GoMethod: "UnitLabel"},
		},
		func() interface{} {
			return &jsiiProxy_Duration{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.GroupVersionKind",
		reflect.TypeOf((*GroupVersionKind)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Helm",
		reflect.TypeOf((*Helm)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "apiObjects", GoGetter: "ApiObjects"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberProperty{JsiiProperty: "releaseName", GoGetter: "ReleaseName"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Helm{}
			_jsii_.InitJsiiProxy(&j.jsiiProxy_Include)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.HelmProps",
		reflect.TypeOf((*HelmProps)(nil)).Elem(),
	)
	_jsii_.RegisterInterface(
		"cdk8s.IAnyProducer",
		reflect.TypeOf((*IAnyProducer)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_IAnyProducer{}
		},
	)
	_jsii_.RegisterClass(
		"cdk8s.Include",
		reflect.TypeOf((*Include)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberProperty{JsiiProperty: "apiObjects", GoGetter: "ApiObjects"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Include{}
			_jsii_.InitJsiiProxy(&j.Type__constructsConstruct)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.IncludeProps",
		reflect.TypeOf((*IncludeProps)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.JsonPatch",
		reflect.TypeOf((*JsonPatch)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_JsonPatch{}
		},
	)
	_jsii_.RegisterClass(
		"cdk8s.Lazy",
		reflect.TypeOf((*Lazy)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "produce", GoMethod: "Produce"},
		},
		func() interface{} {
			return &jsiiProxy_Lazy{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.NameOptions",
		reflect.TypeOf((*NameOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Names",
		reflect.TypeOf((*Names)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Names{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.OwnerReference",
		reflect.TypeOf((*OwnerReference)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Size",
		reflect.TypeOf((*Size)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "toGibibytes", GoMethod: "ToGibibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toKibibytes", GoMethod: "ToKibibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toMebibytes", GoMethod: "ToMebibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toPebibytes", GoMethod: "ToPebibytes"},
			_jsii_.MemberMethod{JsiiMethod: "toTebibytes", GoMethod: "ToTebibytes"},
		},
		func() interface{} {
			return &jsiiProxy_Size{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.SizeConversionOptions",
		reflect.TypeOf((*SizeConversionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"cdk8s.SizeRoundingBehavior",
		reflect.TypeOf((*SizeRoundingBehavior)(nil)).Elem(),
		map[string]interface{}{
			"FAIL": SizeRoundingBehavior_FAIL,
			"FLOOR": SizeRoundingBehavior_FLOOR,
			"NONE": SizeRoundingBehavior_NONE,
		},
	)
	_jsii_.RegisterClass(
		"cdk8s.Testing",
		reflect.TypeOf((*Testing)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Testing{}
		},
	)
	_jsii_.RegisterStruct(
		"cdk8s.TimeConversionOptions",
		reflect.TypeOf((*TimeConversionOptions)(nil)).Elem(),
	)
	_jsii_.RegisterClass(
		"cdk8s.Yaml",
		reflect.TypeOf((*Yaml)(nil)).Elem(),
		nil, // no members
		func() interface{} {
			return &jsiiProxy_Yaml{}
		},
	)
	_jsii_.RegisterEnum(
		"cdk8s.YamlOutputType",
		reflect.TypeOf((*YamlOutputType)(nil)).Elem(),
		map[string]interface{}{
			"FILE_PER_APP": YamlOutputType_FILE_PER_APP,
			"FILE_PER_CHART": YamlOutputType_FILE_PER_CHART,
			"FILE_PER_RESOURCE": YamlOutputType_FILE_PER_RESOURCE,
			"FOLDER_PER_CHART_FILE_PER_RESOURCE": YamlOutputType_FOLDER_PER_CHART_FILE_PER_RESOURCE,
		},
	)
}
