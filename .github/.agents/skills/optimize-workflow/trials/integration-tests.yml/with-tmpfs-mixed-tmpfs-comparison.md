# Workflow Trial Comparison
- **Base Run ID (Trial 1)**: [30280453093](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093) (Status: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093), Runtime: `0:13:32`, Cost: `$1.4727`)
- **New Run ID (Trial 2)**: [30287904253](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253) (Status: [failure](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253), Runtime: `0:18:55`, Cost: `$2.2292`)
- **Runtime Delta**: `+0:05:23` (+39.8%)
- **Cost Delta**: `+$0.7565` (+51.4%)

## Jobs Comparison

### Job: Get PR labels and set runner labels
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025201062) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051015339)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053365` -> `GitHub Actions 1012059573`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:08 | 0:00:06 | -0:00:02 (-25.0%) |

### Job: Enforce CTF Version
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025201080) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051015387)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053366` -> `GitHub Actions 1012059574`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:07 | 0:00:08 | +0:00:01 (+14.3%) |

### Job: Check Paths That Require Tests To Run
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025201171) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051015282)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053367` -> `GitHub Actions 1012059575`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:08 | 0:00:06 | -0:00:02 (-25.0%) |

### Job: Run CCIP v1.6 E2E Tests Setup
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025286828) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051102754)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053405` -> `GitHub Actions 1012059581`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:04 | 0:00:03 | -0:00:01 (-25.0%) |

### Job: Build Chainlink Image (plugins)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025286924) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051102800)
- **Runner**: `['runs-on=30280453093-plugins/cpu=16/memory=64/family=m7i+m8i/spot=co/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-plugins/cpu=16/memory=64/family=m7i+m8i/spot=co/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0554c2e1d7437ceb7--tuyqywbd7a` -> `runs-on--i-0fce843f45c3fe267--kr0i1uui0q`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:05:21 | 0:05:56 | +0:00:35 (+10.9%) |
| **Instance Type** | m8i.4xlarge | m7i-flex.4xlarge | m8i.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0319 | $0.0827 | +$0.0508 (+159.2%) |
| **Savings** | $0.2201 (87.4%) | $0.2113 (71.8%) | - |
| **CPU 1m (Avg)** | 59.21 | 65.74 | avg: 59.21 -> 65.74 (+6.53) |
| **Memory (Avg)** | 5.96 % | 6.00 % | avg: 5.96 % -> 6.00 % (+0.04 %) |
| **Disk I/O (Avg)** | 482.62 MB | 485.33 MB | avg: 482.62 MB -> 485.33 MB (+2.71 MB) |
| **Network I/O (Avg)** | 734.08 MB | 713.65 MB | avg: 734.08 MB -> 713.65 MB (-20.43 MB) |

### Job: Run Core CRE E2E Tests Setup
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025287027) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051102712)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053406` -> `GitHub Actions 1012059580`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:05 | 0:00:04 | -0:00:01 (-20.0%) |

### Job: Build Chainlink Image
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025287046) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051102699)
- **Runner**: `['runs-on=30280453093-core/cpu=16/memory=64/family=m7i+m8i/spot=co/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-core/cpu=16/memory=64/family=m7i+m8i/spot=co/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0d6efd34df3a5081a--mqievmwhv7` -> `runs-on--i-0cd036c29f100a83c--xlkdjubz83`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:05:11 | 0:05:33 | +0:00:22 (+7.1%) |
| **Instance Type** | m8i.4xlarge | m7i-flex.4xlarge | m8i.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0310 | $0.0780 | +$0.0470 (+151.6%) |
| **Savings** | $0.2210 (87.7%) | $0.2160 (73.5%) | - |
| **CPU 1m (Avg)** | 60.76 | 67.49 | avg: 60.76 -> 67.49 (+6.73) |
| **Memory (Avg)** | 6.19 % | 6.11 % | avg: 6.19 % -> 6.11 % (-0.08 %) |
| **Disk I/O (Avg)** | 481.80 MB | 485.31 MB | avg: 481.80 MB -> 485.31 MB (+3.51 MB) |
| **Network I/O (Avg)** | 695.34 MB | 692.03 MB | avg: 695.34 MB -> 692.03 MB (-3.31 MB) |

### Job: Compile CRE & CCIP Tests
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90025335548) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90051166084)
- **Runner**: `['runs-on=30280453093-compile/cpu=32/ram=64/family=c7i+c8i/spot=co/volume=100GB/extras=s3-cache']` -> `['runs-on=30287904253-compile/cpu=32/ram=64/family=c7i+c8i/spot=co/volume=100GB/extras=s3-cache']`
- **Runner Name**: `runs-on--i-06f00fe4cd78180f6--zjb19nnt1f` -> `runs-on--i-0f1aab6399b7d72c1--v9w3rvedrx`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:05 | 0:04:03 | +0:00:58 (+31.4%) |
| **Instance Type** | c8i-flex.8xlarge | c7i-flex.8xlarge | c8i-flex.8xlarge -> c7i-flex.8xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0298 | $0.1025 | +$0.0727 (+244.0%) |
| **Savings** | $0.2982 (90.9%) | $0.3075 (75.0%) | - |
| **CPU 1m (Avg)** | 55.29 | 48.58 | avg: 55.29 -> 48.58 (-6.71) |
| **Memory (Avg)** | 6.64 % | 6.03 % | avg: 6.64 % -> 6.03 % (-0.61 %) |
| **Disk I/O (Avg)** | 5114.41 MB | 4774.37 MB | avg: 5114.41 MB -> 4774.37 MB (-340.04 MB) |
| **Network I/O (Avg)** | 1547.53 MB | 1495.90 MB | avg: 1547.53 MB -> 1495.90 MB (-51.63 MB) |

### Job: Run Core CRE E2E Regression Tests / define-test-matrix
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026881887) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052716267)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053872` -> `GitHub Actions 1012060242`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:05 | 0:00:09 | +0:00:04 (+80.0%) |

### Job: Run Core CRE E2E Tests / define-test-matrix
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026881945) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052716300)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012053871` -> `GitHub Actions 1012060243`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:05 | 0:00:09 | +0:00:04 (+80.0%) |

