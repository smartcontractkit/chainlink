import React, { useState } from 'react'
import { RouteComponentProps } from '@reach/router'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import AdminPrivate from '../containers/Admin/Private'
import Header from '../containers/Admin/Header'
import { DEFAULT_HEADER_HEIGHT } from '../constants'

const styles = (theme: Theme) =>
  createStyles({
    avatar: {
      float: 'right',
    },
    container: {
      overflowX: 'hidden',
    },
    logo: {
      marginRight: theme.spacing.unit * 2,
      width: 200,
    },
  })

interface Props extends RouteComponentProps, WithStyles<typeof styles> {
  children?: any
}

export const Admin = ({ children, classes }: Props) => {
  const [height, setHeight] = useState<number>(DEFAULT_HEADER_HEIGHT)
  const onHeaderResize = (_width: number, height: number) => setHeight(height)

  return (
    <>
      <AdminPrivate />

      <Header onHeaderResize={onHeaderResize} />

      <Grid container spacing={24} className={classes.container}>
        <Grid item xs={12}>
          <Grid container>
            <Grid item xs={12}>
              <main style={{ paddingTop: height }}>{children}</main>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </>
  )
}

export default withStyles(styles)(Admin)
