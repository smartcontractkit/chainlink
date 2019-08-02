import React from 'react'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import { useHooks, useEffect } from 'use-react-hooks'
import KeyValueList from '@chainlink/styleguide/src/components/KeyValueList'
import CardTitle from '@chainlink/styleguide/src/components/Cards/Title'
import JobRunsList from '../../components/JobRuns/List'
import Content from '../../components/Content'
import RegionalNav from './RegionalNav'
import { JobSpecRunsOpts } from '../../api'
import { fetchJob, fetchJobRuns } from '../../actions'
import jobSelector from '../../selectors/job'
import jobRunsByJobIdSelector from '../../selectors/jobRunsByJobId'
import jobsShowRunCountSelector from '../../selectors/jobsShowRunCount'
import { formatInitiators } from '../../utils/jobSpecInitiators'
import matchRouteAndMapDispatchToProps from '../../utils/matchRouteAndMapDispatchToProps'
import TaskList from '../../components/Jobs/TaskList'
import { JobSpec, JobRuns } from '../../../@types/operator_ui'
import { IState } from '../../connectors/redux/reducers'

const renderJobSpec = (job: JobSpec, recentRunsCount: number) => {
  const info = {
    runCount: recentRunsCount,
    initiator: formatInitiators(job.initiators)
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

interface IRecentJobRunsProps {
  job: JobSpec
  recentRuns: JobRuns
  recentRunsCount: number
  showJobRunsCount: number
}

const RecentJobRuns = ({
  job,
  recentRuns,
  recentRunsCount,
  showJobRunsCount
}: IRecentJobRunsProps) => {
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

interface IDetailsProps {
  recentRuns: JobRuns
  recentRunsCount: number
  job?: JobSpec
  showJobRunsCount: number
}

const Details = ({
  job,
  recentRuns,
  recentRunsCount,
  showJobRunsCount
}: IDetailsProps) => {
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
            <Grid item>{renderTaskRuns(job)}</Grid>
            <Grid item>{renderJobSpec(job, recentRunsCount)}</Grid>
          </Grid>
        </Grid>
      </Grid>
    )
  }

  return <div>Fetching...</div>
}

interface IProps {
  jobSpecId: string
  job?: JobSpec
  recentRuns: JobRuns
  recentRunsCount: number
  showJobRunsCount: number
  fetchJob: (id: string) => Promise<any>
  fetchJobRuns: (opts: JobSpecRunsOpts) => Promise<any>
}

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

export const Show = useHooks(
  ({
    jobSpecId,
    job,
    fetchJob,
    fetchJobRuns,
    recentRunsCount,
    recentRuns = [],
    showJobRunsCount = 2
  }: IProps) => {
    useEffect(() => {
      document.title = 'Show Job'
      fetchJob(jobSpecId)
      fetchJobRuns({
        jobSpecId: jobSpecId,
        page: DEFAULT_PAGE,
        size: RECENT_RUNS_COUNT
      })
    }, [])
    return (
      <div>
        {/* TODO: Regional nav should handle job = undefined */}
        {job && <RegionalNav jobSpecId={jobSpecId} job={job} />}
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
)

interface Match {
  params: {
    jobSpecId: string
  }
}

const mapStateToProps = (
  state: IState,
  ownProps: { match: Match; showJobRunsCount: number }
) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const job = jobSelector(state, jobSpecId)
  const recentRuns = jobRunsByJobIdSelector(
    state,
    jobSpecId,
    ownProps.showJobRunsCount
  )
  const recentRunsCount = jobsShowRunCountSelector(state)

  return { jobSpecId, job, recentRuns, recentRunsCount }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob, fetchJobRuns })
)(Show)

export default ConnectedShow
