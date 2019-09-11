import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import React, { useEffect, useState } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { bindActionCreators } from 'redux'
import build from 'redux-object'
import { getJobRuns } from '../../actions/jobRuns'
import List from '../../components/JobRuns/List'
import { ChangePageEvent } from '../../components/Table'
import { State } from '../../reducers'
import { Query } from '../../reducers/search'

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
  query: State['search']['query']
  jobRuns?: JobRun[]
  count: State['jobRunsIndex']['count']
}

interface DispatchProps {
  getJobRuns: (query: Query, page: number, size: number) => void
}

interface Props
  extends WithStyles<typeof styles>,
    OwnProps,
    StateProps,
    DispatchProps {}

const Index = withStyles(styles)(
  ({ getJobRuns, query, rowsPerPage = 10, classes, jobRuns, count }: Props) => {
    const [currentPage, setCurrentPage] = useState(0)
    const onChangePage = (_event: ChangePageEvent, page: number) => {
      setCurrentPage(page)
      getJobRuns(query, page + 1, rowsPerPage)
    }

    useEffect(() => {
      getJobRuns(query, currentPage + 1, rowsPerPage)
    }, [getJobRuns, query, currentPage, rowsPerPage])

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
}: State): JobRun[] | undefined => {
  if (jobRunsIndex.items) {
    return jobRunsIndex.items.map((id: string) => {
      const document = {
        jobRuns: jobRuns.items,
        chainlinkNodes: chainlinkNodes.items,
      }
      return build(document, 'jobRuns', id)
    })
  }
}

const mapDispatchToProps: MapDispatchToProps<
  DispatchProps,
  OwnProps
> = dispatch => bindActionCreators({ getJobRuns }, dispatch)

const mapStateToProps: MapStateToProps<StateProps, OwnProps, State> = state => {
  return {
    query: state.search.query,
    jobRuns: jobRunsSelector(state),
    count: state.jobRunsIndex.count,
  }
}

const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index)

export default ConnectedIndex
