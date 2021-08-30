package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/tools"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Flux monitor external validator suite @flux", func() {
	var (
		s                  *actions.DefaultSuiteSetup
		adapter            environment.ExternalAdapter
		nodes              []client.Chainlink
		rac                contracts.ReadAccessController
		flags              contracts.Flags
		dfv                contracts.DeviationFlaggingValidator
		nodeAddresses      []common.Address
		fluxInstance       contracts.FluxAggregator
		fluxRoundConfirmer *contracts.FluxAggregatorRoundConfirmer
		flagSet            bool
		err                error
	)
	fluxRoundTimeout := time.Second * 30

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(3),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())

			s.Client.ParallelTransactions(true)
		})
		By("Deploying access controller, flags, deviation validator", func() {
			rac, err = s.Deployer.DeployReadAccessController(s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			flags, err = s.Deployer.DeployFlags(s.Wallets.Default(), rac.Address())
			Expect(err).ShouldNot(HaveOccurred())
			dfv, err = s.Deployer.DeployDeviationFlaggingValidator(s.Wallets.Default(), flags.Address(), big.NewInt(0))
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying and funding contract", func() {
			fmOpts := contracts.FluxAggregatorOptions{
				PaymentAmount: big.NewInt(1),
				Validator:     common.HexToAddress(dfv.Address()),
				Timeout:       uint32(30),
				MinSubValue:   big.NewInt(0),
				MaxSubValue:   big.NewInt(1e18),
				Decimals:      uint8(0),
				Description:   "Hardhat Flux Aggregator",
			}
			fluxInstance, err = s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), fmOpts)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.Fund(s.Wallets.Default(), nil, big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Setting access to flags contract", func() {
			err = rac.AddAccess(s.Wallets.Default(), dfv.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Funding Chainlink nodes", func() {
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(
				nodes,
				s.Client,
				s.Wallets.Default(),
				big.NewFloat(2),
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Setting oracle options", func() {
			err = fluxInstance.SetOracles(s.Wallets.Default(),
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				})
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})
		By("Creating flux jobs", func() {
			for _, n := range nodes {
				fluxSpec := &client.FluxMonitorJobSpec{
					Name:              "flux_monitor",
					ContractAddress:   fluxInstance.Address(),
					PollTimerPeriod:   15 * time.Second, // min 15s
					PollTimerDisabled: false,
					ObservationSource: client.ObservationSourceSpecHTTP(fmt.Sprintf("%s/variable", adapter.ClusterURL())),
				}
				_, err = n.CreateJob(fluxSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	})

	Describe("with Flux job", func() {
		It("Sets a flag when value is above threshold", func() {
			err = adapter.SetVariable(1e7)
			Expect(err).ShouldNot(HaveOccurred())
			fluxRoundConfirmer = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
			s.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRoundConfirmer)
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			flagSet, err = flags.GetFlag(context.Background(), fluxInstance.Address())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(flagSet).Should(Equal(false))

			err = adapter.SetVariable(1e8)
			Expect(err).ShouldNot(HaveOccurred())
			fluxRoundConfirmer = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(3), fluxRoundTimeout)
			s.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRoundConfirmer)
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			flagSet, err = flags.GetFlag(context.Background(), fluxInstance.Address())
			Expect(err).ShouldNot(HaveOccurred())
			log.Debug().Bool("Flag", flagSet).Msg("Deviation flag set")
			Expect(flagSet).Should(Equal(true))
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", s.TearDown())
	})
})
