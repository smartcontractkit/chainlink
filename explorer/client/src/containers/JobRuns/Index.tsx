import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { JobRun } from 'explorer/models'
import React, { useEffect, useState } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import build from 'redux-object'
import { fetchJobRuns } from '../../actions/jobRuns'
import List from '../../components/JobRuns/List'
import { ChangePageEvent } from '../../components/Table'
import { AppState } from '../../reducers'
import { DispatchBinding } from '../../utils/types'

const EMPTY_MSG =
  "We couldn't find any results for your search query. Try again with the job id, run id, requester, requester id or transaction hash"

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
  query: AppState['search']['query']
  jobRuns?: JobRun[]
  count: AppState['jobRunsIndex']['count']
}

interface DispatchProps {
  fetchJobRuns: DispatchBinding<typeof fetchJobRuns>
}

interface Props
  extends WithStyles<typeof styles>,
    OwnProps,
    StateProps,
    DispatchProps {}

const Index = withStyles(styles)(
  ({
    fetchJobRuns,
    query,
    rowsPerPage = 10,
    classes,
    jobRuns,
    count,
  }: Props) => {
    const [currentPage, setCurrentPage] = useState(0)
    const onChangePage = (_event: ChangePageEvent, page: number) => {
      setCurrentPage(page)
      fetchJobRuns(query, page + 1, rowsPerPage)
    }

    useEffect(() => {
      fetchJobRuns(query, currentPage + 1, rowsPerPage)
    }, [fetchJobRuns, query, currentPage, rowsPerPage])

    return (
      <div className={classes.container}>
        <List
          currentPage={currentPage}
          jobRuns={jobRuns}
          count={count}
          onChangePage={onChangePage}
          emptyMsg={EMPTY_MSG}
        />
      </div>
    )
  },
)

const jobRunsSelector = ({
  jobRunsIndex,
  jobRuns,
  chainlinkNodes,
}: AppState): JobRun[] | undefined => {
  if (jobRunsIndex.items) {
    return jobRunsIndex.items.map((id: string) => {
      const document = {
        jobRuns: jobRuns.items,
        chainlinkNodes: chainlinkNodes.items,
      }
      return build(document, 'jobRuns', id)
    })
  }

  return
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    query: state.search.query,
    jobRuns: jobRunsSelector(state),
    count: state.jobRunsIndex.count,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchJobRuns,
}

const ConnectedIndex = connect(mapStateToProps, mapDispatchToProps)(Index)

export default ConnectedIndex
