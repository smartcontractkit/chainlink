import React from 'react'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import { useHooks, useEffect } from 'use-react-hooks'
import JobRunsList from '../../components/JobRuns/List'
import KeyValueList from '../../components/KeyValueList'
import Content from '../../components/Content'
import RegionalNav from '../../components/Jobs/RegionalNav'
import CardTitle from '../../components/Cards/Title'
import { fetchJob } from '../../actions'
import jobSelector from '../../selectors/job'
import jobRunsByJobIdSelector from '../../selectors/jobRunsByJobId'
import jobsShowRunCountSelector from '../../selectors/jobsShowRunCount'
import { formatInitiators } from '../../utils/jobSpecInitiators'
import matchRouteAndMapDispatchToProps from '../../utils/matchRouteAndMapDispatchToProps'
import TaskRuns from './TaskRuns'

const renderJobSpec = job => {
  const info = {
    runCount: job.runs && job.runs.length,
    initiator: formatInitiators(job.initiators)
  }

  return <KeyValueList entries={Object.entries(info)} titleize />
}

const renderTaskRuns = job => (
  <Card>
    <CardTitle divider>Task List</CardTitle>
    <TaskRuns taskRuns={job.tasks} />
  </Card>
)

const renderLatestRuns = (
  job,
  recentRuns,
  recentRunsCount,
  showJobRunsCount
) => (
  <React.Fragment>
    <Card>
      <CardTitle divider>Recent Job Runs</CardTitle>
      <JobRunsList
        jobSpecId={job.id}
        runs={recentRuns}
        count={recentRunsCount}
        showJobRunsCount={showJobRunsCount}
      />
    </Card>
  </React.Fragment>
)

interface IDetailsProps {
  recentRuns: any[]
  recentRunsCount: number
  job?: any
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
          {renderLatestRuns(job, recentRuns, recentRunsCount, showJobRunsCount)}
        </Grid>
        <Grid item xs={4}>
          <Grid container direction="column">
            <Grid item>{renderTaskRuns(job)}</Grid>
            <Grid item>{renderJobSpec(job)}</Grid>
          </Grid>
        </Grid>
      </Grid>
    )
  }

  return <div>Fetching...</div>
}

interface IProps {
  jobSpecId: string
  job?: any
  recentRuns: any[]
  recentRunsCount: number
  showJobRunsCount: number
  fetchJob: (string) => Promise<any>
}

export const Show = useHooks(
  ({
    jobSpecId,
    job,
    fetchJob,
    recentRunsCount,
    recentRuns = [],
    showJobRunsCount = 2
  }: IProps) => {
    useEffect(() => {
      document.title = 'Show Job'
      fetchJob(jobSpecId)
    }, [])

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
)

const mapStateToProps = (state, ownProps) => {
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
  matchRouteAndMapDispatchToProps({ fetchJob })
)(Show)

export default ConnectedShow
