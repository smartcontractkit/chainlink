import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import JobRunsList from 'components/JobRunsList'
import PropTypes from 'prop-types'
import React, { Component } from 'react'
import Typography from '@material-ui/core/Typography'
import TablePagination from '@material-ui/core/TablePagination'
import Card from '@material-ui/core/Card'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { connect } from 'react-redux'
import { fetchJobSpecRuns } from 'actions'
import { withStyles } from '@material-ui/core/styles'
import { jobRunsCountSelector, jobRunsSelector } from 'selectors'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'

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
  constructor (props) {
    super(props)
    this.state = {
      page: 0
    }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  componentDidMount () {
    const { jobSpecId, pageSize, fetchJobSpecRuns } = this.props
    const queryPage = this.props.match ? parseInt(this.props.match.params.jobRunsPage, 10) || FIRST_PAGE : FIRST_PAGE
    this.setState({ page: queryPage })
    fetchJobSpecRuns(jobSpecId, queryPage, pageSize)
  }

  componentDidUpdate (prevProps) {
    const prevJobRunsPage = prevProps.match.params.jobRunsPage
    const currentJobRunsPage = this.props.match.params.jobRunsPage

    if (prevJobRunsPage !== currentJobRunsPage) {
      const { pageSize, fetchJobSpecRuns, jobSpecId } = this.props
      this.setState({ page: parseInt(currentJobRunsPage, 10) || FIRST_PAGE })
      fetchJobSpecRuns(jobSpecId, parseInt(currentJobRunsPage, 10) || FIRST_PAGE, pageSize)
    }
  }

  handleChangePage (e, page) {
    const { fetchJobSpecRuns, jobSpecId, pageSize } = this.props
    fetchJobSpecRuns(jobSpecId, page, pageSize)
    this.setState({ page })
  }
  render () {
    const { classes, jobSpecId } = this.props

    return (
      <div>
        <Breadcrumb className={classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href={`/job_specs/${jobSpecId}`}>Job ID: {jobSpecId}</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>Runs</BreadcrumbItem>
        </Breadcrumb>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Runs
        </Typography>

        {renderDetails(this.props, this.state, this.handleChangePage)}
      </div>
    )
  }
}

const renderLatestRuns = (props, state, handleChangePage) => {
  const { jobSpecId, latestJobRuns, jobRunsCount, pageSize } = props
  const TableButtonsWithProps = () => (
    <TableButtons
      {...props}
      count={jobRunsCount}
      onChangePage={handleChangePage}
      page={state.page}
      specID={jobSpecId}
      rowsPerPage={pageSize}
      replaceWith={`/job_specs/${jobSpecId}/runs/page`}
    />
  )
  return (
    <Card>
      <JobRunsList jobSpecId={jobSpecId} runs={latestJobRuns} />
      <TablePagination
        component='div'
        count={jobRunsCount}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[pageSize]}
        page={state.page - 1}
        onChangePage={() => { } /* handler required by component, so make it a no-op */}
        onChangeRowsPerPage={() => { } /* handler required by component, so make it a no-op */}
        ActionsComponent={TableButtonsWithProps}
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

export const ConnectedJobSpecRuns = connect(mapStateToProps, matchRouteAndMapDispatchToProps({ fetchJobSpecRuns }))(
  JobSpecRuns
)

export default withStyles(styles)(ConnectedJobSpecRuns)
