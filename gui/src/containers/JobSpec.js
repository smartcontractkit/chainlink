import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import PrettyJson from 'components/PrettyJson'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import JobRunsList from 'components/JobRunsList'
import formatInitiators from 'utils/formatInitiators'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { fetchJobSpec } from 'actions'
import { jobSpecSelector, latestJobRunsSelector } from 'selectors'

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
  }
})

const renderJobSpec = ({classes, jobSpec, jobRuns}) => (
  <Grid container spacing={40}>
    <Grid item xs={8}>
      <PaddedCard>
        <Typography variant='title' className={classes.definitionTitle}>
          Definition
        </Typography>
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
                  {jobRuns.length}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </PaddedCard>
    </Grid>
  </Grid>
)

const renderLatestRuns = ({classes, jobRuns}) => (
  <React.Fragment>
    <Typography variant='title' className={classes.lastRun}>
      Last Run
    </Typography>

    <JobRunsList runs={jobRuns} />
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
  jobRuns: PropTypes.array.isRequired,
  jobSpec: PropTypes.object
}

JobSpec.defaultProps = {
  jobRuns: []
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const jobSpec = jobSpecSelector(state, jobSpecId)
  const jobRuns = latestJobRunsSelector(state, jobSpecId)

  return {
    jobSpecId,
    jobSpec,
    jobRuns
  }
}

const mapDispatchToProps = (dispatch) => {
  return bindActionCreators({
    fetchJobSpec
  }, dispatch)
}

export const ConnectedJobSpec = connect(mapStateToProps, mapDispatchToProps)(JobSpec)

export default withStyles(styles)(ConnectedJobSpec)
