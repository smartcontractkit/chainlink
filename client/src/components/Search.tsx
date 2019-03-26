import React from 'react'
import { connect } from 'react-redux'
import Paper from '@material-ui/core/Paper'
import IconButton from '@material-ui/core/IconButton'
import SearchIcon from '@material-ui/icons/Search'
import InputBase from '@material-ui/core/InputBase'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import classNames from 'classnames'
import { IState } from '../reducers'

const styles = (theme: Theme) =>
  createStyles({
    form: {
      display: 'flex'
    },
    paper: {
      border: 'solid 1px',
      borderColor: theme.palette.primary.light
    },
    query: {
      boxSizing: 'border-box',
      flexGrow: 1
    }
  })

interface IProps extends WithStyles<typeof styles> {
  className: string
  query?: string
}

const Search = ({ classes, className, query }: IProps) => {
  return (
    <Paper elevation={0} className={classNames(classes.paper, className)}>
      <form method="GET" className={classes.form}>
        <IconButton aria-label="Search" type="submit">
          <SearchIcon />
        </IconButton>
        <InputBase
          className={classes.query}
          defaultValue={query}
          placeholder="Search for something"
          name="search"
        />
      </form>
    </Paper>
  )
}

const mapStateToProps = (state: IState) => {
  return {
    query: state.search.query
  }
}

const mapDispatchToProps = () => ({})

const ConnectedSearch = connect(
  mapStateToProps,
  mapDispatchToProps
)(Search)

export default withStyles(styles)(ConnectedSearch)
