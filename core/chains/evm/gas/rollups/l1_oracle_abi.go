package rollups

/* ABIs for Arbitrum Gas Info and Node Interface precompile contract methods needed for the L1 oracle */
const GetL1BaseFeeEstimateAbiString = `[{"inputs":[],"name":"getL1BaseFeeEstimate","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
const GasEstimateL1ComponentAbiString = `[{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"bool","name":"contractCreation","type":"bool"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"gasEstimateL1Component","outputs":[{"internalType":"uint64","name":"gasEstimateForL1","type":"uint64"},{"internalType":"uint256","name":"baseFee","type":"uint256"},{"internalType":"uint256","name":"l1BaseFeeEstimate","type":"uint256"}],"stateMutability":"payable","type":"function"}]`

/* ABIs for Optimism, Scroll, and Kroma precompile contract methods needed for the L1 oracle */
const L1BaseFeeAbiString = `[{"inputs":[],"name":"l1BaseFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
const GetL1FeeAbiString = `[{"inputs":[{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"getL1Fee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
