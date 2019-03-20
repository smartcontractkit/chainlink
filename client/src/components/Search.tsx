import React from 'react'
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton'
import SearchIcon from '@material-ui/icons/Search'
import InputBase from '@material-ui/core/InputBase'
import { createStyles, Theme, withStyles, WithStyles } from '@material-ui/core/styles'
import classNames from 'classnames'

const styles = (theme: Theme) => createStyles({
  paper: {
    border: 'solid 1px',
    borderColor: theme.palette.primary.light
  }
})

interface IProps extends WithStyles<typeof styles> {
  className: string
}

const Search = (props: IProps) => {
  return (
    <Paper elevation={0} className={classNames(props.classes.paper, props.className)}>
      <IconButton aria-label="Search">
        <SearchIcon />
      </IconButton>
      <InputBase placeholder="Search for something" />
    </Paper>
  )
}

export default withStyles(styles)(Search)
