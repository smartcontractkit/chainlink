import React from 'react'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Grid from '@material-ui/core/Grid'
import Link from '@material-ui/core/Link';
import { createStyles, Theme, withStyles, WithStyles } from '@material-ui/core/styles'
import ConnectedNodes from './ConnectedNodes'

const styles = (theme: Theme) => createStyles({
  appBar: {
    backgroundColor: theme.palette.common.white,
    zIndex: theme.zIndex.modal + 1
  },
  toolbar: {
    paddingLeft: theme.spacing.unit * 5,
    paddingRight: theme.spacing.unit * 5
  }
})

interface IProps extends WithStyles<typeof styles> {
}

const Header = (props: IProps) => {
  return (
    <AppBar className={props.classes.appBar} color="default" position="absolute">
      <Toolbar className={props.classes.toolbar}>
        <Grid container alignItems="center">
          <Grid item xs={12}>
            <Link href="/">
              LINK Stats
            </Link>
            <ConnectedNodes />
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  )
}

export default withStyles(styles)(Header)
