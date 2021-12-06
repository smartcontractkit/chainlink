import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import List from 'components/Transactions/List'
import Content from 'components/Content'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import transactionsSelector from 'selectors/transactions'
import { fetchTransactions } from 'actionCreators'
import { Heading1 } from 'src/components/Heading/Heading1'

export const Index = (props) => {
  document.title = 'Transactions'
  return (
    <Content>
      <Grid container spacing={32}>
        <Grid item xs={12}>
          <Heading1>Transactions</Heading1>
        </Grid>
        <Grid item xs={12}>
          <List
            transactions={props.transactions}
            count={props.count}
            pageSize={props.pageSize}
            fetchTransactions={props.fetchTransactions}
            history={props.history}
            match={props.match}
          />
        </Grid>
      </Grid>
    </Content>
  )
}

Index.propTypes = {
  count: PropTypes.number.isRequired,
  transactions: PropTypes.array,
  pageSize: PropTypes.number,
}

Index.defaultProps = {
  pageSize: 10,
}

const mapStateToProps = (state) => {
  return {
    count: state.transactionsIndex.count,
    transactions: transactionsSelector(state),
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchTransactions }),
)(Index)

export default ConnectedIndex
