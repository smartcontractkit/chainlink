import React, { useState, useEffect } from 'react'
import Paper from '@material-ui/core/Paper'
import IconButton from '@material-ui/core/IconButton'
import SearchIcon from '@material-ui/icons/Search'
import InputBase from '@material-ui/core/InputBase'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import classNames from 'classnames'
import { searchQuery } from '../utils/searchQuery'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    paper: {
      display: 'flex',
      border: 'solid 1px',
      borderColor: palette.grey['300'],
      padding: spacing.unit,
    },
    query: {
      flexGrow: 1,
      boxSizing: 'border-box',
      color: palette.text.primary,
    },
  })

interface Props extends WithStyles<typeof styles> {
  className?: string
  query?: string
}

const SearchBox = ({ classes, className }: Props) => {
  const [query, setQuery] = useState<string>(searchQuery())

  useEffect(() => {
    const onPopState = () => setQuery(searchQuery())
    window.addEventListener('popstate', onPopState)
    return () => window.removeEventListener('popstate', onPopState)
  }, [setQuery])

  return (
    <Paper elevation={0} className={classNames(classes.paper, className)}>
      <IconButton aria-label="Search" type="submit">
        <SearchIcon />
      </IconButton>
      <InputBase
        className={classes.query}
        value={query}
        onChange={e => setQuery(e.target.value)}
        placeholder="Search for Job IDs, Run IDs, Transaction Hashes, or Requesting Addresses"
        name="search"
      />
    </Paper>
  )
}

export default withStyles(styles)(SearchBox)
