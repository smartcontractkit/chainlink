import React from 'react'

import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import WarningIcon from '@material-ui/icons//Warning'

const styles = (theme: Theme) =>
  createStyles({
    paper: {
      backgroundColor: theme.palette.primary.light,
      padding: theme.spacing.unit * 2,
      display: 'flex',
    },
    icon: {
      marginRight: theme.spacing.unit,
    },
  })

interface Props extends WithStyles<typeof styles> {}

export const BetaAlert = withStyles(styles)(({ classes }: Props) => {
  return (
    <Paper className={classes.paper}>
      <WarningIcon className={classes.icon} />
      <Typography>Multi-chain functionality is in Beta</Typography>
    </Paper>
  )
})
