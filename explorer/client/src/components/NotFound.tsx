import React from 'react'
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

interface Props extends WithStyles<typeof styles>, RouteComponentProps {}

export const NotFound: React.FC<Props> = ({ classes }) => {
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
            404
          </Typography>

          <Typography variant="body1">
            The page you are looking for is not here
          </Typography>
        </Paper>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(NotFound)
