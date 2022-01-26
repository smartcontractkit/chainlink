import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import TablePagination from '@material-ui/core/TablePagination'
import List from '../Jobs/JobRunsList'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import Content from 'components/Content'
import { v2 } from 'src/api'
import { transformPipelineJobRun } from '../Jobs/transformJobRuns'
import { Heading1 } from 'src/components/Heading/Heading1'
import Grid from '@material-ui/core/Grid'

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

const fetchRuns = async (page, size) => {
  const response = await v2.runs.getAllJobRuns({ page, size })

  return response.data.map((run) =>
    transformPipelineJobRun(run.attributes.pipelineSpec.jobID)(run),
  )
}

export const Index = (props) => {
  const { jobSpecId, pageSize, match } = props
  const [page, setPage] = useState(FIRST_PAGE)
  const [latestJobRuns, setLatestJobRuns] = useState([])

  useEffect(() => {
    document.title = 'Job Runs'
  }, [])

  useEffect(() => {
    const queryPage = parseInt(match?.params.jobRunsPage, 10) || FIRST_PAGE
    setPage(queryPage)
    fetchRuns(1, pageSize).then(setLatestJobRuns)
  }, [jobSpecId, pageSize, match])

  const handleChangePage = (_, pageNum) => {
    fetchRuns(pageNum, pageSize).then(setLatestJobRuns)
    setPage(pageNum)
  }

  return (
    <Content>
      <Grid container spacing={32}>
        <Grid item xs={12}>
          <Heading1>Runs</Heading1>
        </Grid>

        <Grid item xs={12}>
          {renderDetails(
            { ...props, latestJobRuns },
            { page },
            handleChangePage,
          )}
        </Grid>
      </Grid>
    </Content>
  )
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
  jobRunsCount: PropTypes.number,
  pageSize: PropTypes.number.isRequired,
  pagePath: PropTypes.string.isRequired,
}

Index.defaultProps = {
  pageSize: 25,
}

export default withStyles(styles)(Index)
