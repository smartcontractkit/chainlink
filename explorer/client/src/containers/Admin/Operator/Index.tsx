import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { RouteComponentProps } from '@reach/router'
import { ChainlinkNode } from 'explorer/models'
import React, { useEffect, useState } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import build from 'redux-object'
import { DispatchBinding } from '@chainlink/ts-helpers'
import { fetchAdminOperators } from '../../../actions/adminOperators'
import List from '../../../components/Admin/Operators/List'
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
  loading: boolean
  error: boolean
  count?: AppState['adminOperatorsIndex']['count']
  adminOperators?: ChainlinkNode[]
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
  loading,
  error,
  adminOperators,
  fetchAdminOperators,
  count,
  rowsPerPage = 10,
}) => {
  const [currentPage, setCurrentPage] = useState(0)

  useEffect(() => {
    fetchAdminOperators(currentPage + 1, rowsPerPage)
  }, [rowsPerPage, currentPage, fetchAdminOperators])

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
          loading={loading}
          error={error}
          currentPage={currentPage}
          operators={adminOperators}
          count={count}
          onChangePage={(_, page) => {
            setCurrentPage(page)
          }}
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
  return
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    adminOperators: adminOperatorsSelector(state),
    count: state.adminOperatorsIndex.count,
    loading: state.adminOperators.loading,
    error: state.adminOperators.error,
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
