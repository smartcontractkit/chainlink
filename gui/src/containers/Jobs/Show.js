import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import PaddedCard from 'components/PaddedCard'
import JobRunsList from 'components/JobRuns/List'
import Link from 'components/Link'
import Content from 'components/Content'
import RegionalNav from 'components/Jobs/RegionalNav'
import TimeAgo from 'components/TimeAgo'
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
  return (
    <Grid container spacing={40}>
      <Grid item xs={4}>
        <PaddedCard>
          <Grid container spacing={16}>
            <Grid item xs={12}>
              <Typography variant='subtitle1' color='textSecondary'>ID</Typography>
              <Typography variant='body1' color='inherit'>
                {job.id}
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography variant='subtitle1' color='textSecondary'>Created</Typography>
              <Typography variant='body1' color='inherit'>
                <TimeAgo>{job.createdAt}</TimeAgo>
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Grid container spacing={16}>
                <Grid item xs={6}>
                  <Typography variant='subtitle1' color='textSecondary'>Initiator</Typography>
                  <Typography variant='body1' color='inherit'>
                    {formatInitiators(job.initiators)}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant='subtitle1' color='textSecondary'>Run Count</Typography>
                  <Typography variant='body1' color='inherit'>
                    {job.runs && job.runs.length}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </PaddedCard>
      </Grid>
    </Grid>
  )
}

const renderLatestRuns = ({ job, classes, latestJobRuns, showJobRunsCount }) => (
  <React.Fragment>
    <Typography variant='h5' className={classes.lastRun}>
      Last Run
    </Typography>

    <Card>
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
      <React.Fragment>
        {renderJobSpec(props)}
        {renderLatestRuns(props)}
      </React.Fragment>
    )
  }

  return <div>Fetching...</div>
}

export const Show = useHooks(props => {
  useEffect(() => { fetchJob(jobSpecId) }, [])
  const { jobSpecId, job, fetchJob } = props
  return (
    <div>
      <RegionalNav jobSpecId={jobSpecId} job={job} />
      <Content>
        {renderDetails(props)}
      </Content>
    </div>
  )
}
)

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
  const latestJobRuns = jobRunsByJobIdSelector(state, jobSpecId, ownProps.showJobRunsCount)

  return {jobSpecId, job, latestJobRuns}
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchJob})
)(Show)

export default withStyles(styles)(ConnectedShow)
