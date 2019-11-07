import React, { useState, useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { RouteComponentProps } from '@reach/router'
import build from 'redux-object'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Title from '../../../components/Title'
import List from '../../../components/Admin/Operators/List'
import { ChangePageEvent } from '../../../components/Table'
import { fetchOperators } from '../../../actions/operators'
import { AppState } from '../../../reducers'
import { ChainlinkNode } from 'explorer/models'

const LOADING_MSG = 'Loading operators...'
const EMPTY_MSG = 'There are no operators added to the Explorer yet.'

const styles = ({ breakpoints, spacing }: Theme) =>
  createStyles({
    container: {
      overflow: 'hidden',
      padding: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        padding: spacing.unit * 3,
      },
    },
  })

interface OwnProps {
  rowsPerPage?: number
}

interface StateProps {
  adminOperators?: ChainlinkNode[]
  count: AppState['adminOperatorsIndex']['count']
}

interface DispatchProps {
  fetchOperators: (page: number, size: number) => void
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const Index: React.FC<Props> = ({
  classes,
  adminOperators,
  fetchOperators,
  count,
  rowsPerPage = 10,
}) => {
  const [currentPage, setCurrentPage] = useState(0)
  const onChangePage = (_event: ChangePageEvent, page: number) => {
    setCurrentPage(page)
    fetchOperators(page + 1, rowsPerPage)
  }

  useEffect(() => {
    fetchOperators(currentPage + 1, rowsPerPage)
  }, [rowsPerPage])

  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      className={classes.container}
    >
      <Grid item xs={12}>
        <Title>Endorsed Operators</Title>

        <List
          currentPage={currentPage}
          operators={adminOperators}
          count={count}
          onChangePage={onChangePage}
          emptyMsg={EMPTY_MSG}
          loadingMsg={LOADING_MSG}
        />
      </Grid>
    </Grid>
  )
}

const operatorsSelector = ({
  adminOperatorsIndex,
  adminOperators,
}: AppState): ChainlinkNode[] | undefined => {
  if (adminOperatorsIndex.items) {
    return adminOperatorsIndex.items.map((id: string) => {
      const document = { adminOperators: adminOperators.items }
      return build(document, 'adminOperators', id)
    })
  }
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    adminOperators: operatorsSelector(state),
    count: state.adminOperatorsIndex.count,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchOperators,
}

export const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index)

export default withStyles(styles)(ConnectedIndex)
