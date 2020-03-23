import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { RouteComponentProps } from '@reach/router'
import React, { useState } from 'react'
import Footer from '../components/Admin/Operators/Footer'
import { DEFAULT_HEADER_HEIGHT } from '../constants'
import Header from '../containers/Admin/Header'
import AdminPrivate from '../containers/Admin/Private'

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

interface Props extends RouteComponentProps, WithStyles<typeof styles> {}

export const Admin: React.FC<Props> = ({ children, classes }) => {
  const [height, setHeight] = useState(DEFAULT_HEADER_HEIGHT)
  const onHeaderResize = (_width: number, height: number) => setHeight(height)

  return (
    <AdminPrivate>
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
      <Footer />
    </AdminPrivate>
  )
}

export default withStyles(styles)(Admin)
