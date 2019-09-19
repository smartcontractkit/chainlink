import React from 'react'
import { RouteComponentProps } from '@reach/router'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'

const styles = () =>
  createStyles({
    container: {
      overflow: 'hidden',
    },
  })

interface Props extends RouteComponentProps, WithStyles<typeof styles> {
  children?: any
}

export const AdminMinimal = ({ children, classes }: Props) => {
  return (
    <Grid
      container
      spacing={24}
      alignItems="center"
      className={classes.container}
    >
      <Grid item xs={12}>
        <main>{children}</main>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(AdminMinimal)
