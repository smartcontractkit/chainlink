import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'

const styles = (_theme: Theme) => createStyles({})

interface Props extends WithStyles<typeof styles> {
  path: string
}

export const Admin = (_props: Props) => {
  return (
    <div>
      <h1>Admin</h1>
      <div>SignIn</div>
    </div>
  )
}

export default withStyles(styles)(Admin)
