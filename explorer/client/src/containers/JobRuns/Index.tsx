import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import React, { useEffect, useState } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import build from 'redux-object'
import { DispatchBinding } from '@chainlink/ts-helpers'
import { JobRun } from 'explorer/models'
import { fetchJobRuns } from '../../actions/jobRuns'
import List, { Props as ListProps } from '../../components/JobRuns/List'
import { DEFAULT_ROWS_PER_PAGE } from '../../components/Table'
import { AppState } from '../../reducers'
import { searchQuery } from '../../utils/searchQuery'

const styles = ({ spacing, breakpoints }: Theme) =>
  createStyles({
    container: {
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
      paddingLeft: spacing.unit * 2,
      paddingRight: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        paddingTop: spacing.unit * 3,
        paddingBottom: spacing.unit * 3,
        paddingLeft: spacing.unit * 3,
        paddingRight: spacing.unit * 3,
      },
    },
  })

interface OwnProps {
  rowsPerPage?: number
  path: string
}

interface StateProps {
  loading: AppState['jobRuns']['loading']
  error: AppState['jobRuns']['error']
  jobRuns?: JobRun[]
  count?: AppState['jobRunsIndex']['count']
}

interface DispatchProps {
  fetchJobRuns: DispatchBinding<typeof fetchJobRuns>
}

interface Props
  extends WithStyles<typeof styles>,
    OwnProps,
    StateProps,
    DispatchProps {}

const LOADING_MSG = 'Loading search results...'
const EMPTY_MSG =
  "We couldn't find any results for your search query. Try again with the job id, run id, requester, requester id or transaction hash"

const Index = withStyles(styles)(
  ({
    loading,
    error,
    fetchJobRuns,
    rowsPerPage = DEFAULT_ROWS_PER_PAGE,
    classes,
    jobRuns,
    count,
  }: Props) => {
    const [currentPage, setCurrentPage] = useState(0)
    const onChangePage: ListProps['onChangePage'] = (_event, page) => {
      setCurrentPage(page)
    }

    useEffect(() => {
      fetchJobRuns(searchQuery(), currentPage + 1, rowsPerPage)
    }, [fetchJobRuns, currentPage, rowsPerPage])

    return (
      <div className={classes.container}>
        <List
          loading={loading}
          error={error}
          currentPage={currentPage}
          jobRuns={jobRuns}
          rowsPerPage={rowsPerPage}
          count={count}
          onChangePage={onChangePage}
          loadingMsg={LOADING_MSG}
          emptyMsg={EMPTY_MSG}
        />
      </div>
    )
  },
)

function jobRunsSelector({
  jobRunsIndex,
  jobRuns,
  chainlinkNodes,
}: AppState): JobRun[] | undefined {
  return jobRunsIndex.items?.map(id => {
    const document = {
      jobRuns: jobRuns.items,
      chainlinkNodes: chainlinkNodes.items,
    }
    return build(document, 'jobRuns', id)
  })
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    jobRuns: jobRunsSelector(state),
    count: state.jobRunsIndex.count,
    loading: state.jobRuns.loading,
    error: state.jobRuns.error,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchJobRuns,
}

export default connect(mapStateToProps, mapDispatchToProps)(Index)
