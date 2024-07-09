package rollups

/* ABIs for Arbitrum Gas Info and Node Interface precompile contract methods needed for the L1 oracle */
// ABI found at https://arbiscan.io/address/0x000000000000000000000000000000000000006C#code
const GetL1BaseFeeEstimateAbiString = `[{"inputs":[],"name":"getL1BaseFeeEstimate","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// ABI found at https://arbiscan.io/address/0x00000000000000000000000000000000000000C8#code
const GasEstimateL1ComponentAbiString = `[{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"bool","name":"contractCreation","type":"bool"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"gasEstimateL1Component","outputs":[{"internalType":"uint64","name":"gasEstimateForL1","type":"uint64"},{"internalType":"uint256","name":"baseFee","type":"uint256"},{"internalType":"uint256","name":"l1BaseFeeEstimate","type":"uint256"}],"stateMutability":"payable","type":"function"}]`

/* ABIs for Optimism, Scroll, and Kroma precompile contract methods needed for the L1 oracle */
// All ABIs found at https://optimistic.etherscan.io/address/0xc0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3000f#code
const L1BaseFeeAbiString = `[{"inputs":[],"name":"l1BaseFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
const GetL1FeeAbiString = `[{"inputs":[{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"getL1Fee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// ABIs for OP Stack GasPriceOracle methods needed to calculated encoded gas price
const OPIsEcotoneAbiString = `[{"inputs":[],"name":"isEcotone","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
const OPIsFjordAbiString = `[{"inputs":[],"name":"isFjord","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
const OPBaseFeeScalarAbiString = `[{"inputs":[],"name":"baseFeeScalar","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"}]`
const OPBlobBaseFeeAbiString = `[{"inputs":[],"name":"blobBaseFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
const OPBlobBaseFeeScalarAbiString = `[{"inputs":[],"name":"blobBaseFeeScalar","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"}]`
const OPDecimalsAbiString = `[{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"pure","type":"function"}]`
