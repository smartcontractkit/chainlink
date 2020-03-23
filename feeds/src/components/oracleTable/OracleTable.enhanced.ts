import { connect } from 'react-redux'
import OracleTable from './OracleTable.component'
import {
  aggregatorSelectors,
  aggregatorOperations,
} from 'state/ducks/aggregator'
import { AppState } from 'state'

const mapStateToProps = (state: AppState) => ({
  ethGasPrice: state.aggregator.ethGasPrice,
  latestOraclesState: aggregatorSelectors.latestOraclesState(state),
})

const mapDispatchToProps = {
  fetchEthGasPrice: aggregatorOperations.fetchEthGasPrice,
}

export default connect(mapStateToProps, mapDispatchToProps)(OracleTable)
