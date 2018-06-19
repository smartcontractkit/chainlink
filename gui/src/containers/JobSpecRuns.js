import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import JobRunsList from 'components/JobRunsList'
import PropTypes from 'prop-types'
import React, { Component } from 'react'
import Typography from '@material-ui/core/Typography'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { fetchJobSpecRuns } from 'actions'
import { withStyles } from '@material-ui/core/styles'
import {
  jobRunsCountSelector,
  jobRunsSelector
} from 'selectors'

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

export class JobSpecRuns extends Component {
  componentDidMount () {
    this.props.fetchJobSpecRuns(this.props.jobSpecId)
  }

  render () {
    const {classes, jobSpecId} = this.props

    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>Job ID: {jobSpecId}</BreadcrumbItem>
          <BreadcrumbItem>Runs</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Runs
        </Typography>

        {renderDetails(this.props)}
      </div>
    )
  }
}

const renderLatestRuns = ({jobSpecId, classes, latestJobRuns, jobRunsCount}) => (
  <JobRunsList runs={latestJobRuns} />
)

const renderFetching = () => (
  <div>Fetching...</div>
)

const renderDetails = (props) => {
  if (props.latestJobRuns && props.latestJobRuns.length > 0) {
    return renderLatestRuns(props)
  } else {
    return renderFetching()
  }
}

JobSpecRuns.propTypes = {
  classes: PropTypes.object.isRequired,
  latestJobRuns: PropTypes.array.isRequired,
  jobRunsCount: PropTypes.number.isRequired
}

JobSpecRuns.defaultProps = {
  latestJobRuns: []
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const jobRunsCount = jobRunsCountSelector(state, jobSpecId)
  const latestJobRuns = jobRunsSelector(state, jobSpecId)

  return {
    jobSpecId,
    latestJobRuns,
    jobRunsCount
  }
}

const mapDispatchToProps = (dispatch) => {
  return bindActionCreators({
    fetchJobSpecRuns
  }, dispatch)
}

export const ConnectedJobSpecRuns = connect(mapStateToProps, mapDispatchToProps)(JobSpecRuns)

export default withStyles(styles)(ConnectedJobSpecRuns)
