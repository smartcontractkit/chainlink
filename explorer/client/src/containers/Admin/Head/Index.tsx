import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { RouteComponentProps } from '@reach/router'
import { Head } from 'explorer/models'
import React, { useEffect, useState } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import build from 'redux-object'
import { DispatchBinding } from '@chainlink/ts-helpers'
import { fetchAdminHeads } from '../../../actions/adminHeads'
import List from '../../../components/Admin/Heads/List'
import { ChangePageEvent } from '../../../components/Table'
import Title from '../../../components/Title'
import { AppState } from '../../../reducers'

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
  adminHeads?: Head[]
  count: AppState['adminHeadsIndex']['count']
}

interface DispatchProps {
  fetchAdminHeads: DispatchBinding<typeof fetchAdminHeads>
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const Index: React.FC<Props> = ({
  classes,
  adminHeads,
  fetchAdminHeads,
  count,
  rowsPerPage = 10,
}) => {
  const [currentPage, setCurrentPage] = useState(1)
  const onChangePage = (_event: ChangePageEvent, page: number) => {
    setCurrentPage(page)
  }

  useEffect(() => {
    fetchAdminHeads(currentPage, rowsPerPage)
  }, [rowsPerPage, currentPage, fetchAdminHeads])

  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      className={classes.container}
    >
      <Grid item xs={12}>
        <Title>Heads</Title>

        <List
          currentPage={currentPage}
          heads={adminHeads}
          count={count}
          onChangePage={onChangePage}
        />
      </Grid>
    </Grid>
  )
}

const adminHeadsSelector = ({
  adminHeadsIndex,
  adminHeads,
}: AppState): Head[] | undefined => {
  if (adminHeadsIndex.items) {
    return adminHeadsIndex.items.map(id =>
      build({ adminHeads: adminHeads.items }, 'adminHeads', id),
    )
  }
  return
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    adminHeads: adminHeadsSelector(state),
    count: state.adminHeadsIndex.count,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchAdminHeads,
}

export const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index)

export default withStyles(styles)(ConnectedIndex)
