package workflows

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
)

const WorkflowID = "WorkflowID"
const WorkflowExecutionID = "WorkflowExecutionID"

type keystoneWorkflowContextKey struct{}

var keystoneContextKey = keystoneWorkflowContextKey{}

type KeystoneWorkflowLabels struct {
	WorkflowExecutionID string `json:"workflowExecutionID"`
	WorkflowID          string `json:"workflowID"`
}

var OrderedKeystoneLabels = []string{WorkflowID, WorkflowExecutionID}

var OrderedKeystoneLabelsMap = make(map[string]interface{})

func init() {
	for _, label := range OrderedKeystoneLabels {
		OrderedKeystoneLabelsMap[label] = interface{}(0)
	}
}

func (k *KeystoneWorkflowLabels) ToMap() map[string]string {
	labels := make(map[string]string)

	reflectedLabels := reflect.ValueOf(k)
	for i := range reflectedLabels.NumField() {
		field := reflectedLabels.Field(i)

		// Get the field name (the exported name)
		fieldName := reflectedLabels.Type().Field(i).Name

		// Get the field value
		fieldValue := field.Interface()

		// Cast and populate labels
		strValue, ok := fieldValue.(string)
		if !ok {
			log.Fatalf("Could not convert %v to a string", fieldValue)
		}
		labels[fieldName] = strValue
	}

	return labels
}

func sendLogAsCustomMessage(ctx context.Context, format string, values ...interface{}) {
	msg := fmt.Sprintf(format, values...)

	// OPTION A - Keys are added individually to the context
	for _, label := range OrderedKeystoneLabels {
		val := ctx.Value(label)
		if val != nil {
			msg = fmt.Sprintf("%v.%v", label, msg)
		}
	}

	// OPTION B - One string key is added to the context that stores all labels in json
	// OPTION B.1
	labels := getKeystoneLabelsFromContextUsingMap(ctx)

	// OPTION B.2
	labels = getKeystoneLabelsFromContextUsingReflection(ctx)

	for _, orderedLabelName := range OrderedKeystoneLabels {
		msg = fmt.Sprintf("%v.%v", labels[orderedLabelName], msg)
	}

	// OPTION C - One unexported struct key is added to the context, with public accessors
	structLabels, err := GetKeystoneLabelsFromContext(ctx)
	if err != nil {
		panic("😨")
	}

	labels = structLabels.ToMap()
	for _, orderedLabelName := range OrderedKeystoneLabels {
		msg = fmt.Sprintf("%v.%v", labels[orderedLabelName], msg)
	}

	// Define a custom protobuf payload to emit
	// TODO: add a generalized custom message while beholder can't emit logs
	payload := &pb.TestCustomMessage{
		StringVal: msg,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf")
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		log.Printf("Error emitting message: %v", err)
	}
}

func GetKeystoneLabelsFromContext(ctx context.Context) (KeystoneWorkflowLabels, error) {
	curLabelsAny := ctx.Value(keystoneContextKey)
	curLabels, ok := curLabelsAny.(KeystoneWorkflowLabels)
	if !ok {
		return KeystoneWorkflowLabels{}, fmt.Errorf("context value with keystoneContextKey is not of type KeystoneWorkflowLabels")
	}

	return curLabels, nil
}

// NewKeystoneContext returns a context with the keystoneContextKey loaded. This enables
// the context to be consumed by GetKeystoneLabelsFromContext and KeystoneContextWithLabel.
// labels should not be nil.
func NewKeystoneContext(ctx context.Context, labels KeystoneWorkflowLabels) context.Context {
	return context.WithValue(ctx, keystoneContextKey, labels)
}

// KeystoneContextWithLabel extracts the Keystone Labels set on the passed in immutable context,
// sets the new desired label if valid, and then returns a new context with the updated labels
func KeystoneContextWithLabel(ctx context.Context, key string, value string) (context.Context, error) {
	curLabels, err := GetKeystoneLabelsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if OrderedKeystoneLabelsMap[key] == nil {
		return nil, fmt.Errorf("key %v is not a valid keystone label", key)
	}

	reflectedLabels := reflect.ValueOf(&curLabels).Elem()
	reflectedLabels.FieldByName(key).SetString(value)

	newLabels := reflectedLabels.Interface().(KeystoneWorkflowLabels)
	return context.WithValue(ctx, keystoneContextKey, newLabels), nil
}

// assumes json formatted string value in context
func getKeystoneLabelsFromContextUsingReflection(ctx context.Context) map[string]string {
	jsonLabels, ok := ctx.Value(keystoneContextKey).(string)
	if !ok {
		log.Fatal("KeystoneContextLabel is a type other than string")
	}

	var structuredKeystoneLabels KeystoneWorkflowLabels
	if err := json.Unmarshal([]byte(jsonLabels), &structuredKeystoneLabels); err != nil {
		log.Fatal(err)
	}

	labels := make(map[string]string)

	reflectedLabels := reflect.ValueOf(structuredKeystoneLabels)
	for i := range reflectedLabels.NumField() {
		field := reflectedLabels.Field(i)

		// Get the field name (the exported name)
		fieldName := reflectedLabels.Type().Field(i).Name

		// Get the field value
		fieldValue := field.Interface()

		// Cast and populate labels
		strValue, ok := fieldValue.(string)
		if !ok {
			log.Fatalf("Could not convert %v to a string", fieldValue)
		}
		labels[fieldName] = strValue
	}

	return labels
}

// // assumes json formatted string value in context
func getKeystoneLabelsFromContextUsingMap(ctx context.Context) map[string]string {
	jsonLabels, ok := ctx.Value(keystoneContextKey).(string)
	if !ok {
		log.Fatal("KeystoneContextLabel is a type other than string")
	}

	var rawKeystoneLabels map[string]interface{}
	if err := json.Unmarshal([]byte(jsonLabels), &rawKeystoneLabels); err != nil {
		log.Fatal(err)
	}

	var labels map[string]string
	for key, value := range rawKeystoneLabels {
		strVal, ok := value.(string)
		if !ok {
			log.Fatalf("Failed to convert keystone label %v to string", key)
		}
		labels[key] = strVal
	}

	return labels
}
