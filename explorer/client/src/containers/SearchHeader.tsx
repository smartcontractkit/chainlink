import React from 'react'
import { connect } from 'react-redux'
import Hidden from '@material-ui/core/Hidden'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Header from '../components/Header'
import Logo from '../components/Logo'
import SearchForm from '../components/SearchForm'
import SearchBox from '../components/SearchBox'
import { State } from '../reducers'

const STACKED_LOGO_HEIGHT = 40

const styles = (theme: Theme) =>
  createStyles({
    logoAndSearch: {
      display: 'flex',
      alignItems: 'center',
    },
    logo: {
      marginRight: theme.spacing.unit * 2,
      width: 200,
    },
    stackedLogo: {
      display: 'block',
    },
    searchForm: {
      flexGrow: 1,
    },
    connectedNodes: {
      textAlign: 'right',
    },
  })

interface Props extends WithStyles<typeof styles> {
  onResize: (width: number, height: number) => void
}

const SearchHeader = ({ classes, onResize }: Props) => {
  return (
    <Header onResize={onResize}>
      <Hidden xsDown>
        <Grid container alignItems="center">
          <Grid item sm={12} md={10} lg={9}>
            <div className={classes.logoAndSearch}>
              <Logo className={classes.logo} />
              <SearchForm className={classes.searchForm}>
                <SearchBox />
              </SearchForm>
            </div>
          </Grid>
        </Grid>
      </Hidden>
      <Hidden smUp>
        <Grid container alignItems="center" spacing={0}>
          <Grid item xs={12}>
            <Logo
              className={classes.stackedLogo}
              height={STACKED_LOGO_HEIGHT}
            />
          </Grid>
          <Grid item xs={12}>
            <SearchForm className={classes.searchForm}>
              <SearchBox />
            </SearchForm>
          </Grid>
        </Grid>
      </Hidden>
    </Header>
  )
}

const mapStateToProps = (state: State) => {
  return {
    search: state.search.query,
  }
}

const mapDispatchToProps = () => ({})

const ConnectedSearchHeader = connect(
  mapStateToProps,
  mapDispatchToProps,
)(SearchHeader)

export default withStyles(styles)(ConnectedSearchHeader)