### Job: Run CCIP v1.6 E2E Tests
- **Status**: [skipped](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026882924) -> [skipped](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052717593)
- **Runner**: `[]` -> `[]`
- **Runner Name**: `Unknown` -> `Unknown`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:00 | 0:00:00 | 0:00:00 (0.0%) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_EVM_Read_TxArtifacts  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956696) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807226)
- **Runner**: `['runs-on=30280453093-5-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-5-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-08951f53c8c554285--sqlwvwacu8` -> `runs-on--i-0d7cc6108a7077e8f--a3wsf159tv`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:20 | 0:05:08 | +0:00:48 (+18.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0311 | $0.0360 | +$0.0049 (+15.8%) |
| **Savings** | $0.2209 (87.7%) | $0.0960 (72.7%) | - |
| **CPU 1m (Avg)** | 6.47 | 8.78 | avg: 6.47 -> 8.78 (+2.31) |
| **Memory (Avg)** | 9.71 % | 15.32 % | avg: 9.71 % -> 15.32 % (+5.61 %) |
| **Disk I/O (Avg)** | 542.90 MB | 5289.85 MB | avg: 542.90 MB -> 5289.85 MB (+4746.95 MB) |
| **Network I/O (Avg)** | 1052.97 MB | 1088.28 MB | avg: 1052.97 MB -> 1088.28 MB (+35.31 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_GRPCSource_Lifecycle  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956708) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807248)
- **Runner**: `['runs-on=30280453093-9-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-9-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-06287ba152a25033f--5yi8arfkvl` -> `runs-on--i-0a68e02ef92603ae9--yz9lael7bb`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:47 | 0:02:06 | +0:00:19 (+17.8%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0186 | $0.0342 | +$0.0156 (+83.9%) |
| **Savings** | $0.1074 (85.3%) | $0.0918 (72.9%) | - |
| **CPU 1m (Avg)** | 0.99 | 1.02 | avg: 0.99 -> 1.02 (+0.03) |
| **Memory (Avg)** | 5.51 % | 5.62 % | avg: 5.51 % -> 5.62 % (+0.11 %) |
| **Disk I/O (Avg)** | 467.09 MB | 481.15 MB | avg: 467.09 MB -> 481.15 MB (+14.06 MB) |
| **Network I/O (Avg)** | 661.08 MB | 723.01 MB | avg: 661.08 MB -> 723.01 MB (+61.93 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_EVM_Read_StateQueries  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956709) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807301)
- **Runner**: `['runs-on=30280453093-4-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-4-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0454b9a066555f9be--inzvpa6op7` -> `runs-on--i-00ac6a57f6ed6163d--57j0pg8xtx`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:53 | 0:04:12 | +0:00:19 (+8.2%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0345 | $0.0611 | +$0.0266 (+77.1%) |
| **Savings** | $0.1755 (83.6%) | $0.1489 (70.9%) | - |
| **CPU 1m (Avg)** | 6.54 | 6.58 | avg: 6.54 -> 6.58 (+0.04) |
| **Memory (Avg)** | 8.89 % | 8.98 % | avg: 8.89 % -> 8.98 % (+0.09 %) |
| **Disk I/O (Avg)** | 539.09 MB | 543.27 MB | avg: 539.09 MB -> 543.27 MB (+4.18 MB) |
| **Network I/O (Avg)** | 1029.73 MB | 1059.36 MB | avg: 1029.73 MB -> 1059.36 MB (+29.63 MB) |

### Job: Run Core CRE E2E Tests / TestMustMintVaultJWTForRequest_UsesRawRequestDigest  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956738) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807431)
- **Runner**: `['runs-on=30280453093-16-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-16-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-089a7985fe574a8e3--05az4id9md` -> `runs-on--i-07c2ac78b05f90a1f--1r1y4wp8he`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:58 | 0:02:23 | +0:00:25 (+21.2%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0163 | $0.0183 | +$0.0020 (+12.3%) |
| **Savings** | $0.1097 (87.1%) | $0.0477 (72.3%) | - |
| **CPU 1m (Avg)** | 0.98 | 1.61 | avg: 0.98 -> 1.61 (+0.63) |
| **Memory (Avg)** | 5.46 % | 8.59 % | avg: 5.46 % -> 8.59 % (+3.13 %) |
| **Disk I/O (Avg)** | 478.54 MB | 3306.35 MB | avg: 478.54 MB -> 3306.35 MB (+2827.81 MB) |
| **Network I/O (Avg)** | 653.46 MB | 705.91 MB | avg: 653.46 MB -> 705.91 MB (+52.45 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_EVM_Write_LogTrigger  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956740) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807298)
- **Runner**: `['runs-on=30280453093-2-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-2-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0cd3dfe1f629fe643--czkwcuwkhl` -> `runs-on--i-0b141515a6b416560--1vc7gw8ak4`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:00 | 0:05:16 | +0:01:16 (+31.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0337 | $0.0366 | +$0.0029 (+8.6%) |
| **Savings** | $0.1763 (83.9%) | $0.0954 (72.3%) | - |
| **CPU 1m (Avg)** | 4.50 | 11.04 | avg: 4.50 -> 11.04 (+6.54) |
| **Memory (Avg)** | 9.48 % | 16.52 % | avg: 9.48 % -> 16.52 % (+7.04 %) |
| **Disk I/O (Avg)** | 552.56 MB | 5650.90 MB | avg: 552.56 MB -> 5650.90 MB (+5098.34 MB) |
| **Network I/O (Avg)** | 1087.63 MB | 1101.45 MB | avg: 1087.63 MB -> 1101.45 MB (+13.82 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Aptos_Suite  (workflow-gateway-aptos)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956746) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807326)
- **Runner**: `['runs-on=30280453093-20-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-20-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-02428621544220812--ytcrmczgwa` -> `runs-on--i-0b5a8ee1953b33376--m67l7cpq1d`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:57 | 0:05:50 | +0:00:53 (+17.8%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0347 | $0.0405 | +$0.0058 (+16.7%) |
| **Savings** | $0.2173 (86.2%) | $0.1135 (73.7%) | - |
| **CPU 1m (Avg)** | 3.42 | 3.85 | avg: 3.42 -> 3.85 (+0.43) |
| **Memory (Avg)** | 6.93 % | 9.74 % | avg: 6.93 % -> 9.74 % (+2.81 %) |
| **Disk I/O (Avg)** | 555.16 MB | 6767.71 MB | avg: 555.16 MB -> 6767.71 MB (+6212.55 MB) |
| **Network I/O (Avg)** | 1597.28 MB | 1582.69 MB | avg: 1597.28 MB -> 1582.69 MB (-14.59 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Suite_Bucket_B  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956751) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807312)
- **Runner**: `['runs-on=30280453093-17-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-17-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-05d4ab207f1c3b3f1--aej4bcjr11` -> `runs-on--i-0abc8ae5101389ea3--h91mhredii`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:05:45 | 0:05:32 | -0:00:13 (-3.8%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0402 | $0.0389 | -$0.0013 (-3.2%) |
| **Savings** | $0.2538 (86.3%) | $0.1151 (74.7%) | - |
| **CPU 1m (Avg)** | 1.84 | 2.34 | avg: 1.84 -> 2.34 (+0.5) |
| **Memory (Avg)** | 5.82 % | 10.46 % | avg: 5.82 % -> 10.46 % (+4.64 %) |
| **Disk I/O (Avg)** | 556.54 MB | 5310.60 MB | avg: 556.54 MB -> 5310.60 MB (+4754.06 MB) |
| **Network I/O (Avg)** | 1088.82 MB | 1065.69 MB | avg: 1088.82 MB -> 1065.69 MB (-23.13 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Solana_LogTrigger  (workflow)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956763) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807305)
- **Runner**: `['runs-on=30280453093-23-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-23-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-09bf30c14969562c5--wydlp4weq4` -> `runs-on--i-00f77600225cd3622--ys5ngnb9no`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:18 | 0:04:58 | +0:00:40 (+15.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0308 | $0.0355 | +$0.0047 (+15.3%) |
| **Savings** | $0.1792 (85.3%) | $0.0965 (73.1%) | - |
| **CPU 1m (Avg)** | 2.68 | 2.93 | avg: 2.68 -> 2.93 (+0.25) |
| **Memory (Avg)** | 7.83 % | 14.11 % | avg: 7.83 % -> 14.11 % (+6.28 %) |
| **Disk I/O (Avg)** | 539.63 MB | 5360.76 MB | avg: 539.63 MB -> 5360.76 MB (+4821.13 MB) |
| **Network I/O (Avg)** | 1139.82 MB | 1137.68 MB | avg: 1139.82 MB -> 1137.68 MB (-2.14 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_HTTP_Action_Regression_Suite  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956765) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807337)
- **Runner**: `['runs-on=30280453093-6-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-6-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-057f62e3dcc43a2c4--c4m47qeks9` -> `runs-on--i-0d7288840821c078c--mb6ra6bvdc`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:31 | 0:04:20 | +0:00:49 (+23.2%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0260 | $0.0317 | +$0.0057 (+21.9%) |
| **Savings** | $0.1840 (87.6%) | $0.0783 (71.1%) | - |
| **CPU 1m (Avg)** | 1.61 | 3.17 | avg: 1.61 -> 3.17 (+1.56) |
| **Memory (Avg)** | 6.13 % | 10.13 % | avg: 6.13 % -> 10.13 % (+4 %) |
| **Disk I/O (Avg)** | 529.01 MB | 4631.18 MB | avg: 529.01 MB -> 4631.18 MB (+4102.17 MB) |
| **Network I/O (Avg)** | 977.86 MB | 973.56 MB | avg: 977.86 MB -> 973.56 MB (-4.3 MB) |

### Job: Run Core CRE E2E Tests / TestVaultStaticTopologies_LoadExpectedConfig  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956766) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807367)
- **Runner**: `['runs-on=30280453093-15-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-15-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-042808067506227b1--78ytqd2n7n` -> `runs-on--i-0fda4a3f283c1a7f3--cnnb76aoag`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:53 | 0:02:23 | +0:00:30 (+26.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0152 | $0.0189 | +$0.0037 (+24.3%) |
| **Savings** | $0.1108 (88.0%) | $0.0471 (71.3%) | - |
| **CPU 1m (Avg)** | 0.83 | 1.62 | avg: 0.83 -> 1.62 (+0.79) |
| **Memory (Avg)** | 5.29 % | 9.55 % | avg: 5.29 % -> 9.55 % (+4.26 %) |
| **Disk I/O (Avg)** | 479.21 MB | 3484.30 MB | avg: 479.21 MB -> 3484.30 MB (+3005.09 MB) |
| **Network I/O (Avg)** | 683.31 MB | 726.90 MB | avg: 683.31 MB -> 726.90 MB (+43.59 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Beholder_Suite  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956771) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807345)
- **Runner**: `['runs-on=30280453093-7-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-7-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-032c8dc3ee82673f4--fxyucopgat` -> `runs-on--i-0329c78dc1e931929--vvrvdgocau`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:08 | 0:05:28 | +0:02:20 (+74.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0236 | $0.0384 | +$0.0148 (+62.7%) |
| **Savings** | $0.1444 (86.0%) | $0.1156 (75.0%) | - |
| **CPU 1m (Avg)** | 9.15 | 10.48 | avg: 9.15 -> 10.48 (+1.33) |
| **Memory (Avg)** | 5.90 % | 7.63 % | avg: 5.90 % -> 7.63 % (+1.73 %) |
| **Disk I/O (Avg)** | 528.31 MB | 6525.76 MB | avg: 528.31 MB -> 6525.76 MB (+5997.45 MB) |
| **Network I/O (Avg)** | 894.72 MB | 917.90 MB | avg: 894.72 MB -> 917.90 MB (+23.18 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Suite_Bucket_C  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956774) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807380)
- **Runner**: `['runs-on=30280453093-1-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-1-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0bb4f8c3d349b2fd6--72s7vrb0m7` -> `runs-on--i-0183f48eee253c169--4mdoa2ja3o`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:36 | 0:04:44 | +0:01:08 (+31.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0265 | $0.0340 | +$0.0075 (+28.3%) |
| **Savings** | $0.1835 (87.4%) | $0.0980 (74.2%) | - |
| **CPU 1m (Avg)** | 5.62 | 9.84 | avg: 5.62 -> 9.84 (+4.22) |
| **Memory (Avg)** | 12.61 % | 19.00 % | avg: 12.61 % -> 19.00 % (+6.39 %) |
| **Disk I/O (Avg)** | 548.86 MB | 5088.21 MB | avg: 548.86 MB -> 5088.21 MB (+4539.35 MB) |
| **Network I/O (Avg)** | 1087.11 MB | 1078.98 MB | avg: 1087.11 MB -> 1078.98 MB (-8.13 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Solana_Read_Accounts  (workflow)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956776) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807462)
- **Runner**: `['runs-on=30280453093-24-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-24-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-005407f1a16d20d39--7f9xilxocn` -> `runs-on--i-009868d5902cc7f24--6z1ad1lj8m`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:51 | 0:04:40 | +0:00:49 (+21.2%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0282 | $0.0331 | +$0.0049 (+17.4%) |
| **Savings** | $0.1818 (86.5%) | $0.0989 (74.9%) | - |
| **CPU 1m (Avg)** | 1.51 | 2.57 | avg: 1.51 -> 2.57 (+1.06) |
| **Memory (Avg)** | 7.40 % | 13.93 % | avg: 7.40 % -> 13.93 % (+6.53 %) |
| **Disk I/O (Avg)** | 537.40 MB | 5063.05 MB | avg: 537.40 MB -> 5063.05 MB (+4525.65 MB) |
| **Network I/O (Avg)** | 1092.33 MB | 1110.95 MB | avg: 1092.33 MB -> 1110.95 MB (+18.62 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Module_Cache  (workflow-gateway-cache-test)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956780) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807424)
- **Runner**: `['runs-on=30280453093-26-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-26-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0eda2c2eb5d2ae6ab--eufppqsnz5` -> `runs-on--i-06cd11b6aa0884f38--6gvxk4viqt`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:54 | 0:05:15 | +0:01:21 (+34.6%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0283 | $0.0364 | +$0.0081 (+28.6%) |
| **Savings** | $0.1817 (86.5%) | $0.0956 (72.4%) | - |
| **CPU 1m (Avg)** | 3.17 | 6.99 | avg: 3.17 -> 6.99 (+3.82) |
| **Memory (Avg)** | 8.44 % | 17.28 % | avg: 8.44 % -> 17.28 % (+8.84 %) |
| **Disk I/O (Avg)** | 531.99 MB | 5169.39 MB | avg: 531.99 MB -> 5169.39 MB (+4637.4 MB) |
| **Network I/O (Avg)** | 979.96 MB | 1064.34 MB | avg: 979.96 MB -> 1064.34 MB (+84.38 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Suite_Bucket_A  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956785) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807285)
- **Runner**: `['runs-on=30280453093-0-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-0-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-055e3a06d909155eb--d9zj0sjut1` -> `runs-on--i-0a9387ecfb7704102--r49w9qatrw`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:56 | 0:05:18 | +0:01:22 (+34.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0344 | $0.0370 | +$0.0026 (+7.6%) |
| **Savings** | $0.1756 (83.6%) | $0.0950 (72.0%) | - |
| **CPU 1m (Avg)** | 7.34 | 10.13 | avg: 7.34 -> 10.13 (+2.79) |
| **Memory (Avg)** | 9.80 % | 16.53 % | avg: 9.80 % -> 16.53 % (+6.73 %) |
| **Disk I/O (Avg)** | 546.49 MB | 5486.44 MB | avg: 546.49 MB -> 5486.44 MB (+4939.95 MB) |
| **Network I/O (Avg)** | 1067.70 MB | 1072.58 MB | avg: 1067.70 MB -> 1072.58 MB (+4.88 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Sharding  (workflow-gateway-sharded)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956792) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807856)
- **Runner**: `['runs-on=30280453093-25-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-25-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-01da8c3cbca2aeaf3--pf44osje8j` -> `runs-on--i-0b7b60fb991e3497a--kwdp2eqaav`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:48 | 0:05:06 | +0:01:18 (+34.2%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0336 | $0.0355 | +$0.0019 (+5.7%) |
| **Savings** | $0.1764 (84.0%) | $0.0965 (73.1%) | - |
| **CPU 1m (Avg)** | 5.42 | 12.74 | avg: 5.42 -> 12.74 (+7.32) |
| **Memory (Avg)** | 8.12 % | 15.71 % | avg: 8.12 % -> 15.71 % (+7.59 %) |
| **Disk I/O (Avg)** | 531.23 MB | 5141.36 MB | avg: 531.23 MB -> 5141.36 MB (+4610.13 MB) |
| **Network I/O (Avg)** | 981.27 MB | 1043.70 MB | avg: 981.27 MB -> 1043.70 MB (+62.43 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Stellar_Suite  (workflow-gateway-stellar)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956795) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807346)
- **Runner**: `['runs-on=30280453093-21-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-21-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0f59ec975366ce98e--xi5avcm0q2` -> `runs-on--i-0e3dda61353ba3b4b--fljeklvns6`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:59 | 0:05:32 | +0:00:33 (+11.0%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0425 | $0.0390 | -$0.0035 (-8.2%) |
| **Savings** | $0.2095 (83.1%) | $0.1150 (74.7%) | - |
| **CPU 1m (Avg)** | 2.71 | 3.04 | avg: 2.71 -> 3.04 (+0.33) |
| **Memory (Avg)** | 7.06 % | 11.70 % | avg: 7.06 % -> 11.70 % (+4.64 %) |
| **Disk I/O (Avg)** | 579.85 MB | 5946.79 MB | avg: 579.85 MB -> 5946.79 MB (+5366.94 MB) |
| **Network I/O (Avg)** | 1304.75 MB | 1236.02 MB | avg: 1304.75 MB -> 1236.02 MB (-68.73 MB) |

### Job: Run Core CRE E2E Tests / TestMultiGatewayTopology_LoadExpectedConfig  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956796) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807262)
- **Runner**: `['runs-on=30280453093-11-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-11-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-084b02434646fb51a--l5rp0r0yns` -> `runs-on--i-0c5ba563e92fa43ef--abo1u5rl53`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:50 | 0:02:26 | +0:00:36 (+32.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0186 | $0.0188 | +$0.0002 (+1.1%) |
| **Savings** | $0.1074 (85.3%) | $0.0472 (71.4%) | - |
| **CPU 1m (Avg)** | 0.82 | 1.91 | avg: 0.82 -> 1.91 (+1.09) |
| **Memory (Avg)** | 4.78 % | 9.65 % | avg: 4.78 % -> 9.65 % (+4.87 %) |
| **Disk I/O (Avg)** | 455.87 MB | 3610.25 MB | avg: 455.87 MB -> 3610.25 MB (+3154.38 MB) |
| **Network I/O (Avg)** | 617.54 MB | 765.03 MB | avg: 617.54 MB -> 765.03 MB (+147.49 MB) |

### Job: Run Core CRE E2E Tests / TestVaultJSONOmitUnpopulatedEnabled_CRESettingDefaultsDisabled  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956800) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807324)
- **Runner**: `['runs-on=30280453093-13-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-13-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0c66d42fd4f515f7e--q3drr9w483` -> `runs-on--i-05969e3ff7ffb58a4--rvur2tqc0d`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:52 | 0:02:28 | +0:00:36 (+32.1%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0151 | $0.0190 | +$0.0039 (+25.8%) |
| **Savings** | $0.1109 (88.0%) | $0.0470 (71.2%) | - |
| **CPU 1m (Avg)** | 0.74 | 2.12 | avg: 0.74 -> 2.12 (+1.38) |
| **Memory (Avg)** | 5.38 % | 9.05 % | avg: 5.38 % -> 9.05 % (+3.67 %) |
| **Disk I/O (Avg)** | 471.26 MB | 3437.90 MB | avg: 471.26 MB -> 3437.90 MB (+2966.64 MB) |
| **Network I/O (Avg)** | 676.55 MB | 721.59 MB | avg: 676.55 MB -> 721.59 MB (+45.04 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_EVM_Read_HeavyCalls  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956808) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807320)
- **Runner**: `['runs-on=30280453093-3-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-3-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0be263432baf9e17a--dr0k9jm58q` -> `runs-on--i-02ffde8130d7089ce--6ab35fmr7v`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:47 | 0:05:11 | +0:01:24 (+37.0%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0337 | $0.0360 | +$0.0023 (+6.8%) |
| **Savings** | $0.1763 (84.0%) | $0.0960 (72.7%) | - |
| **CPU 1m (Avg)** | 4.15 | 5.98 | avg: 4.15 -> 5.98 (+1.83) |
| **Memory (Avg)** | 7.80 % | 13.78 % | avg: 7.80 % -> 13.78 % (+5.98 %) |
| **Disk I/O (Avg)** | 521.82 MB | 5231.65 MB | avg: 521.82 MB -> 5231.65 MB (+4709.83 MB) |
| **Network I/O (Avg)** | 957.78 MB | 1075.08 MB | avg: 957.78 MB -> 1075.08 MB (+117.3 MB) |

### Job: Run Core CRE E2E Tests / TestVaultSignedResponseRequestIDEnabled_CRESettingDefaultsDisabled  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956811) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807381)
- **Runner**: `['runs-on=30280453093-14-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-14-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0e1ba20432b392f59--hmgvzch0e2` -> `runs-on--i-062bd44bb807d54e0--1higass9c0`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:48 | 0:02:36 | +0:00:48 (+44.4%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0185 | $0.0197 | +$0.0012 (+6.5%) |
| **Savings** | $0.1075 (85.3%) | $0.0683 (77.7%) | - |
| **CPU 1m (Avg)** | 1.12 | 1.93 | avg: 1.12 -> 1.93 (+0.81) |
| **Memory (Avg)** | 4.16 % | 8.72 % | avg: 4.16 % -> 8.72 % (+4.56 %) |
| **Disk I/O (Avg)** | 454.64 MB | 3526.46 MB | avg: 454.64 MB -> 3526.46 MB (+3071.82 MB) |
| **Network I/O (Avg)** | 611.65 MB | 759.62 MB | avg: 611.65 MB -> 759.62 MB (+147.97 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Suite_Bucket_B  (workflow-gateway-capabilities-vault-jwt_auth-enabled)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956840) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807339)
- **Runner**: `['runs-on=30280453093-18-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-18-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-07f82cc1a24008ba1--kfd4dph5a1` -> `runs-on--i-09dda128d1a647657--jjx5k4gsas`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:05:16 | 0:05:41 | +0:00:25 (+7.9%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.4xlarge | m8ib.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0371 | $0.0802 | +$0.0431 (+116.2%) |
| **Savings** | $0.2149 (85.3%) | $0.2138 (72.7%) | - |
| **CPU 1m (Avg)** | 2.08 | 0.95 | avg: 2.08 -> 0.95 (-1.13) |
| **Memory (Avg)** | 6.35 % | 6.18 % | avg: 6.35 % -> 6.18 % (-0.17 %) |
| **Disk I/O (Avg)** | 553.85 MB | 554.99 MB | avg: 553.85 MB -> 554.99 MB (+1.14 MB) |
| **Network I/O (Avg)** | 1072.00 MB | 1070.87 MB | avg: 1072.00 MB -> 1070.87 MB (-1.13 MB) |

### Job: Run Core CRE E2E Tests / TestVaultOptimizationsEnabled_CRESettingDefaultsDisabled  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956888) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807338)
- **Runner**: `['runs-on=30280453093-12-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-12-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-08b00d402dc7301c2--62q16roa62` -> `runs-on--i-03afe6495ff1559c2--4j2tlcpc2q`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:01:50 | 0:02:27 | +0:00:37 (+33.6%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0189 | $0.0197 | +$0.0008 (+4.2%) |
| **Savings** | $0.1071 (85.0%) | $0.0683 (77.6%) | - |
| **CPU 1m (Avg)** | 0.91 | 1.66 | avg: 0.91 -> 1.66 (+0.75) |
| **Memory (Avg)** | 5.56 % | 8.13 % | avg: 5.56 % -> 8.13 % (+2.57 %) |
| **Disk I/O (Avg)** | 468.20 MB | 3321.94 MB | avg: 468.20 MB -> 3321.94 MB (+2853.74 MB) |
| **Network I/O (Avg)** | 661.26 MB | 693.04 MB | avg: 661.26 MB -> 693.04 MB (+31.78 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Suite_Bucket_B  (workflow-gateway-capabilities-vault-optimizations-enabled)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956892) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807363)
- **Runner**: `['runs-on=30280453093-19-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-19-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-099705f34d08a982c--sdubhl1vkw` -> `runs-on--i-056519d5a6f793c53--hfky3qesyl`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:06:09 | 0:07:24 | +0:01:15 (+20.3%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.4xlarge | m8ib.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0421 | $0.1016 | +$0.0595 (+141.3%) |
| **Savings** | $0.2519 (85.7%) | $0.2344 (69.8%) | - |
| **CPU 1m (Avg)** | 1.96 | 1.49 | avg: 1.96 -> 1.49 (-0.47) |
| **Memory (Avg)** | 6.22 % | 5.94 % | avg: 6.22 % -> 5.94 % (-0.28 %) |
| **Disk I/O (Avg)** | 562.27 MB | 567.59 MB | avg: 562.27 MB -> 567.59 MB (+5.32 MB) |
| **Network I/O (Avg)** | 1126.38 MB | 1143.57 MB | avg: 1126.38 MB -> 1143.57 MB (+17.19 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_GRPCSource_AuthRejection  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956930) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807409)
- **Runner**: `['runs-on=30280453093-10-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-10-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-042671285000cd924--2gz5p0ijq2` -> `runs-on--i-0bcf732fa52b20c5a--2wi20xa5sj`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:19 | 0:03:05 | +0:00:46 (+33.1%) |
| **Instance Type** | m8i.4xlarge | m7i-flex.2xlarge | m8i.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0155 | $0.0233 | +$0.0078 (+50.3%) |
| **Savings** | $0.1105 (87.7%) | $0.0647 (73.5%) | - |
| **CPU 1m (Avg)** | 0.99 | 1.95 | avg: 0.99 -> 1.95 (+0.96) |
| **Memory (Avg)** | 6.12 % | 9.36 % | avg: 6.12 % -> 9.36 % (+3.24 %) |
| **Disk I/O (Avg)** | 501.82 MB | 4040.87 MB | avg: 501.82 MB -> 4040.87 MB (+3539.05 MB) |
| **Network I/O (Avg)** | 815.52 MB | 842.43 MB | avg: 815.52 MB -> 842.43 MB (+26.91 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_DurableEmitter  (workflow-gateway-capabilities)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026956947) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807263)
- **Runner**: `['runs-on=30280453093-8-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-8-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0e3c7aba95920543d--l6a7ec1v10` -> `runs-on--i-0d275404089d13eb1--wrhhpmdn99`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:35 | 0:03:27 | +0:00:52 (+33.5%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0244 | $0.0253 | +$0.0009 (+3.7%) |
| **Savings** | $0.1436 (85.5%) | $0.0627 (71.3%) | - |
| **CPU 1m (Avg)** | 1.62 | 2.55 | avg: 1.62 -> 2.55 (+0.93) |
| **Memory (Avg)** | 6.52 % | 10.01 % | avg: 6.52 % -> 10.01 % (+3.49 %) |
| **Disk I/O (Avg)** | 507.43 MB | 4282.76 MB | avg: 507.43 MB -> 4282.76 MB (+3775.33 MB) |
| **Network I/O (Avg)** | 881.05 MB | 895.71 MB | avg: 881.05 MB -> 895.71 MB (+14.66 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_HTTP_Action_Multi_Gateway  (workflow-gateway-capabilities-multi-gateway)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957130) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807304)
- **Runner**: `['runs-on=30280453093-27-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-27-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-09f25fdb3fa1b09de--n6dbxgoflv` -> `runs-on--i-00373a119ae91af82--xmsfz19obk`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:45 | 0:04:23 | +0:01:38 (+59.4%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0250 | $0.0314 | +$0.0064 (+25.6%) |
| **Savings** | $0.1430 (85.1%) | $0.0786 (71.4%) | - |
| **CPU 1m (Avg)** | 1.30 | 6.80 | avg: 1.30 -> 6.80 (+5.5) |
| **Memory (Avg)** | 7.43 % | 11.15 % | avg: 7.43 % -> 11.15 % (+3.72 %) |
| **Disk I/O (Avg)** | 525.10 MB | 4740.61 MB | avg: 525.10 MB -> 4740.61 MB (+4215.51 MB) |
| **Network I/O (Avg)** | 920.81 MB | 964.95 MB | avg: 920.81 MB -> 964.95 MB (+44.14 MB) |

### Job: Run Core CRE E2E Tests / Test_CRE_V2_Solana_Write  (workflow)
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957170) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052807404)
- **Runner**: `['runs-on=30280453093-22-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-22-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-03bd7c54159401039--wxhibiaryl` -> `runs-on--i-0cd3b24d0393917ca--2hglfl8lmp`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:08 | 0:04:55 | +0:00:47 (+19.0%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0354 | $0.0343 | -$0.0011 (-3.1%) |
| **Savings** | $0.1746 (83.1%) | $0.0977 (74.0%) | - |
| **CPU 1m (Avg)** | 2.31 | 4.58 | avg: 2.31 -> 4.58 (+2.27) |
| **Memory (Avg)** | 8.35 % | 14.58 % | avg: 8.35 % -> 14.58 % (+6.23 %) |
| **Disk I/O (Avg)** | 550.88 MB | 5149.74 MB | avg: 550.88 MB -> 5149.74 MB (+4598.86 MB) |
| **Network I/O (Avg)** | 1173.84 MB | 1130.93 MB | avg: 1173.84 MB -> 1130.93 MB (-42.91 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_HTTP_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957387) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784268)
- **Runner**: `['runs-on=30280453093-2-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-2-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-037408c32f085ff30--fa5tk1aee3` -> `runs-on--i-0dcb6aa90f6fdc887--dmfain0d8c`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:32 | 0:04:04 | +0:01:32 (+60.5%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.2xlarge | m8ib.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0198 | $0.0299 | +$0.0101 (+51.0%) |
| **Savings** | $0.1482 (88.2%) | $0.0801 (72.8%) | - |
| **CPU 1m (Avg)** | 1.05 | 6.08 | avg: 1.05 -> 6.08 (+5.03) |
| **Memory (Avg)** | 5.98 % | 12.24 % | avg: 5.98 % -> 12.24 % (+6.26 %) |
| **Disk I/O (Avg)** | 550.18 MB | 4740.50 MB | avg: 550.18 MB -> 4740.50 MB (+4190.32 MB) |
| **Network I/O (Avg)** | 829.22 MB | 953.00 MB | avg: 829.22 MB -> 953.00 MB (+123.78 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_BalanceAt_Invalid_Address_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957394) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784175)
- **Runner**: `['runs-on=30280453093-3-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-3-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0cc6232483e1ed05e--iw8oehaihm` -> `runs-on--i-0555ad9f66071adc9--0hyxnxttmf`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:58 | 0:04:52 | +0:00:54 (+22.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0346 | $0.0707 | +$0.0361 (+104.3%) |
| **Savings** | $0.1754 (83.5%) | $0.1813 (71.9%) | - |
| **CPU 1m (Avg)** | 12.34 | 18.75 | avg: 12.34 -> 18.75 (+6.41) |
| **Memory (Avg)** | 12.48 % | 12.35 % | avg: 12.48 % -> 12.35 % (-0.13 %) |
| **Disk I/O (Avg)** | 597.46 MB | 596.00 MB | avg: 597.46 MB -> 596.00 MB (-1.46 MB) |
| **Network I/O (Avg)** | 1083.99 MB | 1054.02 MB | avg: 1083.99 MB -> 1054.02 MB (-29.97 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_GetTransactionReceipt_Invalid_Hash_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957464) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784307)
- **Runner**: `['runs-on=30280453093-11-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-11-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-079fc87dc2530c564--q4a91lbnv3` -> `runs-on--i-09b5bc117878a2750--wwwna60fvb`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:29 | 0:04:37 | +0:00:08 (+3.0%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0387 | $0.0664 | +$0.0277 (+71.6%) |
| **Savings** | $0.2133 (84.6%) | $0.1856 (73.7%) | - |
| **CPU 1m (Avg)** | 17.12 | 20.80 | avg: 17.12 -> 20.80 (+3.68) |
| **Memory (Avg)** | 13.72 % | 10.44 % | avg: 13.72 % -> 10.44 % (-3.28 %) |
| **Disk I/O (Avg)** | 593.83 MB | 594.95 MB | avg: 593.83 MB -> 594.95 MB (+1.12 MB) |
| **Network I/O (Avg)** | 1082.84 MB | 1039.64 MB | avg: 1082.84 MB -> 1039.64 MB (-43.2 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_HTTP_Action_CRUD_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957476) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784218)
- **Runner**: `['runs-on=30280453093-18-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-18-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0d181a8d564491f99--m2chdvpmbj` -> `runs-on--i-075ec105a600326d1--ls7qtfswvb`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:57 | 0:03:18 | +0:00:21 (+11.9%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0275 | $0.0471 | +$0.0196 (+71.3%) |
| **Savings** | $0.1405 (83.6%) | $0.1209 (72.0%) | - |
| **CPU 1m (Avg)** | 3.71 | 5.60 | avg: 3.71 -> 5.60 (+1.89) |
| **Memory (Avg)** | 8.34 % | 8.12 % | avg: 8.34 % -> 8.12 % (-0.22 %) |
| **Disk I/O (Avg)** | 562.15 MB | 579.27 MB | avg: 562.15 MB -> 579.27 MB (+17.12 MB) |
| **Network I/O (Avg)** | 885.40 MB | 927.57 MB | avg: 885.40 MB -> 927.57 MB (+42.17 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_CallContract_Invalid_Balance_Reader_Contract_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957479) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784392)
- **Runner**: `['runs-on=30280453093-5-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-5-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-08c4f1a424a426f49--6uv7l00k0o` -> `runs-on--i-0638a6d6f1d523630--9vwc9b1tmc`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:27 | 0:04:45 | +0:00:18 (+6.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0387 | $0.0688 | +$0.0301 (+77.8%) |
| **Savings** | $0.2133 (84.6%) | $0.1832 (72.7%) | - |
| **CPU 1m (Avg)** | 10.43 | 16.36 | avg: 10.43 -> 16.36 (+5.93) |
| **Memory (Avg)** | 10.69 % | 12.39 % | avg: 10.69 % -> 12.39 % (+1.7 %) |
| **Disk I/O (Avg)** | 595.31 MB | 595.60 MB | avg: 595.31 MB -> 595.60 MB (+0.29 MB) |
| **Network I/O (Avg)** | 1074.82 MB | 1050.06 MB | avg: 1074.82 MB -> 1050.06 MB (-24.76 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_WriteReport_Invalid_Receiver_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957503) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784230)
- **Runner**: `['runs-on=30280453093-13-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-13-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-087beb29a402df5fa--6391oiovvc` -> `runs-on--i-0318c53b6d0b0d5c9--j2d95r2yqi`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:20 | 0:04:43 | +0:01:23 (+41.5%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0300 | $0.0674 | +$0.0374 (+124.7%) |
| **Savings** | $0.1380 (82.2%) | $0.1846 (73.2%) | - |
| **CPU 1m (Avg)** | 4.82 | 5.81 | avg: 4.82 -> 5.81 (+0.99) |
| **Memory (Avg)** | 8.18 % | 8.55 % | avg: 8.18 % -> 8.55 % (+0.37 %) |
| **Disk I/O (Avg)** | 575.47 MB | 591.27 MB | avg: 575.47 MB -> 591.27 MB (+15.8 MB) |
| **Network I/O (Avg)** | 955.19 MB | 999.16 MB | avg: 955.19 MB -> 999.16 MB (+43.97 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_Cron_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957510) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784090)
- **Runner**: `['runs-on=30280453093-1-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-1-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0a0ca15161acdaa00--i1bdmqr8dh` -> `runs-on--i-05a93b030d5cea58c--f9mhg5a53q`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:58 | 0:04:14 | +0:01:16 (+42.7%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0267 | $0.0294 | +$0.0027 (+10.1%) |
| **Savings** | $0.1413 (84.1%) | $0.0806 (73.3%) | - |
| **CPU 1m (Avg)** | 0.85 | 3.77 | avg: 0.85 -> 3.77 (+2.92) |
| **Memory (Avg)** | 9.13 % | 13.01 % | avg: 9.13 % -> 13.01 % (+3.88 %) |
| **Disk I/O (Avg)** | 568.65 MB | 4832.09 MB | avg: 568.65 MB -> 4832.09 MB (+4263.44 MB) |
| **Network I/O (Avg)** | 948.39 MB | 1049.74 MB | avg: 948.39 MB -> 1049.74 MB (+101.35 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_GetTransactionByHash_Invalid_Hash_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957515) -> [failure](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784241)
- **Runner**: `['runs-on=30280453093-10-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-10-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-05480c03a3013f4c6--6gel083ef3` -> `runs-on--i-046e33cccc5100635--qza22ghfwx`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:30 | 0:07:15 | +0:02:45 (+61.1%) |
| **Instance Type** | m8i-flex.4xlarge | m7i-flex.2xlarge | m8i-flex.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0281 | $0.0491 | +$0.0210 (+74.7%) |
| **Savings** | $0.2239 (88.9%) | $0.1269 (72.1%) | - |
| **CPU 1m (Avg)** | 19.02 | 23.02 | avg: 19.02 -> 23.02 (+4) |
| **Memory (Avg)** | 13.70 % | 23.59 % | avg: 13.70 % -> 23.59 % (+9.89 %) |
| **Disk I/O (Avg)** | 599.44 MB | 6147.00 MB | avg: 599.44 MB -> 6147.00 MB (+5547.56 MB) |
| **Network I/O (Avg)** | 1073.99 MB | 1159.30 MB | avg: 1073.99 MB -> 1159.30 MB (+85.31 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_HeaderByNumber_Invalid_Block_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957544) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784220)
- **Runner**: `['runs-on=30280453093-12-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-12-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0bd94cc3ac26a0ddf--40cyjsta8d` -> `runs-on--i-0e48db99d20f36078--ygqbmo1bsf`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:20 | 0:05:03 | +0:01:43 (+51.5%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0299 | $0.0358 | +$0.0059 (+19.7%) |
| **Savings** | $0.1381 (82.2%) | $0.0962 (72.9%) | - |
| **CPU 1m (Avg)** | 5.21 | 7.55 | avg: 5.21 -> 7.55 (+2.34) |
| **Memory (Avg)** | 7.94 % | 15.01 % | avg: 7.94 % -> 15.01 % (+7.07 %) |
| **Disk I/O (Avg)** | 568.95 MB | 5272.91 MB | avg: 568.95 MB -> 5272.91 MB (+4703.96 MB) |
| **Network I/O (Avg)** | 945.49 MB | 1049.82 MB | avg: 945.49 MB -> 1049.82 MB (+104.33 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_Consensus_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957548) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784163)
- **Runner**: `['runs-on=30280453093-0-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-0-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0b973914166a8ee24--54cw5vrw66` -> `runs-on--i-0af039854e6671028--hzlbpf1iod`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:18 | 0:03:40 | +0:00:22 (+11.1%) |
| **Instance Type** | m8ib.4xlarge | m7i-flex.4xlarge | m8ib.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0249 | $0.0542 | +$0.0293 (+117.7%) |
| **Savings** | $0.1851 (88.2%) | $0.1558 (74.2%) | - |
| **CPU 1m (Avg)** | 4.40 | 6.15 | avg: 4.40 -> 6.15 (+1.75) |
| **Memory (Avg)** | 8.36 % | 8.11 % | avg: 8.36 % -> 8.11 % (-0.25 %) |
| **Disk I/O (Avg)** | 570.57 MB | 579.22 MB | avg: 570.57 MB -> 579.22 MB (+8.65 MB) |
| **Network I/O (Avg)** | 931.25 MB | 950.45 MB | avg: 931.25 MB -> 950.45 MB (+19.2 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_WriteReport_Invalid_Gas_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957553) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784210)
- **Runner**: `['runs-on=30280453093-16-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-16-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0813dfacc12f6ddba--nlqpkgra4e` -> `runs-on--i-09a8b09b0299f4b59--74140idj85`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:18 | 0:04:05 | +0:00:47 (+23.7%) |
| **Instance Type** | m8i.4xlarge | m7i-flex.4xlarge | m8i.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0210 | $0.0600 | +$0.0390 (+185.7%) |
| **Savings** | $0.1470 (87.5%) | $0.1500 (71.4%) | - |
| **CPU 1m (Avg)** | 4.04 | 2.74 | avg: 4.04 -> 2.74 (-1.3) |
| **Memory (Avg)** | 6.67 % | 6.98 % | avg: 6.67 % -> 6.98 % (+0.31 %) |
| **Disk I/O (Avg)** | 577.16 MB | 587.52 MB | avg: 577.16 MB -> 587.52 MB (+10.36 MB) |
| **Network I/O (Avg)** | 962.13 MB | 966.99 MB | avg: 962.13 MB -> 966.99 MB (+4.86 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_FilterLogs_Invalid_ToBlock_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957558) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784154)
- **Runner**: `['runs-on=30280453093-9-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-9-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-062c1d26e44a48487--d4e24uli7t` -> `runs-on--i-00f64b298ee4dcda0--wdvbmi654w`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:49 | 0:05:15 | +0:01:26 (+37.6%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0335 | $0.0376 | +$0.0041 (+12.2%) |
| **Savings** | $0.1765 (84.0%) | $0.0944 (71.5%) | - |
| **CPU 1m (Avg)** | 10.61 | 18.21 | avg: 10.61 -> 18.21 (+7.6) |
| **Memory (Avg)** | 9.74 % | 15.02 % | avg: 9.74 % -> 15.02 % (+5.28 %) |
| **Disk I/O (Avg)** | 590.86 MB | 5218.03 MB | avg: 590.86 MB -> 5218.03 MB (+4627.17 MB) |
| **Network I/O (Avg)** | 1039.64 MB | 1040.46 MB | avg: 1039.64 MB -> 1040.46 MB (+0.82 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_CallContract_Invalid_Addr_To_Read_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957577) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784285)
- **Runner**: `['runs-on=30280453093-4-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-4-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0f2d201466ca6bf07--pp2g7zwuzp` -> `runs-on--i-05fb54e9d1d506471--pc85s7l9k9`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:23 | 0:05:42 | +0:01:19 (+30.0%) |
| **Instance Type** | m8i.4xlarge | m7i-flex.2xlarge | m8i.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0269 | $0.0404 | +$0.0135 (+50.2%) |
| **Savings** | $0.1831 (87.2%) | $0.1136 (73.8%) | - |
| **CPU 1m (Avg)** | 10.42 | 13.40 | avg: 10.42 -> 13.40 (+2.98) |
| **Memory (Avg)** | 10.95 % | 17.59 % | avg: 10.95 % -> 17.59 % (+6.64 %) |
| **Disk I/O (Avg)** | 598.04 MB | 5417.37 MB | avg: 598.04 MB -> 5417.37 MB (+4819.33 MB) |
| **Network I/O (Avg)** | 1075.24 MB | 1072.37 MB | avg: 1075.24 MB -> 1072.37 MB (-2.87 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_EstimateGas_Invalid_To_Address_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957580) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784372)
- **Runner**: `['runs-on=30280453093-6-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-6-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0150698b91dbc36c8--216dvcqzwb` -> `runs-on--i-0f8f4c6588f006dd9--39t7qcteno`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:21 | 0:05:41 | +0:01:20 (+30.7%) |
| **Instance Type** | m8i-flex.4xlarge | m7i-flex.2xlarge | m8i-flex.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0272 | $0.0396 | +$0.0124 (+45.6%) |
| **Savings** | $0.1828 (87.0%) | $0.1144 (74.3%) | - |
| **CPU 1m (Avg)** | 11.41 | 13.86 | avg: 11.41 -> 13.86 (+2.45) |
| **Memory (Avg)** | 10.94 % | 18.14 % | avg: 10.94 % -> 18.14 % (+7.2 %) |
| **Disk I/O (Avg)** | 597.73 MB | 5554.05 MB | avg: 597.73 MB -> 5554.05 MB (+4956.32 MB) |
| **Network I/O (Avg)** | 1066.16 MB | 1102.30 MB | avg: 1066.16 MB -> 1102.30 MB (+36.14 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_WriteReport_Corrupt_Receiver_Address_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957584) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784525)
- **Runner**: `['runs-on=30280453093-15-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-15-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0e679889ed9acab3c--ucv7jwgras` -> `runs-on--i-0b95d924f320e1290--50tyae64lc`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:17 | 0:04:38 | +0:01:21 (+41.1%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0298 | $0.0329 | +$0.0031 (+10.4%) |
| **Savings** | $0.1382 (82.2%) | $0.0991 (75.0%) | - |
| **CPU 1m (Avg)** | 2.17 | 7.42 | avg: 2.17 -> 7.42 (+5.25) |
| **Memory (Avg)** | 6.98 % | 13.75 % | avg: 6.98 % -> 13.75 % (+6.77 %) |
| **Disk I/O (Avg)** | 568.47 MB | 4891.74 MB | avg: 568.47 MB -> 4891.74 MB (+4323.27 MB) |
| **Network I/O (Avg)** | 928.66 MB | 998.72 MB | avg: 928.66 MB -> 998.72 MB (+70.06 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_Stellar_ReadContract_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957590) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784170)
- **Runner**: `['runs-on=30280453093-19-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-19-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0a8e47dc2487a0f1e--41cvxj387s` -> `runs-on--i-01488a8df6ae11775--ljzw7tmxej`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:43 | 0:05:46 | +0:01:03 (+22.3%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0402 | $0.0408 | +$0.0006 (+1.5%) |
| **Savings** | $0.2118 (84.1%) | $0.1132 (73.5%) | - |
| **CPU 1m (Avg)** | 1.84 | 2.31 | avg: 1.84 -> 2.31 (+0.47) |
| **Memory (Avg)** | 6.40 % | 11.11 % | avg: 6.40 % -> 11.11 % (+4.71 %) |
| **Disk I/O (Avg)** | 623.20 MB | 5643.98 MB | avg: 623.20 MB -> 5643.98 MB (+5020.78 MB) |
| **Network I/O (Avg)** | 1269.57 MB | 1182.75 MB | avg: 1269.57 MB -> 1182.75 MB (-86.82 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_FilterLogs_Invalid_Addresses_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957604) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784196)
- **Runner**: `['runs-on=30280453093-7-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-7-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-08b617f00f2d74905--ynkv0r6zc1` -> `runs-on--i-0d24043eef86aa523--oqk415xs7e`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:04:22 | 0:06:13 | +0:01:51 (+42.4%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0380 | $0.0438 | +$0.0058 (+15.3%) |
| **Savings** | $0.1720 (81.9%) | $0.1102 (71.6%) | - |
| **CPU 1m (Avg)** | 12.91 | 19.54 | avg: 12.91 -> 19.54 (+6.63) |
| **Memory (Avg)** | 13.00 % | 19.70 % | avg: 13.00 % -> 19.70 % (+6.7 %) |
| **Disk I/O (Avg)** | 596.46 MB | 5641.98 MB | avg: 596.46 MB -> 5641.98 MB (+5045.52 MB) |
| **Network I/O (Avg)** | 1078.94 MB | 1097.14 MB | avg: 1078.94 MB -> 1097.14 MB (+18.2 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_FilterLogs_Invalid_FromBlock_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957608) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784225)
- **Runner**: `['runs-on=30280453093-8-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-8-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0f64840b9c7ddd516--9uuks5p2i5` -> `runs-on--i-0fb9d890f81c64703--dj0m5hrt4o`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:03:56 | 0:04:46 | +0:00:50 (+21.2%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0348 | $0.0685 | +$0.0337 (+96.8%) |
| **Savings** | $0.1752 (83.4%) | $0.1835 (72.8%) | - |
| **CPU 1m (Avg)** | 11.77 | 23.08 | avg: 11.77 -> 23.08 (+11.31) |
| **Memory (Avg)** | 11.27 % | 11.94 % | avg: 11.27 % -> 11.94 % (+0.67 %) |
| **Disk I/O (Avg)** | 582.31 MB | 598.15 MB | avg: 582.31 MB -> 598.15 MB (+15.84 MB) |
| **Network I/O (Avg)** | 1001.40 MB | 1055.74 MB | avg: 1001.40 MB -> 1055.74 MB (+54.34 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_LogTrigger_Invalid_Address_Regression
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957661) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784244)
- **Runner**: `['runs-on=30280453093-17-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-17-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']`
- **Runner Name**: `runs-on--i-0851ec8a91b9e8d6a--gictvxfhwc` -> `runs-on--i-0e866802e9a0ee025--5oxewgafjh`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:33 | 0:03:17 | +0:00:44 (+28.8%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.4xlarge | m8in.4xlarge -> m7i-flex.4xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0240 | $0.0501 | +$0.0261 (+108.7%) |
| **Savings** | $0.1440 (85.7%) | $0.1179 (70.2%) | - |
| **CPU 1m (Avg)** | 4.39 | 6.20 | avg: 4.39 -> 6.20 (+1.81) |
| **Memory (Avg)** | 7.10 % | 6.21 % | avg: 7.10 % -> 6.21 % (-0.89 %) |
| **Disk I/O (Avg)** | 551.85 MB | 565.81 MB | avg: 551.85 MB -> 565.81 MB (+13.96 MB) |
| **Network I/O (Avg)** | 836.95 MB | 879.44 MB | avg: 836.95 MB -> 879.44 MB (+42.49 MB) |

### Job: Run Core CRE E2E Regression Tests / Test_CRE_V2_EVM_WriteReport_Failing_On_Receiver
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90026957684) -> [success](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90052784178)
- **Runner**: `['runs-on=30280453093-14-1/cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs']` -> `['runs-on=30287904253-14-1/cpu=8/ram=32/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache']`
- **Runner Name**: `runs-on--i-0c0819c7f4318c8c1--4gx447qovd` -> `runs-on--i-044c4f60f2ff9ef1f--baq0qywr67`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:02:29 | 0:02:20 | -0:00:09 (-6.0%) |
| **Instance Type** | m8in.4xlarge | m7i-flex.2xlarge | m8in.4xlarge -> m7i-flex.2xlarge |
| **Lifecycle** | spot | on-demand | - |
| **Cost** | $0.0232 | $0.0191 | -$0.0041 (-17.7%) |
| **Savings** | $0.1448 (86.2%) | $0.0469 (71.1%) | - |
| **CPU 1m (Avg)** | 0.62 | 1.72 | avg: 0.62 -> 1.72 (+1.1) |
| **Memory (Avg)** | 4.10 % | 9.42 % | avg: 4.10 % -> 9.42 % (+5.32 %) |
| **Disk I/O (Avg)** | 544.52 MB | 3253.76 MB | avg: 544.52 MB -> 3253.76 MB (+2709.24 MB) |
| **Network I/O (Avg)** | 786.26 MB | 680.45 MB | avg: 786.26 MB -> 680.45 MB (-105.81 MB) |

### Job: ETH Smoke Tests
- **Status**: [success](https://github.com/smartcontractkit/chainlink/actions/runs/30280453093/job/90028794924) -> [failure](https://github.com/smartcontractkit/chainlink/actions/runs/30287904253/job/90054795602)
- **Runner**: `['ubuntu-latest']` -> `['ubuntu-latest']`
- **Runner Name**: `GitHub Actions 1012054635` -> `GitHub Actions 1012060787`

| Metric | Trial 1 (Base) | Trial 2 (New) | Delta / Change |
| --- | --- | --- | --- |
| **Duration** | 0:00:03 | 0:00:02 | -0:00:01 (-33.3%) |