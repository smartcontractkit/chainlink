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
import { IconButton } from '@material-ui/core'
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft'
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight'
import FirstPageIcon from '@material-ui/icons/FirstPage'
import LastPageIcon from '@material-ui/icons/LastPage'

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
    marginLeft: theme.spacing.unit * 2.5
  }
})

const TableButtons = props => {
  const lastPage = Math.ceil(props.count / props.rowsPerPage)
  const firstPage = 1
  const currentPage = props.page
  const handlePage = page => {
    page = Math.min(page, lastPage)
    page = Math.max(page, firstPage)
    const curry = e => {
      if(props.history)
        props.history.replace(`/job_specs/${props.specID}/runs/page/${page}`)
      props.onChangePage(e, page)
    }
    return curry
  }
  return (
    <div className={props.classes.customButtons}>
      <IconButton onClick={handlePage(firstPage)} disabled={currentPage === firstPage} aria-label='First Page'>
        <FirstPageIcon />
      </IconButton>
      <IconButton onClick={handlePage(currentPage-1)} disabled={currentPage === firstPage} aria-label='Previous Page'>
        <KeyboardArrowLeft />
      </IconButton>
      <IconButton onClick={handlePage(currentPage+1)} disabled={currentPage >= lastPage} aria-label='Next Page'>
        <KeyboardArrowRight />
      </IconButton>
      <IconButton onClick={handlePage(lastPage)} disabled={currentPage >= lastPage} aria-label='Last Page'>
        <LastPageIcon />
      </IconButton>
    </div>
  )
}

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
    const firstPage = 1
    if (this.props.match.params.jobRunsPage) {
      const START_PAGE = this.props.match.params.jobRunsPage
      this.setState({ page: START_PAGE })
      fetchJobSpecRuns(jobSpecId, START_PAGE, pageSize)
    } else {
      this.setState({ page: firstPage })
      fetchJobSpecRuns(jobSpecId, firstPage, pageSize)
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
        page={state.page-1}
        onChangePage={() => {} /* handler required by component, so make it a no-op */}
        onChangeRowsPerPage={() => {} /* handler required by component, so make it a no-op */}
        ActionsComponent={withStyles(styles)(TableButtonsWithProps)}
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
