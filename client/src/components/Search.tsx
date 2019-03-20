import React from 'react'
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton'
import SearchIcon from '@material-ui/icons/Search'
import InputBase from '@material-ui/core/InputBase'
import { createStyles, Theme, withStyles, WithStyles } from '@material-ui/core/styles'
import classNames from 'classnames'

interface IProps extends WithStyles<typeof styles> {
  className: string,
  searchParams: URLSearchParams
}

const styles = (theme: Theme) => createStyles({
  paper: {
    border: 'solid 1px',
    borderColor: theme.palette.primary.light
  }
})

const SEARCH_PARAM: string = 'search'

const search = (searchParams: URLSearchParams): string | undefined => {
  const search = searchParams.get(SEARCH_PARAM)
  if (search) {
    return search
  }
}

const Search = (props: IProps) => {
  return (
    <Paper elevation={0} className={classNames(props.classes.paper, props.className)}>
      <form method="GET">
        <IconButton aria-label="Search">
          <SearchIcon />
        </IconButton>
        <InputBase
          value={search(props.searchParams)}
          placeholder="Search for something"
          name="search"
        />
      </form>
    </Paper>
  )
}

Search.defaultProps = {
  searchParams: (new URL(document.location.toString())).searchParams
}

export default withStyles(styles)(Search)
