import React, { useEffect } from 'react'
import { bindActionCreators, Dispatch } from 'redux'
import { connect } from 'react-redux'
import { RouteComponentProps } from '@reach/router'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import Paper from '@material-ui/core/Paper'
import { fetchOperators } from '../../../actions/operators'
import { State } from '../../../reducers'

const styles = ({ breakpoints, spacing }: Theme) =>
  createStyles({
    container: {
      overflow: 'hidden',
      padding: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        padding: spacing.unit * 3,
      },
    },
    paper: {
      padding: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        padding: spacing.unit * 3,
      },
    },
  })

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface OwnProps {}

interface StateProps {
  authenticated: boolean
}

interface DispatchProps {
  fetchOperators: () => any
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const Index = ({ classes, fetchOperators }: Props) => {
  useEffect(() => {
    fetchOperators()
  }, [])

  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      className={classes.container}
    >
      <Grid item xs={12}>
        <Paper className={classes.paper}>
          <Typography variant="h3" gutterBottom>
            Operators
          </Typography>

          <Typography variant="body1">
            Add list of operators when working on story [#166832915]
          </Typography>
        </Paper>
      </Grid>
    </Grid>
  )
}

function mapStateToProps(state: State): StateProps {
  return {
    authenticated: state.adminAuth.allowed,
  }
}

function mapDispatchToProps(dispatch: Dispatch): DispatchProps {
  return bindActionCreators({ fetchOperators }, dispatch)
}

export const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index)

export default withStyles(styles)(ConnectedIndex)
