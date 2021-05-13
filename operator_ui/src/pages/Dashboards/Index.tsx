import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import RecentlyCreatedJobs from 'components/Jobs/RecentlyCreated'
import Footer from 'components/Footer'
import Content from 'components/Content'
import {
  fetchRecentJobRuns,
  fetchRecentlyCreatedJobs,
  fetchAccountBalance,
} from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'
import dashboardJobRunsCountSelector from 'selectors/dashboardJobRunsCount'
import recentJobRunsSelector from 'selectors/recentJobRuns'
import recentlyCreatedJobsSelector from 'selectors/recentlyCreatedJobs'

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
}

export const Index = ({
  recentJobRunsCount = 2,
  recentlyCreatedPageSize = 2,
}: Props) => {
  const dispatch = useDispatch()
  const accountBalance = useSelector(accountBalanceSelector)
  const jobRunsCount = useSelector(dashboardJobRunsCountSelector)
  const recentJobRuns = useSelector(recentJobRunsSelector)
  const recentlyCreatedJobs = useSelector(recentlyCreatedJobsSelector)

  useEffect(() => {
    document.title = 'Dashboard'
  }, [])

  useEffect(() => {
    dispatch(fetchAccountBalance())
  }, [dispatch])

  useEffect(() => {
    dispatch(fetchRecentJobRuns(recentJobRunsCount))
  }, [dispatch, recentJobRunsCount])

  useEffect(() => {
    dispatch(fetchRecentlyCreatedJobs(recentlyCreatedPageSize))
  }, [dispatch, recentlyCreatedPageSize])

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

export default Index
