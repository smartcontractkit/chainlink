import React from 'react'
// import React, { useState, useEffect } from 'react'
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
// import { fetchAdminOperators } from '../../../actions/adminOperators'
import { AppState } from '../../../reducers'
// import { DispatchBinding } from '../../../utils/types'

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

interface OwnProps extends RouteComponentProps<{ operatorId: string }> {
  // rowsPerPage?: number
  // operatorId: string
}

interface StateProps {
  operator: ChainlinkNode
  // operator: any
}

interface DispatchProps {
  // fetchAdminOperators: DispatchBinding<typeof fetchAdminOperators>
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const Show: React.FC<Props> = ({
  classes,
  operator,
  //   adminOperators,
  //   fetchAdminOperators,
  //   count,
  // rowsPerPage = 10,
}) => {
  //   const [currentPage, setCurrentPage] = useState(0)
  //   const onChangePage = (_event: ChangePageEvent, page: number) => {
  //     setCurrentPage(page)
  //     fetchAdminOperators(page + 1, rowsPerPage)
  //   }
  //   useEffect(() => {
  //     fetchAdminOperators(currentPage + 1, rowsPerPage)
  //   }, [rowsPerPage])
  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      // className={classes.container}
    >
      <Grid item xs={12}>
        <Title>Operator Show</Title>
        <Operator
          // currentPage={currentPage}
          operator={operator}
          // count={count}
          // onChangePage={onChangePage}
          // emptyMsg={EMPTY_MSG}
          // loadingMsg={LOADING_MSG}
        />
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
  //   fetchAdminOperators,
}

export const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Show)

export default withStyles(styles)(ConnectedShow)

// const Show = () => {
//   return (
//     <>
//       <h1>Operator Show</h1>
//     </>
//   )
// }

// export default Show
