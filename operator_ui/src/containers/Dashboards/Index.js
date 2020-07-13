import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import RecentlyCreatedJobs from 'components/Jobs/RecentlyCreated'
import Footer from 'components/Footer'
import Content from 'components/Content'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import {
  fetchRecentJobRuns,
  fetchRecentlyCreatedJobs,
  fetchAccountBalance,
} from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'
import dashboardJobRunsCountSelector from 'selectors/dashboardJobRunsCount'
import recentJobRunsSelector from 'selectors/recentJobRuns'
import recentlyCreatedJobsSelector from 'selectors/recentlyCreatedJobs'

export const Index = ({
  accountBalance,
  fetchAccountBalance,
  fetchRecentJobRuns,
  fetchRecentlyCreatedJobs,
  jobRunsCount,
  recentJobRuns,
  recentJobRunsCount,
  recentlyCreatedJobs,
  recentlyCreatedPageSize,
}) => {
  useEffect(() => {
    document.title = 'Dashboard'
    fetchAccountBalance()
    fetchRecentJobRuns(recentJobRunsCount)
    fetchRecentlyCreatedJobs(recentlyCreatedPageSize)
  }, [
    fetchAccountBalance,
    fetchRecentJobRuns,
    fetchRecentlyCreatedJobs,
    recentJobRunsCount,
    recentlyCreatedPageSize,
  ])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Activity
            runs={recentJobRuns}
            pageSize={recentJobRunsCount}
            count={jobRunsCount}
          />
        </Grid>
        <Grid item xs={3}>
          <Grid container>
            <Grid item xs={12}>
              <TokenBalanceCard
                title="Link Balance"
                value={accountBalance?.linkBalance}
              />
            </Grid>
            <Grid item xs={12}>
              <TokenBalanceCard
                title="Ether Balance"
                value={accountBalance?.ethBalance}
              />
            </Grid>
            <Grid item xs={12}>
              <RecentlyCreatedJobs jobs={recentlyCreatedJobs} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Footer />
    </Content>
  )
}

Index.propTypes = {
  accountBalance: PropTypes.object,
  recentJobRunsCount: PropTypes.number.isRequired,
  jobRunsCount: PropTypes.number,
  recentJobRuns: PropTypes.array,
  recentlyCreatedJobs: PropTypes.array,
  recentlyCreatedPageSize: PropTypes.number,
}

Index.defaultProps = {
  recentJobRunsCount: 2,
  recentlyCreatedPageSize: 2,
}

const mapStateToProps = state => {
  return {
    accountBalance: accountBalanceSelector(state),
    jobRunsCount: dashboardJobRunsCountSelector(state),
    recentJobRuns: recentJobRunsSelector(state),
    recentlyCreatedJobs: recentlyCreatedJobsSelector(state),
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({
    fetchAccountBalance,
    fetchRecentJobRuns,
    fetchRecentlyCreatedJobs,
  }),
)(Index)

export default ConnectedIndex
