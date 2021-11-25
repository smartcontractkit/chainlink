import React from 'react'

import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'

import Paper from '@material-ui/core/Paper'

const styles = (theme: Theme) =>
  createStyles({
    paper: {
      display: 'flex',
      margin: `${theme.spacing.unit * 2.5}px 0`,
      padding: theme.spacing.unit * 3,
    },
    content: {
      flex: 1,
      width: '100%',
    },
    actions: {
      marginTop: -theme.spacing.unit * 1.5,
      marginLeft: -theme.spacing.unit * 4,
      [theme.breakpoints.up('sm')]: {
        marginLeft: 0,
        marginRight: -theme.spacing.unit * 1.5,
      },
    },
  })

interface Props extends WithStyles<typeof styles> {
  actions?: React.ReactNode
}

// DetailsCard provides a box to present details about an object.
//
// This should be used in conjunction with a Grid and additonally provides space
// for actions which should contain a icon button that opens a menu.
//
// We may make this more specific to the use case in the future
export const DetailsCard = withStyles(styles)(
  ({ actions, children, classes }: React.PropsWithChildren<Props>) => {
    return (
      <Paper className={classes.paper}>
        <div className={classes.content}>{children}</div>

        {actions && <div className={classes.actions}>{actions}</div>}
      </Paper>
    )
  },
)
