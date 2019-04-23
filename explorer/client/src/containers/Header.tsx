import React from 'react'
import { connect } from 'react-redux'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import ReactResizeDetector from 'react-resize-detector'
import Logo from '../components/Logo'
import SearchForm from '../components/SearchForm'
import SearchBox from '../components/SearchBox'
import { IState } from '../reducers'

const styles = (theme: Theme) =>
  createStyles({
    appBar: {
      backgroundColor: theme.palette.common.white,
      zIndex: theme.zIndex.modal + 1
    },
    toolbar: {
      padding: theme.spacing.unit * 5,
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2
    },
    logoAndSearch: {
      display: 'flex',
      alignItems: 'center'
    },
    logo: {
      width: 200
    },
    searchForm: {
      flexGrow: 1
    },
    search: {
      marginLeft: theme.spacing.unit * 2
    },
    connectedNodes: {
      textAlign: 'right'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  onResize: (width: number, height: number) => void
}

const Header = ({ classes, onResize }: IProps) => {
  return (
    <AppBar className={classes.appBar} color="default">
      <ReactResizeDetector
        refreshMode="debounce"
        refreshRate={200}
        handleWidth
        onResize={onResize}>
        <Toolbar className={classes.toolbar}>
          <Grid container alignItems="center">
            <Grid item xs={8}>
              <div className={classes.logoAndSearch}>
                <Logo className={classes.logo} />
                <SearchForm className={classes.searchForm}>
                  <SearchBox className={classes.search} />
                </SearchForm>
              </div>
            </Grid>
          </Grid>
        </Toolbar>
      </ReactResizeDetector>
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
