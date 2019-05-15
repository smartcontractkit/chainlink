import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import JobRunsList from 'components/JobRuns/List'
import KeyValueList from 'components/KeyValueList'
import Content from 'components/Content'
import RegionalNav from 'components/Jobs/RegionalNav'
import CardTitle from 'components/Cards/Title'
import { fetchJob } from 'actions'
import jobSelector from 'selectors/job'
import jobRunsByJobIdSelector from 'selectors/jobRunsByJobId'
import { formatInitiators } from 'utils/jobSpecInitiators'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { useHooks, useEffect } from 'use-react-hooks'

const renderJobSpec = ({ job }) => {
  const info = {
    runCount: job.runs && job.runs.length,
    initiator: formatInitiators(job.initiators)
  }

  return <KeyValueList entries={Object.entries(info)} titleize />
}

const renderTaskList = ({ job }) => {
  const list = Object.assign(
    {},
    ...job.tasks.map(item => ({ [item.type]: '' }))
  )
  return (
    <Card>
      <CardTitle divider>Task List</CardTitle>
      <KeyValueList entries={Object.entries(list)} titleize />
    </Card>
  )
}

const renderLatestRuns = ({ job, latestJobRuns, showJobRunsCount }) => (
  <React.Fragment>
    <Card>
      <CardTitle divider>Recent Job Runs</CardTitle>
      <JobRunsList
        jobSpecId={job.id}
        jobRuns={job.runs}
        runs={latestJobRuns}
        showJobRunsCount={showJobRunsCount}
      />
    </Card>
  </React.Fragment>
)

const renderDetails = props => {
  if (props.job) {
    return (
      <Grid container spacing={24}>
        <Grid item xs={8}>
          {renderLatestRuns(props)}
        </Grid>
        <Grid item xs={4}>
          <Grid container direction="column">
            <Grid item>{renderTaskList(props)}</Grid>
            <Grid item>{renderJobSpec(props)}</Grid>
          </Grid>
        </Grid>
      </Grid>
    )
  }

  return <div>Fetching...</div>
}

export const Show = useHooks(props => {
  useEffect(() => {
    document.title = 'Show Job'
    fetchJob(jobSpecId)
  }, [])
  const { jobSpecId, job, fetchJob } = props
  return (
    <div>
      <RegionalNav jobSpecId={jobSpecId} job={job} />
      <Content>{renderDetails(props)}</Content>
    </div>
  )
})

Show.propTypes = {
  latestJobRuns: PropTypes.array.isRequired,
  job: PropTypes.object,
  showJobRunsCount: PropTypes.number
}

Show.defaultProps = {
  latestJobRuns: [],
  showJobRunsCount: 2
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const job = jobSelector(state, jobSpecId)
  const latestJobRuns = jobRunsByJobIdSelector(
    state,
    jobSpecId,
    ownProps.showJobRunsCount
  )

  return { jobSpecId, job, latestJobRuns }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob })
)(Show)

export default ConnectedShow
