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
import Operator from '../../../components/Admin/Operators/Operator'
// import { ChangePageEvent } from '../../../components/Table'
import { fetchAdminOperator } from '../../../actions/adminOperators'
import { AppState } from '../../../reducers'
import { DispatchBinding } from '../../../utils/types'

// const LOADING_MSG = 'Loading operators...'
// const EMPTY_MSG = 'There are no operators added to the Explorer yet.'

const styles = ({ breakpoints, spacing }: Theme) =>
  createStyles({
    // container: {
    //   overflow: 'hidden',
    //   padding: spacing.unit * 2,
    //   [breakpoints.up('sm')]: {
    //     padding: spacing.unit * 3,
    //   },
    // },
  })

interface OwnProps extends RouteComponentProps<{ operatorId: string }> {}

interface StateProps {
  operator: ChainlinkNode
}

interface DispatchProps {
  fetchAdminOperator: DispatchBinding<typeof fetchAdminOperator>
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

const Show: React.FC<Props> = ({
  classes,
  operator,
  operatorId,
  fetchAdminOperator,
}) => {
  useEffect(() => {
    if (operatorId) {
      fetchAdminOperator(operatorId)
    }
  }, [operatorId])

  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      // className={classes.container}
    >
      <Grid item xs={12}>
        <Title>Operator Details</Title>
        <Operator operator={operator} />
      </Grid>
    </Grid>
  )
}

const adminOperatorSelector = (
  operatorId: string | undefined,
  state: AppState,
): ChainlinkNode => {
  return build(state.adminOperators, 'items', operatorId)
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
  ownProps,
) => {
  return {
    operator: adminOperatorSelector(ownProps.operatorId, state),
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchAdminOperator,
}

const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Show)
const StyledShow = withStyles(styles)(ConnectedShow)

export default StyledShow
