import React from 'react'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Grid from '@material-ui/core/Grid'
import Link from '@material-ui/core/Link';
import { withStyles } from '@material-ui/core/styles'

const styles = (theme: any) => {
  return {
    appBar: {
      backgroundColor: theme.palette.common.white,
      zIndex: theme.zIndex.modal + 1
    },
    toolbar: {
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5
    }
  }
}

const Header = (props: any) => {
  return (
    <AppBar className={props.classes.appBar} color="default" position="absolute">
      <Toolbar className={props.classes.toolbar}>
        <Grid container alignItems="center">
          <Grid item xs={12}>
            <Link href="/">
              LINK Stats
            </Link>
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  )
}

export default withStyles(styles)(Header)
