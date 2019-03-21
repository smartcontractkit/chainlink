import React from 'react'
import { connect } from 'react-redux'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Grid from '@material-ui/core/Grid'
import { createStyles, Theme, withStyles, WithStyles } from '@material-ui/core/styles'
import Logo from '../components/Logo'
import Search from '../components/Search'
import ConnectedNodes from '../components/ConnectedNodes'
import { IState } from '../reducers'

const styles = (theme: Theme) => createStyles({
  appBar: {
    backgroundColor: theme.palette.common.white,
    zIndex: theme.zIndex.modal + 1
  },
  toolbar: {
    paddingLeft: theme.spacing.unit * 5,
    paddingRight: theme.spacing.unit * 5
  },
  logoAndSearch: {
    display: 'flex',
    alignItems: 'center'
  },
  logo: {
    width: 150
  },
  search: {
    flexGrow: 1
  },
  connectedNodes: {
    textAlign: 'right'
  }
})

interface IProps extends WithStyles<typeof styles> {
}

const Header = (props: IProps) => {
  return (
    <AppBar className={props.classes.appBar} color="default" position="absolute">
      <Toolbar className={props.classes.toolbar}>
        <Grid container alignItems="center">
          <Grid item xs={8}>
            <div className={props.classes.logoAndSearch}>
              <Logo className={props.classes.logo} />
              <Search className={props.classes.search} />
            </div>
          </Grid>

          <Grid item xs={4}>
            <ConnectedNodes className={props.classes.connectedNodes} />
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  )
}

const mapStateToProps = (state: IState) => {
  return {
    search: state.search.query
  }
}

const mapDispatchToProps = () => ({})

const ConnectedHeader = connect(
  mapStateToProps,
  mapDispatchToProps
)(Header)

export default withStyles(styles)(ConnectedHeader)
