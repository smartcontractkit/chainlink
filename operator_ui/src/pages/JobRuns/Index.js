import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import TablePagination from '@material-ui/core/TablePagination'
import List from '../Jobs/JobRunsList'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import Title from 'components/Title'
import Content from 'components/Content'

const styles = (theme) => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5,
  },
})

const renderLatestRuns = (props, state, handleChangePage) => {
  const { jobSpecId, latestJobRuns, jobRunsCount = 0, pageSize } = props
  const pagePath = props.pagePath.replace(':jobSpecId', jobSpecId)

  const TableButtonsWithProps = () => (
    <TableButtons
      history={props.history}
      count={jobRunsCount}
      onChangePage={handleChangePage}
      page={state.page}
      specID={jobSpecId}
      rowsPerPage={pageSize}
      replaceWith={pagePath}
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

const Fetching = () => <div>Fetching...</div>

const renderDetails = (props, state, handleChangePage) => {
  if (props.latestJobRuns) {
    return renderLatestRuns(props, state, handleChangePage)
  } else {
    return <Fetching />
  }
}

export const Index = (props) => {
  const { jobSpecId, pageSize, match } = props
  const [page, setPage] = useState(FIRST_PAGE)

  useEffect(() => {
    document.title = 'Job Runs'
  }, [])

  useEffect(() => {
    const queryPage = parseInt(match?.params.jobRunsPage, 10) || FIRST_PAGE
    setPage(queryPage)
    // TODO - Fetch V2 jobs runs here
  }, [jobSpecId, pageSize, match])
  const handleChangePage = (_, pageNum) => {
    // TODO - Fetch V2 jobs runs here
    setPage(pageNum)
  }

  return (
    <Content>
      <Title>Runs</Title>

      {renderDetails(props, { page }, handleChangePage)}
    </Content>
  )
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
  latestJobRuns: PropTypes.array,
  jobRunsCount: PropTypes.number,
  pageSize: PropTypes.number.isRequired,
  pagePath: PropTypes.string.isRequired,
}

Index.defaultProps = {
  latestJobRuns: [],
  pageSize: 25,
}

export default withStyles(styles)(Index)
