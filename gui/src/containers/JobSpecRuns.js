import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import JobRunsList from 'components/JobRunsList'
import PropTypes from 'prop-types'
import React, { Component } from 'react'
import Typography from '@material-ui/core/Typography'
import TablePagination from '@material-ui/core/TablePagination'
import Card from '@material-ui/core/Card'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { fetchJobSpecRuns } from 'actions'
import { withStyles } from '@material-ui/core/styles'
import { jobRunsCountSelector, jobRunsSelector } from 'selectors'
import { IconButton } from '@material-ui/core';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import FirstPageIcon from '@material-ui/icons/FirstPage';
import LastPageIcon from '@material-ui/icons/LastPage';

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  customButtons: {
    flexShrink: 0,
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing.unit * 2.5,
  }
})

export class JobSpecRuns extends Component {
  constructor(props) {
    super(props)
    this.state = {
      page: 0
    }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  componentDidMount() {
    const { jobSpecId, pageSize, fetchJobSpecRuns } = this.props
    if (this.props.match.params.jobRunsPage) {
      const START_PAGE = this.props.match.params.jobRunsPage
      this.setState({ page: START_PAGE - 1 })
      fetchJobSpecRuns(jobSpecId, START_PAGE, pageSize)
    } else {
      this.setState({ page: 0})
      fetchJobSpecRuns(jobSpecId, 1, pageSize)
    }
  }

  handleChangePage(e, page) {
    const { fetchJobSpecRuns, jobSpecId, pageSize } = this.props
    this.props.history.replace(`/job_specs/${jobSpecId}/runs/page/${page+1}`)
    fetchJobSpecRuns(jobSpecId, page + 1, pageSize)
    this.setState({ page })
  }
  render() {
    const { classes, jobSpecId } = this.props

    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href="/">Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href={`/job_specs/${jobSpecId}`}>Job ID: {jobSpecId}</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>Runs</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant="display2" color="inherit" className={classes.title}>
          Runs
        </Typography>

        {renderDetails(this.props, this.state, this.handleChangePage)}
      </div>
    )
  }
}

const renderLatestRuns = ({ jobSpecId, latestJobRuns, jobRunsCount, pageSize }, state, handleChangePage) => {
  return(
    <Card>
      <JobRunsList jobSpecId={jobSpecId} runs={latestJobRuns} />
      <TablePagination
        component="div"
        count={jobRunsCount}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[pageSize]}
        page={state.page}
        onChangePage={handleChangePage}
        onChangeRowsPerPage={() => {} /* handler required by component, so make it a no-op */}
      />
    </Card>
  )
}

const renderFetching = () => <div>Fetching...</div>

const renderDetails = (props, state, handleChangePage) => {
  if (props.latestJobRuns && props.latestJobRuns.length > 0) {
    return renderLatestRuns(props, state, handleChangePage)
  } else {
    return renderFetching()
  }
}

JobSpecRuns.propTypes = {
  classes: PropTypes.object.isRequired,
  latestJobRuns: PropTypes.array.isRequired,
  jobRunsCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired
}

JobSpecRuns.defaultProps = {
  latestJobRuns: [],
  pageSize: 10
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

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchJobSpecRuns
    },
    dispatch
  )
}

export const ConnectedJobSpecRuns = connect(mapStateToProps, mapDispatchToProps)(JobSpecRuns)

export default withStyles(styles)(ConnectedJobSpecRuns)
