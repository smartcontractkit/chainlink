package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

const WorkflowID = "WorkflowID"
const WorkflowExecutionID = "WorkflowExecutionID"

type keystoneWorkflowContextKey struct{}

var keystoneContextKey = keystoneWorkflowContextKey{}

type KeystoneWorkflowLabels struct {
	WorkflowExecutionID string
	WorkflowID          string
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

	labels[WorkflowID] = k.WorkflowID
	labels[WorkflowExecutionID] = k.WorkflowExecutionID

	return labels
}

// GetKeystoneLabelsFromContext extracts the KeystoneWorkflowLabels struct set on the
// unexported keystoneContextKey. Call NewKeystoneContext first before usage -
// if the key is unset or the value is not of the expected type GetKeystoneLabelsFromContext will error.
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

func sendLogAsCustomMessage(ctx context.Context, format string, values ...interface{}) error {
	msg, err := composeLabeledMsg(ctx, format, values...)
	if err != nil {
		return fmt.Errorf("sendLogAsCustomMessag failed: %w", err)
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

	return nil
}

func composeLabeledMsg(ctx context.Context, format string, values ...interface{}) (string, error) {
	msg := fmt.Sprintf(format, values...)

	structLabels, err := GetKeystoneLabelsFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("composing labeled message failed: %w", err)
	}

	labels := structLabels.ToMap()

	// Populate labeled message in reverse
	numLabels := len(OrderedKeystoneLabels)
	for i := range numLabels {
		msg = fmt.Sprintf("%v.%v", labels[OrderedKeystoneLabels[numLabels-1-i]], msg)
	}

	return msg, nil
}
