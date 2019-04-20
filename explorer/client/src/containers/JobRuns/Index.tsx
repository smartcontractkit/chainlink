import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch, Action } from 'redux'
import { denormalize } from 'normalizr'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import List from '../../components/JobRuns/List'
import { getJobRuns } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { JobRun } from '../../entities'

const styles = ({ spacing }: Theme) =>
  createStyles({
    container: {
      margin: spacing.unit * 5
    }
  })

interface IProps extends WithStyles<typeof styles> {
  query?: string
  jobRuns?: IJobRun[]
  getJobRuns: Function
  path: string
  page: number
  size: number
}

const Index = withStyles(styles)(
  ({ query, page, size, jobRuns, getJobRuns, classes }: IProps) => {
    useEffect(() => {
      getJobRuns(query, page, size)
    }, [])

    return (
      <div className={classes.container}>
        <List jobRuns={jobRuns} />
      </div>
    )
  }
)

const jobRunsSelector = ({
  jobRunsIndex,
  jobRuns
}: IState): IJobRun[] | undefined => {
  if (jobRunsIndex.items) {
    return denormalize(jobRunsIndex.items, [JobRun], { jobRuns: jobRuns.items })
  }
}

const mapStateToProps = (state: IState) => ({
  query: state.search.query,
  jobRuns: jobRunsSelector(state),
  page: 1,
  size: 300
})

const mapDispatchToProps = (dispatch: Dispatch<Action<any>>) =>
  bindActionCreators({ getJobRuns }, dispatch)

const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps
)(Index)

export default ConnectedIndex
