import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import RecentActivity from 'components/Dashboards/RecentActivity'
import TokenBalance from 'components/TokenBalance'
import RecentlyCreatedJobs from 'components/Jobs/RecentlyCreated'
import Footer from 'components/Footer'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import {
  fetchRecentJobRuns,
  fetchRecentlyCreatedJobs,
  fetchAccountBalance
} from 'actions'
import accountBalanceSelector from 'selectors/accountBalance'
import recentJobRunsSelector from 'selectors/recentJobRuns'
import recentlyCreatedJobsSelector from 'selectors/recentlyCreatedJobs'

const styles = theme => ({
  index: {
    paddingTop: theme.spacing.unit * 5,
    paddingBottom: theme.spacing.unit * 5
  }
})

export class Index extends Component {
  componentDidMount () {
    const {props} = this
    props.fetchAccountBalance()
    props.fetchRecentJobRuns(props.recentJobRunsCount)
    props.fetchRecentlyCreatedJobs(props.recentlyCreatedPageSize)
  }

  render () {
    const {props} = this
    return (
      <div className={props.classes.index}>
        <Grid container spacing={40}>
          <Grid item xs={9}>
            <RecentActivity runs={props.recentJobRuns} />
          </Grid>
          <Grid item xs={3}>
            <Grid container spacing={24}>
              <Grid item xs={12}>
                <TokenBalance
                  title='Link Balance'
                  value={props.accountBalance && props.accountBalance.linkBalance}
                />
              </Grid>
              <Grid item xs={12}>
                <TokenBalance
                  title='Ether Balance'
                  value={props.accountBalance && props.accountBalance.ethBalance}
                />
              </Grid>
              <Grid item xs={12}>
                <RecentlyCreatedJobs jobs={props.recentlyCreatedJobs} />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
        <Footer />
      </div>
    )
  }
}

Index.propTypes = {
  accountBalance: PropTypes.object,
  recentJobRunsCount: PropTypes.number.isRequired,
  recentJobRuns: PropTypes.array,
  recentlyCreatedJobs: PropTypes.array,
  recentlyCreatedPageSize: PropTypes.number
}

Index.defaultProps = {
  recentJobRunsCount: 2,
  recentlyCreatedPageSize: 2
}

const mapStateToProps = state => {
  return {
    accountBalance: accountBalanceSelector(state),
    recentJobRuns: recentJobRunsSelector(state),
    recentlyCreatedJobs: recentlyCreatedJobsSelector(state)
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({
    fetchAccountBalance,
    fetchRecentJobRuns,
    fetchRecentlyCreatedJobs
  })
)(Index)

export default withStyles(styles)(ConnectedIndex)
