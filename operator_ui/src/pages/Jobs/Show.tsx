import { CardTitle, KeyValueList } from '@chainlink/styleguide'
import {
  createStyles,
  Theme,
  Typography,
  WithStyles,
  withStyles,
} from '@material-ui/core'
import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import { RouteComponentProps } from 'react-router-dom'
import { fetchJob, fetchJobRuns } from 'actionCreators'
import Content from 'components/Content'
import JobRunsList from 'components/JobRuns/List'
import TaskList from 'components/Jobs/TaskList'
import { AppState } from 'src/reducers'
import { JobRuns, JobSpec } from 'operator_ui'
import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import jobSelector from 'selectors/job'
import jobRunsByJobIdSelector from 'selectors/jobRunsByJobId'
import jobsShowRunCountSelector from 'selectors/jobsShowRunCount'
import { GWEI_PER_TOKEN } from 'utils/constants'
import formatMinPayment from 'utils/formatWeiAsset'
import { formatInitiators } from 'utils/jobSpecInitiators'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import RegionalNav from './RegionalNav'

const renderJobSpec = (job: JobSpec, recentRunsCount: number) => {
  const info = {
    runCount: recentRunsCount,
    initiator: formatInitiators(job.initiators),
    minimumPayment: `${formatMinPayment(job.minPayment) || 0} Link`,
  }

  return (
    <KeyValueList showHead={false} entries={Object.entries(info)} titleize />
  )
}

const renderTaskRuns = (job: JobSpec) => (
  <Card>
    <CardTitle divider>Task List</CardTitle>
    <TaskList tasks={job.tasks} />
  </Card>
)

interface RecentJobRunsProps {
  job: JobSpec
  recentRuns: JobRuns
  recentRunsCount: number
  showJobRunsCount: number
}

const totalLinkEarned = (job: JobSpec) => {
  const zero = '0.000000'
  const unformatted = job.earnings && (job.earnings / GWEI_PER_TOKEN).toString()
  const formatted =
    unformatted &&
    (unformatted.length >= 3 ? unformatted : (unformatted + '.').padEnd(8, '0'))
  return formatted || zero
}

const chartCardStyles = (theme: Theme) =>
  createStyles({
    wrapper: {
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 2,
      marginBottom: theme.spacing.unit * 2,
    },
    paymentText: {
      color: theme.palette.secondary.main,
      fontWeight: 450,
    },
    earnedText: {
      color: theme.palette.text.secondary,
      fontSize: theme.spacing.unit * 2,
    },
  })

interface ChartProps extends WithStyles<typeof chartCardStyles> {
  job: JobSpec
}

const ChartArea = withStyles(chartCardStyles)(
  ({ classes, job }: ChartProps) => (
    <Card>
      <Grid item className={classes.wrapper}>
        <Typography className={classes.paymentText} variant="h5">
          Link Payment
        </Typography>
        <Typography className={classes.earnedText}>
          {totalLinkEarned(job)}
        </Typography>
      </Grid>
    </Card>
  ),
)

const RecentJobRuns: React.FC<RecentJobRunsProps> = ({
  job,
  recentRuns,
  recentRunsCount,
  showJobRunsCount,
}) => {
  return (
    <Card>
      <CardTitle divider>Recent Job Runs</CardTitle>

      <JobRunsList
        jobSpecId={job.id}
        runs={recentRuns}
        count={recentRunsCount}
        showJobRunsCount={showJobRunsCount}
      />
    </Card>
  )
}

interface DetailsProps {
  recentRuns: JobRuns
  recentRunsCount: number
  job?: JobSpec
  showJobRunsCount: number
}

const Details: React.FC<DetailsProps> = ({
  job,
  recentRuns,
  recentRunsCount,
  showJobRunsCount,
}) => {
  if (job) {
    return (
      <Grid container spacing={24}>
        <Grid item xs={8}>
          <RecentJobRuns
            job={job}
            recentRuns={recentRuns}
            recentRunsCount={recentRunsCount}
            showJobRunsCount={showJobRunsCount}
          />
        </Grid>
        <Grid item xs={4}>
          <Grid container direction="column">
            <Grid item xs>
              <ChartArea job={job} />
            </Grid>
            <Grid item xs>
              {renderTaskRuns(job)}
            </Grid>
            <Grid item xs>
              {renderJobSpec(job, recentRunsCount)}
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    )
  }

  return <div>Fetching...</div>
}

type OwnProps = {
  showJobRunsCount: number
}

export type ShowComponentProps = OwnProps &
  RouteComponentProps<{
    jobSpecId: string
  }>

interface Props {
  jobSpecId: string
  job?: JobSpec
  recentRuns: JobRuns
  recentRunsCount: number
  showJobRunsCount: number
  fetchJob: (id: string) => Promise<any>
  fetchJobRuns: (opts: any) => Promise<any>
}

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const Show: React.FC<Props> = ({
  jobSpecId,
  job,
  fetchJob,
  fetchJobRuns,
  recentRunsCount,
  recentRuns = [],
  showJobRunsCount = 2,
}) => {
  useEffect(() => {
    document.title = 'Show Job'
    fetchJob(jobSpecId)
    fetchJobRuns({
      jobSpecId,
      page: DEFAULT_PAGE,
      size: RECENT_RUNS_COUNT,
    })
  }, [fetchJob, fetchJobRuns, jobSpecId])
  return (
    <div>
      <RegionalNav jobSpecId={jobSpecId} job={job} />
      <Content>
        <Details
          job={job}
          recentRuns={recentRuns}
          recentRunsCount={recentRunsCount}
          showJobRunsCount={showJobRunsCount}
        />
      </Content>
    </div>
  )
}

const mapStateToProps = (state: AppState, ownProps: ShowComponentProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const job = jobSelector(state, jobSpecId)
  const recentRuns = jobRunsByJobIdSelector(
    state,
    jobSpecId,
    ownProps.showJobRunsCount,
  )
  const recentRunsCount = jobsShowRunCountSelector(state)

  return { jobSpecId, job, recentRuns, recentRunsCount }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob, fetchJobRuns }),
)(Show)

export default ConnectedShow
