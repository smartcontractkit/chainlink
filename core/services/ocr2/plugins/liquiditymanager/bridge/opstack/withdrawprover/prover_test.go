package withdrawprover

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l2_output_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/mocks/mock_optimism_l2_output_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/mocks/mock_optimism_portal"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_prover_GetFPAC(t *testing.T) {
	type fields struct {
		optimismPortal *mock_optimism_portal.OptimismPortalInterface
	}
	type args struct {
		ctx context.Context //nolint:containedctx
	}
	tests := []struct {
		name    string
		fields  fields
		expect  func(t *testing.T, fields fields, args args)
		assert  func(t *testing.T, fields fields)
		args    args
		want    bool
		wantErr bool
	}{
		{
			"success, v2.5.0",
			fields{
				optimismPortal: mock_optimism_portal.NewOptimismPortalInterface(t),
			},
			func(t *testing.T, fields fields, args args) {
				fields.optimismPortal.On("Version", mock.Anything).Return("2.5.0", nil)
			},
			func(t *testing.T, fields fields) {
				fields.optimismPortal.AssertExpectations(t)
			},
			args{
				ctx: testutils.Context(t),
			},
			false, // 2.5.0 does not have fault proofs
			false,
		},
		{
			"success, v3.0.0",
			fields{
				optimismPortal: mock_optimism_portal.NewOptimismPortalInterface(t),
			},
			func(t *testing.T, fields fields, args args) {
				fields.optimismPortal.On("Version", mock.Anything).Return("3.0.0", nil)
			},
			func(t *testing.T, fields fields) {
				fields.optimismPortal.AssertExpectations(t)
			},
			args{
				ctx: testutils.Context(t),
			},
			true, // 3.0.0 and beyond do have fault proofs
			false,
		},
		{
			"rpc error",
			fields{
				optimismPortal: mock_optimism_portal.NewOptimismPortalInterface(t),
			},
			func(t *testing.T, fields fields, args args) {
				fields.optimismPortal.On("Version", mock.Anything).Return("", errors.New("rpc error"))
			},
			func(t *testing.T, fields fields) {
				fields.optimismPortal.AssertExpectations(t)
			},
			args{
				ctx: testutils.Context(t),
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &prover{
				optimismPortal: tt.fields.optimismPortal,
			}
			tt.expect(t, tt.fields, tt.args)
			defer tt.assert(t, tt.fields)
			got, err := p.GetFPAC(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_prover_makeStateTrieProof(t *testing.T) {
	type fields struct {
		l2Client *evmclientmocks.Client
	}
	type args struct {
		ctx           context.Context //nolint:containedctx
		l2BlockNumber *big.Int
		address       common.Address
		slot          [32]byte
	}
	tests := []struct {
		name    string
		fields  fields
		expect  func(t *testing.T, fields fields, args args)
		assert  func(t *testing.T, fields fields)
		args    args
		want    stateTrieProof
		wantErr bool
	}{
		// tx: https://sepolia-optimism.etherscan.io/tx/0x14e41dac648d2e1c166ca5c44af09c7c1da684b994ae74dc11303b1ac4bd057c
		{
			"success",
			fields{
				l2Client: evmclientmocks.NewClient(t),
			},
			func(t *testing.T, fields fields, args args) {
				fields.l2Client.On("CallContext",
					mock.Anything,
					mock.Anything,
					"eth_getProof",
					args.address,
					[]string{hexutil.Encode(args.slot[:])},
					hexutil.EncodeBig(args.l2BlockNumber)).
					Return(nil).
					Run(func(args mock.Arguments) {
						respArg := args.Get(1).(*getProofResponse)
						respArg.AccountProof = []hexutil.Bytes{
							hexutil.MustDecode("0xf90211a024e2e559aef2165bfb046a112794fd534e36094b368ff3d70551cf3bdddc82caa08f7b1fa9d398a884bcb1f87b6f43f006b6106c76c4c9126db9a94d4e53ebd8cda0d20f2b640f6a6c070bdbdf11c52912a20d8021250b88d870234a191200607358a05fefa66d31fb7471c2141f7e1212051e684962018b05ba313e9c9c0f337b1271a0540e28bf43afa08069d9569268a9e1e46a56cb04eec4c7df38c2fe649e9b06d2a0bb90aa38bee9aef45c09258836a6b262e9a86f8ba93cb8f91d4cd4203604ffeda0e13893e3c5c659f06845e1a218db9a021b8a5c39b48aeb6813ee046aa8460b12a074abf4864c3a30f9467c24faa4fe0b40065b3363313b6ce2242f61da368a31a1a0e11ee492da64d0c21ac467930ac8ecc049b682aabae217615c53314c3b21dae1a078607fad3708b39bfac5729bdc5b35999be345498247e0791b1f691ac230a4b6a078f235c2a84a078906d21394b6d9e09d623fda4e7b5c473605a246b2f3dfa675a0efd79fd338272aeaee54973cdb538687911844a0df7ba441ca251f3f088011dca03de418148ff7fc6066685ce0abedb3c5849a3dd61b8d6b47dd41775cb10a5dc7a0e2eed53317186ec0e6ff9c430483a9133ddb42c838577740ec073ac159d0e0d5a01194c0f25d9ba2d8294b4e4cb11fe71a4c1566801869f285b24dfaefe746e54ca05a252db0913308f77036adf37c522f270fc804514aa4d39c51225635786ed19780"),
							hexutil.MustDecode("0xf90211a0e8cb47d47c2f31e2d875b77fc256338774990b1d4f296fd4798c0dd580624b4ca0862dd31e18cbfdae5215834cd6720652be168eb7e403e23f7f0d434e45d481fba07fb5d4113b65708ba199e45c74674806b76bec26beb50605f1be4199872064d8a0c65b6c0803b9a22296d787815f9b01a1c449e1ecacde0bf614bd072b5894b313a03201551747e4d6da1d66e679cc08572c6abba5494996088a989eb3e85895d66aa01827646847784668472a42f09e3c55c103fd75fcac2f919897113f61fae292e1a0e117814cb964c3a7c040b21b4b52e4d40f197f7c8900fcd5c386d8c8ddb4c427a04028a79f7165231af5366c53c635dd3fe27af6d913f029b837ae6c2ff59649f9a010df94976d67c8e29cc793f92ef68d6764676a74ededafe16d32d3e3ca7488d5a066d6faac86c1844342c5da3715b1ebdfc9a68a819ceacfca82527a23eddac1bda01adfc361e060b565259abc7f473f382292aef20e4b9902dddb1835568fae2b3ca0d1da4776b1044796ab6ab19f9f2f739ec1595918c4570ca25e264463933d0401a019f31843f55e5a13920d92309ebe20cec5db45ba6337aa57153594b0127c02dfa015b5c0f1c0dbbb3c8a9fa15027ed11ad8aad3bba9ccbf03cad912ea7ffd4d0dda031afbf38e403d55bc9fdf5adc9a04fa02468744059f0469f89f9e82c45c2eebda0ad8cb1320f8f8050bc9796cc9f7c9871b4cd3208425cfeec4ba1adcf011a2d4680"),
							hexutil.MustDecode("0xf90211a0ce236ad4396beefe6ce878ba75d1adbc64664cb6e70e9b0a956a2512ba54acc8a00f0e27384ecb7dc9d114bbc958a6a82f96ee12608a8b412d5143e18a49e7e2aaa085284b8d0b8f081f1cd5a5fe31351e461ddfe4d8400f081420082e0fc6c24cc3a02aefed7ec7df717c20cca036c9a1d8009330131c93e390b34db96aeed560dda6a0d49c6d93f90e3b7cb86b3de94e585a972df928cb17cf9abcc9f6541b33c2d2a4a00d41f6163e879300398bdbaa04b8729b407702f250e7054a554a0fab85b253d0a05951876dcdb04db592a2eadefc87ddcc0f10fb65ccf4e6ce23d05074591e4f52a01d0f22736229c0894aea2dab3476b3d64dfbbb3345df4f518dfeebb8108de47ea02d5ed850b78caaa427db357ef42859cb86601a7f1b890c256f87eb853ee81a60a0f60ba8fbffb6b8229fc7b755e8211ad6a2c6a04761707257117ea8f0ca4f159fa08c1d02fe6f7349d5c648b6ff8d014fddd1b269cc7c021c46326ef5333689cbfca082bf6e558b5ac1d8c48b16557c21ec61627698c777730a05095fd84610b49e61a0e44e37a007a82eda3e675a494f9239dd0f7c83eec09034f8800fccb3849c6341a0bd873a50e5bf2054d96461aadae699a628cac35265a74a7e1f021bcc1f13714ba0ea4cf96963ea2cb93d405ced3feb15566be0520deb1ac40a1b58f2364a6849cca023304a93cbb26f067fe262b4c2baee91f1f3010aac474fba2b65b2d5eb64a20280"),
							hexutil.MustDecode("0xf901f1a0b5584bacb3378841c750cd7259e283e9df6c825246efe92233f49f75145dd7b7a0bc98870fa2bbc4774444ae972fa8941899606a7fcb5ab2d18fcf73218cdf9565a06ac9a816de084bd761fe5f84146b17a59c16b5e8371ed3331e75c39c2b40a842a0c1590fc370a6a7a43486d2e6fc0aa7763d51778090274b130e8773a68b4cda2880a09b126049b8f9516e75f4980f95e8e5482c989dcbc2c545b797cda9c9d651776ba063cea23007f409b1cf99a46147fb195310bd953245c69b55a5170adc3b3e3845a0c31802cba87a226a07d3731ca2bce0422130c8dc62242284b4834fe3c4a3fcd6a0b81c708b51a779c72e3dee040e3b2dad44bca2a24e19dd287a72314cc5d8fdefa02d32182e1d00d6f26c766fa4de677bec312fb1de29db623db07c2b32df5e0980a0f089d911597ed9db6f2804545b438c590dd9a22348012ac61c384b7e191449d3a047ee59a1a3e5f2ca0eec4be3135e19b7340073e2e8786352b2d8f857304b2b4ba022888e04917eb715a60b3596fce61084e802667ba5469ab67a867d6a174f0378a0da7233a9ffd8c03331fc1b9b4f9ef96211fe32cfe89caf2102cf14629ea3cd67a0fcc27a918124452e2fda4913792718d99acb335e0a5aba3c7a0edabf1f7a8850a038e12f16e2ee1b73c7cc93c5de8f60e85e48a9f530f4c761fcd7805baf30ec7e80"),
							hexutil.MustDecode("0xf87180a04c97cb36dfc02b85ff08305900d7b85938ebf278e7438cf6da161e0f3bc8c86ea09e29f4c403bf4c0f9406b2e430379f6489b739b50009f1f8340c5028c7aea02480808080808080a07e3098276a782634c2f70240d7085a7e151668e72fcb14dae91427d1b380b231808080808080"),
							hexutil.MustDecode("0xf8719e30b0147f4cc0e0156d993334777d699c312c2fe454f8b3fa338ed309f4a0b850f84e808a1e92638544b22f3a6b25a099f1d24db21811cbfd50ec024ad875b743aa2fd7b622fbeee274893cb20d67fea01f958654ab06a152993e7a0ae7b6dbb0d4b19265cc9337b8789fe1353bd9dc35"),
						}
						respArg.StorageHash = common.HexToHash("0x99f1d24db21811cbfd50ec024ad875b743aa2fd7b622fbeee274893cb20d67fe")
						val := hexutil.Big{}
						require.NoError(t, val.UnmarshalJSON([]byte(`"0x1"`)))
						respArg.StorageProof = []StorageEntry{
							{
								Key:   hexutil.MustDecode("0x4ee0474e73e1e4e24aa4502f8f2be6821e35dbe367736ca942566749559b2e34"),
								Value: val,
								Proof: []hexutil.Bytes{
									hexutil.MustDecode("0xf90211a035477fbb557c4fb0957fd6cc7377b8dcde1ed58c5697391e4d878380ba3d5f1ba0b26e7c3bfd3749cc5198cb21f81d6200f14c3d7124ca3f9d7b3b9cb8c4211382a05fcec73e0874b9b64d1ca3cc2e8859f58a17a9ffd7f83c028d55f6ab62abe27ea0a6b8d73577801114b703631da0c090e5a7044626c8706e62fa02613127df0afea0a5dd0666ff31ff95320d8e0d687bd2ddc9a514a499e64f17a16a193f4321a0b7a0c9c5ae0c69ea5dd542c56ac0079107879680af4565094fe4b79c63ada85fa3cfa0fc2a7900cd2d0d01cf1e85eb313444309147be2b419be5a4f7356d2668472421a0b92e4635bd9fc24a1d6eda550846ccb074d35b936caf75b7a690a33be6b80417a011c4a7a9ffd57349ef193844e55b8924d62ce9745c3ffc74bcddc11058383307a05825864325ed6bb6e4fdf7d35cc3912550534b2246e6292f9888fdb1cd450e8ba0877dac314f529f22d428281f857ca157f2a6e7db8ff93a10ee1c29750aa440cba03298812964b7595f8fc4bf856035f27056aa16712f85387d6bc6805bd6a1162ba0f809c362b672fb2abea666dc362338778d9e1c6516824c9b731ad077d8151003a0c883efb28212fd4c8e3a5ecb06d344701c86560d882ae867e6ea9759035569f8a03b6cf713fa7d18d6dc69eaafb5436fbfbc7169932b35d7b8015c3e3e37c40111a06262cce3877e42f7cc59e86780eaa441b953e25b4b118992c19269f5191fdae380"),
									hexutil.MustDecode("0xf90211a072f38710cbaa74a3561ee76c45d163fb1886a9256d30ce5b4e0d73304a951331a067547a84a49d5be460c8f48da20254a662ba48e06f26652650b25b8fc7dfc72fa00f792dd6c752ede547714ff4ae4a2b2948b222588ee0340a0a5ba1b8896ce4c6a0c11f1fd9b6b94f4dc2ad7be2de593369ade4c8bd26ec898f70562281ac75abbda0696d3fb880aaa3ddc7ad280d06ae6d6e2f31dad763823b9edf0ed20844e8455da0271c99964a116538f8ed79edbbbf190d08225993c5548fd7cb966c27ba89b010a0bb9a9ff70f473b65847ed4d14d364dd57084e655be6762bd1b39fce4729b167da0ce4c544a62abc0f5fa529ba1b2242e4fe96e0ea59a22d763182d94448e5bb9bda086a659714f804dd6cca55fd9be54b4d97e3be99154a5d048a87e2912fbd86876a0bcbbb6e51bea8065cdfab19f88adf8658e3e82d21ece75d1868a869b6264bd4ca0c47fa199c47cd50ab3174c8f70895adf406db192ae8aa26e60aff11ef826eb4ba0fe9d69a5f8d8e579302402e035a18cf0721b9798988eecdce417a46595457266a00771a06c9d33b27e797fb9a4e34a0ada71c44d7e5dd32195a86a771ad854425ea0a046d28c60ae27cc8928d600ba3f9cfb5735909e1c080d7d301db6cb5b649063a0a13a7c1a4a4bcc3c0fc7f32a254af5fadfcd0ca7bd3c37afaea0a33d5d91f01ba04481ceaae96bdc14128bf414238b18b5a81f0ef94a70a23749c5f4ad16aef4cf80"),
									hexutil.MustDecode("0xf8f180808080808080a002682e3c25cf51d441b6adeb3c1f123a38d46dea8e7b25959b0677317c6d460e80a033fcc0519a9473348c7dd32b4124523d05622728a622a9beabf3e5fa14a8c5d8a018de3e474b028433e4795e015c065587bf216a911f30528ad17f60ab805248d5a0b8c97dd3317e94a95cf181b6828e2211030dee251d6d98cc5d52d96c0a5a8f9fa0a400f40c63a8487c963c3c808ef88825919f60b113e46e7750a283242410736980a00a9a919dacb0c6a821038f01092eb066cdc8d9d3487be818dc2898c614006d53a041baf9eaeec69b50910213e0d7c2eff6eb16b674874fe3d0522fffc0cd17c39b80"),
									hexutil.MustDecode("0xe19f34aa7bea8a53dc373c48131065f2dd7d6c357b53011d84220202183bd6e03601"),
								},
							},
						}
					})
			},
			func(t *testing.T, fields fields) {
			},
			args{
				ctx:           testutils.Context(t),
				l2BlockNumber: big.NewInt(9042600),
				address:       common.HexToAddress("0x4200000000000000000000000000000000000016"),
				slot:          common.HexToHash("0x4ee0474e73e1e4e24aa4502f8f2be6821e35dbe367736ca942566749559b2e34"),
			},
			stateTrieProof{
				AccountProof: [][]byte{
					hexutil.MustDecode("0xf90211a024e2e559aef2165bfb046a112794fd534e36094b368ff3d70551cf3bdddc82caa08f7b1fa9d398a884bcb1f87b6f43f006b6106c76c4c9126db9a94d4e53ebd8cda0d20f2b640f6a6c070bdbdf11c52912a20d8021250b88d870234a191200607358a05fefa66d31fb7471c2141f7e1212051e684962018b05ba313e9c9c0f337b1271a0540e28bf43afa08069d9569268a9e1e46a56cb04eec4c7df38c2fe649e9b06d2a0bb90aa38bee9aef45c09258836a6b262e9a86f8ba93cb8f91d4cd4203604ffeda0e13893e3c5c659f06845e1a218db9a021b8a5c39b48aeb6813ee046aa8460b12a074abf4864c3a30f9467c24faa4fe0b40065b3363313b6ce2242f61da368a31a1a0e11ee492da64d0c21ac467930ac8ecc049b682aabae217615c53314c3b21dae1a078607fad3708b39bfac5729bdc5b35999be345498247e0791b1f691ac230a4b6a078f235c2a84a078906d21394b6d9e09d623fda4e7b5c473605a246b2f3dfa675a0efd79fd338272aeaee54973cdb538687911844a0df7ba441ca251f3f088011dca03de418148ff7fc6066685ce0abedb3c5849a3dd61b8d6b47dd41775cb10a5dc7a0e2eed53317186ec0e6ff9c430483a9133ddb42c838577740ec073ac159d0e0d5a01194c0f25d9ba2d8294b4e4cb11fe71a4c1566801869f285b24dfaefe746e54ca05a252db0913308f77036adf37c522f270fc804514aa4d39c51225635786ed19780"),
					hexutil.MustDecode("0xf90211a0e8cb47d47c2f31e2d875b77fc256338774990b1d4f296fd4798c0dd580624b4ca0862dd31e18cbfdae5215834cd6720652be168eb7e403e23f7f0d434e45d481fba07fb5d4113b65708ba199e45c74674806b76bec26beb50605f1be4199872064d8a0c65b6c0803b9a22296d787815f9b01a1c449e1ecacde0bf614bd072b5894b313a03201551747e4d6da1d66e679cc08572c6abba5494996088a989eb3e85895d66aa01827646847784668472a42f09e3c55c103fd75fcac2f919897113f61fae292e1a0e117814cb964c3a7c040b21b4b52e4d40f197f7c8900fcd5c386d8c8ddb4c427a04028a79f7165231af5366c53c635dd3fe27af6d913f029b837ae6c2ff59649f9a010df94976d67c8e29cc793f92ef68d6764676a74ededafe16d32d3e3ca7488d5a066d6faac86c1844342c5da3715b1ebdfc9a68a819ceacfca82527a23eddac1bda01adfc361e060b565259abc7f473f382292aef20e4b9902dddb1835568fae2b3ca0d1da4776b1044796ab6ab19f9f2f739ec1595918c4570ca25e264463933d0401a019f31843f55e5a13920d92309ebe20cec5db45ba6337aa57153594b0127c02dfa015b5c0f1c0dbbb3c8a9fa15027ed11ad8aad3bba9ccbf03cad912ea7ffd4d0dda031afbf38e403d55bc9fdf5adc9a04fa02468744059f0469f89f9e82c45c2eebda0ad8cb1320f8f8050bc9796cc9f7c9871b4cd3208425cfeec4ba1adcf011a2d4680"),
					hexutil.MustDecode("0xf90211a0ce236ad4396beefe6ce878ba75d1adbc64664cb6e70e9b0a956a2512ba54acc8a00f0e27384ecb7dc9d114bbc958a6a82f96ee12608a8b412d5143e18a49e7e2aaa085284b8d0b8f081f1cd5a5fe31351e461ddfe4d8400f081420082e0fc6c24cc3a02aefed7ec7df717c20cca036c9a1d8009330131c93e390b34db96aeed560dda6a0d49c6d93f90e3b7cb86b3de94e585a972df928cb17cf9abcc9f6541b33c2d2a4a00d41f6163e879300398bdbaa04b8729b407702f250e7054a554a0fab85b253d0a05951876dcdb04db592a2eadefc87ddcc0f10fb65ccf4e6ce23d05074591e4f52a01d0f22736229c0894aea2dab3476b3d64dfbbb3345df4f518dfeebb8108de47ea02d5ed850b78caaa427db357ef42859cb86601a7f1b890c256f87eb853ee81a60a0f60ba8fbffb6b8229fc7b755e8211ad6a2c6a04761707257117ea8f0ca4f159fa08c1d02fe6f7349d5c648b6ff8d014fddd1b269cc7c021c46326ef5333689cbfca082bf6e558b5ac1d8c48b16557c21ec61627698c777730a05095fd84610b49e61a0e44e37a007a82eda3e675a494f9239dd0f7c83eec09034f8800fccb3849c6341a0bd873a50e5bf2054d96461aadae699a628cac35265a74a7e1f021bcc1f13714ba0ea4cf96963ea2cb93d405ced3feb15566be0520deb1ac40a1b58f2364a6849cca023304a93cbb26f067fe262b4c2baee91f1f3010aac474fba2b65b2d5eb64a20280"),
					hexutil.MustDecode("0xf901f1a0b5584bacb3378841c750cd7259e283e9df6c825246efe92233f49f75145dd7b7a0bc98870fa2bbc4774444ae972fa8941899606a7fcb5ab2d18fcf73218cdf9565a06ac9a816de084bd761fe5f84146b17a59c16b5e8371ed3331e75c39c2b40a842a0c1590fc370a6a7a43486d2e6fc0aa7763d51778090274b130e8773a68b4cda2880a09b126049b8f9516e75f4980f95e8e5482c989dcbc2c545b797cda9c9d651776ba063cea23007f409b1cf99a46147fb195310bd953245c69b55a5170adc3b3e3845a0c31802cba87a226a07d3731ca2bce0422130c8dc62242284b4834fe3c4a3fcd6a0b81c708b51a779c72e3dee040e3b2dad44bca2a24e19dd287a72314cc5d8fdefa02d32182e1d00d6f26c766fa4de677bec312fb1de29db623db07c2b32df5e0980a0f089d911597ed9db6f2804545b438c590dd9a22348012ac61c384b7e191449d3a047ee59a1a3e5f2ca0eec4be3135e19b7340073e2e8786352b2d8f857304b2b4ba022888e04917eb715a60b3596fce61084e802667ba5469ab67a867d6a174f0378a0da7233a9ffd8c03331fc1b9b4f9ef96211fe32cfe89caf2102cf14629ea3cd67a0fcc27a918124452e2fda4913792718d99acb335e0a5aba3c7a0edabf1f7a8850a038e12f16e2ee1b73c7cc93c5de8f60e85e48a9f530f4c761fcd7805baf30ec7e80"),
					hexutil.MustDecode("0xf87180a04c97cb36dfc02b85ff08305900d7b85938ebf278e7438cf6da161e0f3bc8c86ea09e29f4c403bf4c0f9406b2e430379f6489b739b50009f1f8340c5028c7aea02480808080808080a07e3098276a782634c2f70240d7085a7e151668e72fcb14dae91427d1b380b231808080808080"),
					hexutil.MustDecode("0xf8719e30b0147f4cc0e0156d993334777d699c312c2fe454f8b3fa338ed309f4a0b850f84e808a1e92638544b22f3a6b25a099f1d24db21811cbfd50ec024ad875b743aa2fd7b622fbeee274893cb20d67fea01f958654ab06a152993e7a0ae7b6dbb0d4b19265cc9337b8789fe1353bd9dc35"),
				},
				StorageProof: [][]byte{
					hexutil.MustDecode("0xf90211a035477fbb557c4fb0957fd6cc7377b8dcde1ed58c5697391e4d878380ba3d5f1ba0b26e7c3bfd3749cc5198cb21f81d6200f14c3d7124ca3f9d7b3b9cb8c4211382a05fcec73e0874b9b64d1ca3cc2e8859f58a17a9ffd7f83c028d55f6ab62abe27ea0a6b8d73577801114b703631da0c090e5a7044626c8706e62fa02613127df0afea0a5dd0666ff31ff95320d8e0d687bd2ddc9a514a499e64f17a16a193f4321a0b7a0c9c5ae0c69ea5dd542c56ac0079107879680af4565094fe4b79c63ada85fa3cfa0fc2a7900cd2d0d01cf1e85eb313444309147be2b419be5a4f7356d2668472421a0b92e4635bd9fc24a1d6eda550846ccb074d35b936caf75b7a690a33be6b80417a011c4a7a9ffd57349ef193844e55b8924d62ce9745c3ffc74bcddc11058383307a05825864325ed6bb6e4fdf7d35cc3912550534b2246e6292f9888fdb1cd450e8ba0877dac314f529f22d428281f857ca157f2a6e7db8ff93a10ee1c29750aa440cba03298812964b7595f8fc4bf856035f27056aa16712f85387d6bc6805bd6a1162ba0f809c362b672fb2abea666dc362338778d9e1c6516824c9b731ad077d8151003a0c883efb28212fd4c8e3a5ecb06d344701c86560d882ae867e6ea9759035569f8a03b6cf713fa7d18d6dc69eaafb5436fbfbc7169932b35d7b8015c3e3e37c40111a06262cce3877e42f7cc59e86780eaa441b953e25b4b118992c19269f5191fdae380"),
					hexutil.MustDecode("0xf90211a072f38710cbaa74a3561ee76c45d163fb1886a9256d30ce5b4e0d73304a951331a067547a84a49d5be460c8f48da20254a662ba48e06f26652650b25b8fc7dfc72fa00f792dd6c752ede547714ff4ae4a2b2948b222588ee0340a0a5ba1b8896ce4c6a0c11f1fd9b6b94f4dc2ad7be2de593369ade4c8bd26ec898f70562281ac75abbda0696d3fb880aaa3ddc7ad280d06ae6d6e2f31dad763823b9edf0ed20844e8455da0271c99964a116538f8ed79edbbbf190d08225993c5548fd7cb966c27ba89b010a0bb9a9ff70f473b65847ed4d14d364dd57084e655be6762bd1b39fce4729b167da0ce4c544a62abc0f5fa529ba1b2242e4fe96e0ea59a22d763182d94448e5bb9bda086a659714f804dd6cca55fd9be54b4d97e3be99154a5d048a87e2912fbd86876a0bcbbb6e51bea8065cdfab19f88adf8658e3e82d21ece75d1868a869b6264bd4ca0c47fa199c47cd50ab3174c8f70895adf406db192ae8aa26e60aff11ef826eb4ba0fe9d69a5f8d8e579302402e035a18cf0721b9798988eecdce417a46595457266a00771a06c9d33b27e797fb9a4e34a0ada71c44d7e5dd32195a86a771ad854425ea0a046d28c60ae27cc8928d600ba3f9cfb5735909e1c080d7d301db6cb5b649063a0a13a7c1a4a4bcc3c0fc7f32a254af5fadfcd0ca7bd3c37afaea0a33d5d91f01ba04481ceaae96bdc14128bf414238b18b5a81f0ef94a70a23749c5f4ad16aef4cf80"),
					hexutil.MustDecode("0xf8f180808080808080a002682e3c25cf51d441b6adeb3c1f123a38d46dea8e7b25959b0677317c6d460e80a033fcc0519a9473348c7dd32b4124523d05622728a622a9beabf3e5fa14a8c5d8a018de3e474b028433e4795e015c065587bf216a911f30528ad17f60ab805248d5a0b8c97dd3317e94a95cf181b6828e2211030dee251d6d98cc5d52d96c0a5a8f9fa0a400f40c63a8487c963c3c808ef88825919f60b113e46e7750a283242410736980a00a9a919dacb0c6a821038f01092eb066cdc8d9d3487be818dc2898c614006d53a041baf9eaeec69b50910213e0d7c2eff6eb16b674874fe3d0522fffc0cd17c39b80"),
					hexutil.MustDecode("0xe19f34aa7bea8a53dc373c48131065f2dd7d6c357b53011d84220202183bd6e03601"),
				},
				StorageValue: big.NewInt(1),
				StorageRoot:  common.HexToHash("0x99f1d24db21811cbfd50ec024ad875b743aa2fd7b622fbeee274893cb20d67fe"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &prover{
				l2Client: tt.fields.l2Client,
			}
			tt.expect(t, tt.fields, tt.args)
			defer tt.assert(t, tt.fields)
			got, err := p.makeStateTrieProof(tt.args.ctx, tt.args.l2BlockNumber, tt.args.address, tt.args.slot)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_prover_getMessageBedrockOutput(t *testing.T) {
	type fields struct {
		optimismPortal *mock_optimism_portal.OptimismPortalInterface
		l2OutputOracle *mock_optimism_l2_output_oracle.OptimismL2OutputOracleInterface
	}
	type args struct {
		ctx           context.Context //nolint:containedctx
		l2BlockNumber *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		expect  func(t *testing.T, fields fields, args args)
		assert  func(t *testing.T, fields fields)
		args    args
		want    bedrockOutput
		wantErr bool
	}{
		// tx: https://sepolia-optimism.etherscan.io/tx/0x14e41dac648d2e1c166ca5c44af09c7c1da684b994ae74dc11303b1ac4bd057c
		{
			"success",
			fields{
				optimismPortal: mock_optimism_portal.NewOptimismPortalInterface(t),
				l2OutputOracle: mock_optimism_l2_output_oracle.NewOptimismL2OutputOracleInterface(t),
			},
			func(t *testing.T, fields fields, args args) {
				fields.optimismPortal.On("Version", mock.Anything).Return("2.5.0", nil)
				fields.l2OutputOracle.On("GetL2OutputIndexAfter", mock.Anything, args.l2BlockNumber).Return(big.NewInt(75354), nil)
				fields.l2OutputOracle.On("GetL2Output", mock.Anything, big.NewInt(75354)).Return(optimism_l2_output_oracle.TypesOutputProposal{
					OutputRoot:    common.HexToHash("0x6d57c409f16da462d7f25abec017baaea894b8a3e2127d3fd56e2bfe2f4eb692"),
					Timestamp:     big.NewInt(1709887860),
					L2BlockNumber: big.NewInt(9042600),
				}, nil)
			},
			func(t *testing.T, fields fields) {
				fields.optimismPortal.AssertExpectations(t)
				fields.l2OutputOracle.AssertExpectations(t)
			},
			args{
				ctx:           testutils.Context(t),
				l2BlockNumber: big.NewInt(9042540),
			},
			bedrockOutput{
				OutputRoot:    common.HexToHash("0x6d57c409f16da462d7f25abec017baaea894b8a3e2127d3fd56e2bfe2f4eb692"),
				L2OutputIndex: big.NewInt(75354),
				L1Timestamp:   big.NewInt(1709887860),
				L2BlockNumber: big.NewInt(9042600),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &prover{
				optimismPortal: tt.fields.optimismPortal,
				l2OutputOracle: tt.fields.l2OutputOracle,
			}
			tt.expect(t, tt.fields, tt.args)
			defer tt.assert(t, tt.fields)
			got, err := p.getMessageBedrockOutput(tt.args.ctx, tt.args.l2BlockNumber)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
