import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRuns } from 'actions'
import jobRunsSelector from 'selectors/jobRuns'
import jobRunsCountSelector from 'selectors/jobRunsCount'
import { FIRST_PAGE, GenericList } from 'components/GenericList'
import Title from 'components/Title'
import { useHooks, useEffect, useState } from 'use-react-hooks'
import Content from 'components/Content'

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5,
  },
})

const buildItems = runs =>
  runs.map(r => [
    { type: 'link', text: r.id, to: `/jobs/${r.jobId}/runs/id/${r.id}` },
    { type: 'time_ago', text: r.createdAt },
    { type: 'status', text: r.status },
  ])

export const Index = useHooks(props => {
  const { jobSpecId, fetchJobRuns, pageSize } = props
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    document.title = 'Job Runs'
    const queryPage = props.match
      ? parseInt(props.match.params.jobRunsPage, 10) || FIRST_PAGE
      : FIRST_PAGE
    setPage(queryPage - 1)
    fetchJobRuns({ jobSpecId, page: queryPage, size: pageSize })
  }, [])
  const handleChangePage = (e, pageNum) => {
    if (e) {
      setPage(pageNum)
      fetchJobRuns({ jobSpecId, page: pageNum + 1, size: pageSize })
      if (props.history)
        props.history.push(`/jobs/${jobSpecId}/runs/page/${pageNum + 1}`)
    }
  }
  return (
    <Content>
      <Title>Runs</Title>
      {props.latestJobRuns ? (
        <GenericList
          emptyMsg="No jobs have been run yet"
          headers={['Id', 'Created', 'Status']}
          items={buildItems(props.latestJobRuns)}
          onChangePage={handleChangePage}
          count={props.jobRunsCount}
          currentPage={page}
          rowsPerPage={pageSize}
        />
      ) : (
        <div>Fetching...</div>
      )}
    </Content>
  )
})

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

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const jobRunsCount = jobRunsCountSelector(state)
  const latestJobRuns = jobRunsSelector(state, jobSpecId)

  return {
    jobSpecId,
    latestJobRuns,
    jobRunsCount,
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRuns }),
)(Index)

export default withStyles(styles)(ConnectedIndex)
