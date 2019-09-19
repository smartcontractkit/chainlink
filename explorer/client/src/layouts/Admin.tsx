import React, { useState } from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Header from '../components/Header'
import Logo from '../components/Admin/Logo'
import { DEFAULT_HEADER_HEIGHT } from '../constants'

const styles = (theme: Theme) =>
  createStyles({
    logo: {
      marginRight: theme.spacing.unit * 2,
      width: 200,
    },
  })

interface Props extends WithStyles<typeof styles> {
  children?: any
  path: string
}

export const Admin = ({ children, classes }: Props) => {
  const [height, setHeight] = useState<number>(DEFAULT_HEADER_HEIGHT)
  const onHeaderResize = (_width: number, height: number) => setHeight(height)

  return (
    <Grid container spacing={24}>
      <Grid item xs={12}>
        <Header onResize={onHeaderResize}>
          <Logo className={classes.logo} />
        </Header>

        <main style={{ paddingTop: height }}>{children}</main>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(Admin)
