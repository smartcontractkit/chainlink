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
import { ChainlinkNode } from 'explorer/models'
import Title from '../../../components/Title'
import List from '../../../components/Admin/Operators/List'
import { ChangePageEvent } from '../../../components/Table'
import { fetchAdminOperators } from '../../../actions/adminOperators'
import { AppState } from '../../../reducers'
import { DispatchBinding } from '../../../utils/types'

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
  fetchAdminOperators: DispatchBinding<typeof fetchAdminOperators>
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
  fetchAdminOperators,
  count,
  rowsPerPage = 10,
}) => {
  const [currentPage, setCurrentPage] = useState(1)
  const onChangePage = (_event: ChangePageEvent, page: number) => {
    setCurrentPage(page)
  }

  useEffect(() => {
    fetchAdminOperators(currentPage, rowsPerPage)
  }, [rowsPerPage, currentPage])

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
        />
      </Grid>
    </Grid>
  )
}

const adminOperatorsSelector = ({
  adminOperatorsIndex,
  adminOperators,
}: AppState): ChainlinkNode[] | undefined => {
  if (adminOperatorsIndex.items) {
    return adminOperatorsIndex.items.map(id =>
      build({ adminOperators: adminOperators.items }, 'adminOperators', id),
    )
  }
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    adminOperators: adminOperatorsSelector(state),
    count: state.adminOperatorsIndex.count,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchAdminOperators,
}

export const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index)

export default withStyles(styles)(ConnectedIndex)
