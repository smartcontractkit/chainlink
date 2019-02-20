import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import JobRunsList from 'components/JobRuns/List'
import Link from 'components/Link'
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

const styles = theme => ({
  lastRun: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  showMore: {
    marginTop: theme.spacing.unit * 3,
    marginLeft: theme.spacing.unit * 3,
    display: 'block'
  }
})

const renderJobSpec = ({ job }) => {
  const info = {
    runCount: job.runs && job.runs.length,
    initiator: formatInitiators(job.initiators)
  }

  return <KeyValueList entries={Object.entries(info)} titleize />
}

const renderLatestRuns = ({
  job,
  classes,
  latestJobRuns,
  showJobRunsCount
}) => (
  <React.Fragment>
    <Card>
      <CardTitle divider>Recent Job Runs</CardTitle>
      <JobRunsList jobSpecId={job.id} runs={latestJobRuns} />
    </Card>
    {job.runs && job.runs.length > showJobRunsCount && (
      <Link to={`/jobs/${job.id}/runs`} className={classes.showMore}>
        Show More
      </Link>
    )}
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
          {renderJobSpec(props)}
        </Grid>
      </Grid>
    )
  }

  return <div>Fetching...</div>
}

export const Show = useHooks(props => {
  useEffect(() => {
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
  classes: PropTypes.object.isRequired,
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

export default withStyles(styles)(ConnectedShow)
