import React from 'react'
import IconButton from '@material-ui/core/IconButton'
import FirstPageIcon from '@material-ui/icons/FirstPage'
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft'
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight'
import LastPageIcon from '@material-ui/icons/LastPage'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'

const styles = (theme: Theme) =>
  createStyles({
    root: {
      flexShrink: 0,
      color: theme.palette.text.secondary,
      marginLeft: theme.spacing.unit * 2.5,
    },
  })

type PageEvent = React.MouseEvent<HTMLButtonElement> | null

interface Props extends WithStyles<typeof styles> {
  count: number
  page: number
  rowsPerPage: number
  theme: Theme
  onChangePage: (event: PageEvent, page: number) => void
}

class PaginationActions extends React.Component<Props> {
  handleFirstPageButtonClick = (event: any) => {
    this.props.onChangePage(event, 0)
  }

  handleBackButtonClick = (event: any) => {
    this.props.onChangePage(event, this.props.page - 1)
  }

  handleNextButtonClick = (event: any) => {
    this.props.onChangePage(event, this.props.page + 1)
  }

  handleLastPageButtonClick = (event: any) => {
    this.props.onChangePage(
      event,
      Math.max(0, Math.ceil(this.props.count / this.props.rowsPerPage) - 1),
    )
  }

  render() {
    const { classes, count, page, rowsPerPage, theme } = this.props

    return (
      <div className={classes.root}>
        <IconButton
          disabled={page === 0}
          aria-label="First Page"
          onClick={this.handleFirstPageButtonClick}
        >
          {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
        </IconButton>
        <IconButton
          disabled={page === 0}
          aria-label="Previous Page"
          onClick={this.handleBackButtonClick}
        >
          {theme.direction === 'rtl' ? (
            <KeyboardArrowRight />
          ) : (
            <KeyboardArrowLeft />
          )}
        </IconButton>
        <IconButton
          disabled={page >= Math.ceil(count / rowsPerPage) - 1}
          aria-label="Next Page"
          onClick={this.handleNextButtonClick}
        >
          {theme.direction === 'rtl' ? (
            <KeyboardArrowLeft />
          ) : (
            <KeyboardArrowRight />
          )}
        </IconButton>
        <IconButton
          disabled={page >= Math.ceil(count / rowsPerPage) - 1}
          aria-label="Last Page"
          onClick={this.handleLastPageButtonClick}
        >
          {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
        </IconButton>
      </div>
    )
  }
}

export default withStyles(styles, { withTheme: true })(PaginationActions)
