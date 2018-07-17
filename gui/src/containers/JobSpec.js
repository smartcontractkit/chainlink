import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import PrettyJson from 'components/PrettyJson'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import Card from '@material-ui/core/Card'
import JobRunsList from 'components/JobRunsList'
import formatInitiators from 'utils/formatInitiators'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import Link from 'components/Link'
import Copy from 'components/CopyJobSpec'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { fetchJobSpec } from 'actions'
import {
  jobSpecSelector,
  jobRunsSelector,
  jobRunsCountSelector
} from 'selectors'
import { LATEST_JOB_RUNS_COUNT } from 'connectors/redux/reducers/jobRuns'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  definitionTitle: {
    marginBottom: theme.spacing.unit * 3
  },
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
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

const renderJobSpec = ({classes, jobSpec, latestJobRuns, jobRunsCount}) => (
  <Grid container spacing={40}>
    <Grid item xs={8}>
      <PaddedCard>
        <Typography variant='title' className={classes.definitionTitle}>
          Definition
        </Typography>
        <Copy JobSpec={JSON.stringify(jobSpecDefinition(jobSpec), null, "\t")}/>     
	<PrettyJson object={jobSpecDefinition(jobSpec)} />
      </PaddedCard>
    </Grid>
    <Grid item xs={4}>
      <PaddedCard>
        <Grid container spacing={16}>
          <Grid item xs={12}>
            <Typography variant='subheading' color='textSecondary'>ID</Typography>
            <Typography variant='body1' color='inherit'>
              {jobSpec.id}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant='subheading' color='textSecondary'>Created</Typography>
            <Typography variant='body1' color='inherit'>
              {jobSpec.createdAt}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Grid container spacing={16}>
              <Grid item xs={6}>
                <Typography variant='subheading' color='textSecondary'>Initiator</Typography>
                <Typography variant='body1' color='inherit'>
                  {formatInitiators(jobSpec.initiators)}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant='subheading' color='textSecondary'>Run Count</Typography>
                <Typography variant='body1' color='inherit'>
                  {jobRunsCount}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </PaddedCard>
    </Grid>
  </Grid>
)

const renderLatestRuns = ({jobSpecId, classes, latestJobRuns, jobRunsCount}) => (
  <React.Fragment>
    <Typography variant='title' className={classes.lastRun}>
      Last Run
    </Typography>

    <Card>
      <JobRunsList jobSpecId={jobSpecId} runs={latestJobRuns} />
    </Card>
    {jobRunsCount > LATEST_JOB_RUNS_COUNT &&
      <Link to={`/job_specs/${jobSpecId}/runs`} className={classes.showMore}>
        Show More
      </Link>
    }
  </React.Fragment>
)

const renderFetching = () => (
  <div>Fetching...</div>
)

const renderDetails = (props) => {
  if (props.jobSpec) {
    return (
      <React.Fragment>
        {renderJobSpec(props)}
        {renderLatestRuns(props)}
      </React.Fragment>
    )
  } else {
    return renderFetching()
  }
}

export class JobSpec extends Component {
  componentDidMount () {
    this.props.fetchJobSpec(this.props.jobSpecId)
  }

  render () {
    const {classes, jobSpecId} = this.props

    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>Job ID: {jobSpecId}</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Job Spec Detail
        </Typography>

        {renderDetails(this.props)}
      </div>
    )
  }
}

JobSpec.propTypes = {
  classes: PropTypes.object.isRequired,
  latestJobRuns: PropTypes.array.isRequired,
  jobSpec: PropTypes.object,
  jobRunsCount: PropTypes.number
}

JobSpec.defaultProps = {
  latestJobRuns: []
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const jobSpec = jobSpecSelector(state, jobSpecId)
  const jobRunsCount = jobRunsCountSelector(state, jobSpecId)
  const latestJobRuns = jobRunsSelector(state)

  return {
    jobSpecId,
    jobSpec,
    latestJobRuns,
    jobRunsCount
  }
}

const mapDispatchToProps = (dispatch) => {
  return bindActionCreators({
    fetchJobSpec
  }, dispatch)
}

export const ConnectedJobSpec = connect(mapStateToProps, mapDispatchToProps)(JobSpec)

export default withStyles(styles)(ConnectedJobSpec)
