import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch } from 'redux'
import build from 'redux-object'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import List from '../../components/JobRuns/List'
import { getJobRuns } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { ChangePageEvent } from '../../components/Table'

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
        paddingRight: spacing.unit * 3
      }
    }
  })

interface IProps extends WithStyles<typeof styles> {
  path: string
  rowsPerPage: number
  query?: string
  jobRuns?: IJobRun[]
  count?: number
  getJobRuns: Function
}

const Index = withStyles(styles)((props: IProps) => {
  const [currentPage, setCurrentPage] = useState(0)
  const onChangePage = (_event: ChangePageEvent, page: number) => {
    setCurrentPage(page)
    props.getJobRuns(props.query, page + 1, props.rowsPerPage)
  }

  useEffect(() => {
    props.getJobRuns(props.query, currentPage + 1, props.rowsPerPage)
  }, [props.query])

  return (
    <div className={props.classes.container}>
      <List
        currentPage={currentPage}
        jobRuns={props.jobRuns}
        count={props.count}
        onChangePage={onChangePage}
        emptyMsg={EMPTY_MSG}
      />
    </div>
  )
})

const jobRunsSelector = ({
  jobRunsIndex,
  jobRuns,
  chainlinkNodes
}: IState): IJobRun[] | undefined => {
  if (jobRunsIndex.items) {
    return jobRunsIndex.items.map((id: string) => {
      const document = {
        jobRuns: jobRuns.items,
        chainlinkNodes: chainlinkNodes.items
      }
      return build(document, 'jobRuns', id)
    })
  }
}

const jobRunsCountSelector = (state: IState) => {
  return state.jobRunsIndex.count
}

const mapStateToProps = (state: IState) => ({
  rowsPerPage: 10,
  query: state.search.query,
  jobRuns: jobRunsSelector(state),
  count: jobRunsCountSelector(state)
})

const mapDispatchToProps = (dispatch: Dispatch<any>) =>
  bindActionCreators({ getJobRuns }, dispatch)

const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps
)(Index)

export default ConnectedIndex
