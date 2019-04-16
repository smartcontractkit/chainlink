import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import TablePagination from '@material-ui/core/TablePagination'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRuns } from 'actions'
import jobRunsSelector from 'selectors/jobRuns'
import jobRunsCountSelector from 'selectors/jobRunsCount'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import List from 'components/JobRuns/List'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import Title from 'components/Title'
import Content from 'components/Content'
import { useHooks, useEffect, useState } from 'use-react-hooks'

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const renderLatestRuns = (props, state, handleChangePage) => {
  const { jobSpecId, latestJobRuns, jobRunsCount, pageSize } = props
  const TableButtonsWithProps = () => (
    <TableButtons
      history={props.history}
      count={jobRunsCount}
      onChangePage={handleChangePage}
      page={state.page}
      specID={jobSpecId}
      rowsPerPage={pageSize}
      replaceWith={`/jobs/${jobSpecId}/runs/page`}
    />
  )
  return (
    <Card>
      <List
        jobSpecId={jobSpecId}
        runs={latestJobRuns}
        showJobRunsCount={jobRunsCount}
      />
      <TablePagination
        component="div"
        count={jobRunsCount}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[pageSize]}
        page={state.page - 1}
        onChangePage={
          () => {} /* handler required by component, so make it a no-op */
        }
        onChangeRowsPerPage={
          () => {} /* handler required by component, so make it a no-op */
        }
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

export const Index = useHooks(props => {
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    document.title = 'Job Runs'
    const queryPage = props.match
      ? parseInt(props.match.params.jobRunsPage, 10) || FIRST_PAGE
      : FIRST_PAGE
    setPage(queryPage)
    fetchJobRuns(jobSpecId, queryPage, pageSize)
  }, [])
  const { classes, jobSpecId, fetchJobRuns, pageSize } = props
  const handleChangePage = (e, pageNum) => {
    fetchJobRuns(jobSpecId, pageNum, pageSize)
    setPage(pageNum)
  }

  return (
    <Content>
      <Breadcrumb className={classes.breadcrumb}>
        <BreadcrumbItem href="/">Dashboard</BreadcrumbItem>
        <BreadcrumbItem>&gt;</BreadcrumbItem>
        <BreadcrumbItem href={`/jobs/${jobSpecId}`}>
          Job ID: {jobSpecId}
        </BreadcrumbItem>
        <BreadcrumbItem>&gt;</BreadcrumbItem>
        <BreadcrumbItem>Runs</BreadcrumbItem>
      </Breadcrumb>
      <Title>Runs</Title>
      {renderDetails(props, { page }, handleChangePage)}
    </Content>
  )
})

Index.propTypes = {
  classes: PropTypes.object.isRequired,
  latestJobRuns: PropTypes.array,
  jobRunsCount: PropTypes.number,
  pageSize: PropTypes.number.isRequired
}

Index.defaultProps = {
  latestJobRuns: [],
  pageSize: 10
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const jobRunsCount = jobRunsCountSelector(state)
  const latestJobRuns = jobRunsSelector(state, jobSpecId)

  return {
    jobSpecId,
    latestJobRuns,
    jobRunsCount
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRuns })
)(Index)

export default withStyles(styles)(ConnectedIndex)
