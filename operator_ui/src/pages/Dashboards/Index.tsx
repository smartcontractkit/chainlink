import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import WarningIcon from '@material-ui/icons/Warning'
import { withStyles, WithStyles, Theme } from '@material-ui/core/styles'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import RecentlyCreatedJobs from 'components/Jobs/RecentlyCreated'
import Footer from 'components/Footer'
import Content from 'components/Content'
import Paper from '@material-ui/core/Paper'
import {
  fetchRecentJobRuns,
  fetchRecentlyCreatedJobs,
  fetchAccountBalance,
} from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'
import dashboardJobRunsCountSelector from 'selectors/dashboardJobRunsCount'
import recentJobRunsSelector from 'selectors/recentJobRuns'
import recentlyCreatedJobsSelector from 'selectors/recentlyCreatedJobs'

const styles = (theme: Theme) => ({
  root: {
    ...theme.mixins.gutters(),
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
    backgroundColor: 'rgb(255, 213, 153)',
  },
})

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
  classes: WithStyles<typeof styles>['classes']
}

export const Index = ({
  recentJobRunsCount = 2,
  recentlyCreatedPageSize = 2,
  classes,
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
        {recentlyCreatedJobs && recentlyCreatedJobs.length > 0 && (
          <Grid item xs={12}>
            <Paper className={classes.root}>
              <p
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  margin: 0,
                  marginBottom: 8,
                }}
              >
                <WarningIcon style={{ marginRight: 8, color: '#ff9800' }} />
                <Typography variant="h5">Found legacy job specs</Typography>
              </p>
              <Typography>
                The JSON style of job spec is now deprecated and support for
                jobs using this format will be REMOVED in an upcoming release.
                You should migrate all these jobs to V2 (TOML) format. For help
                doing this, please check the{' '}
                <a href="https://docs.chain.link/docs/jobs/">docs</a>. To test
                your node to see how it would behave after support for these
                jobs is removed, you may set ENABLE_LEGACY_JOB_PIPELINE=false
              </Typography>
            </Paper>
          </Grid>
        )}
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

export default withStyles(styles)(Index)
