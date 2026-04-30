package environment

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	chipingresspb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// chipDemoLoadTestProto is the raw .proto for schema subject chip-demo-pb.DemoClientPayload.
// Keep in sync with core/services/beholder/chip_load_test_demo.proto and atlas chip-ingress demo client.
const chipDemoLoadTestProto = `syntax = "proto3";

option go_package = "github.com/smartcontractkit/chainlink/v2/core/services/beholder;beholder";

package pb;

message DemoClientPayload {
  string id = 1;
  string domain = 2;
  string entity = 3;
  int64 batch_num = 4;
  int64 message_num = 5;
  int64 batch_position = 6;
}
`

// registerChipDemoLoadTestSchema registers the chip-demo protobuf used by DurableEmitter load tests
// (TestTPS_* with CHIP_INGRESS_TEST_ADDR) against the local CRE Beholder Chip Ingress.
// It uses the same demo basic-auth account as atlas/chip-ingress docker-compose (CE_SA_CHIP_INGRESS_DEMO_CLIENT).
func registerChipDemoLoadTestSchema(ctx context.Context, chipGRPCAddress string) error {
	if strings.TrimSpace(chipGRPCAddress) == "" {
		return errors.New("chip gRPC address is empty")
	}

	opts := []chipingress.Opt{
		chipingress.WithInsecureConnection(),
		chipingress.WithBasicAuth("chip-ingress-demo-client", "password"),
	}
	c, err := chipingress.NewClient(chipGRPCAddress, opts...)
	if err != nil {
		return errors.Wrap(err, "chipingress.NewClient")
	}
	defer func() { _ = c.Close() }()

	regCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err = c.RegisterSchemas(regCtx, &chipingresspb.Schema{
		Subject: "chip-demo-pb.DemoClientPayload",
		Schema:  chipDemoLoadTestProto,
		Format:  chipingresspb.SchemaType_PROTOBUF,
	})
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "already") || strings.Contains(msg, "exists") || strings.Contains(msg, "duplicate") {
			framework.L.Info().Msg("chip-demo load-test schema already registered (chip-demo-pb.DemoClientPayload)")
			return nil
		}
		return errors.Wrap(err, "RegisterSchemas chip-demo-pb.DemoClientPayload")
	}
	framework.L.Info().Msg("registered chip-demo load-test schema (chip-demo-pb.DemoClientPayload) for durable emitter / external Chip tests")
	return nil
}
