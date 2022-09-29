import { IconButton, withStyles } from '@material-ui/core'
import FirstPageIcon from '@material-ui/icons/FirstPage'
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft'
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight'
import LastPageIcon from '@material-ui/icons/LastPage'
import React from 'react'

const styles = (theme) => ({
  customButtons: {
    flexShrink: 0,
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing.unit * 2.5,
  },
})

const TableButtons = (props) => {
  const firstPage = 1
  const currentPage = props.page
  const lastPage = Math.ceil(props.count / props.rowsPerPage)
  const handlePage = (page) => (e) => {
    page = Math.min(page, lastPage)
    page = Math.max(page, firstPage)
    if (props.history) props.history.push(`${props.replaceWith}/${page}`)
    props.onChangePage(e, page)
  }

  return (
    <div className={props.classes.customButtons}>
      <IconButton
        onClick={handlePage(firstPage)}
        disabled={currentPage <= firstPage}
        aria-label="First Page"
      >
        <FirstPageIcon />
      </IconButton>
      <IconButton
        onClick={handlePage(currentPage - 1)}
        disabled={currentPage <= firstPage}
        aria-label="Previous Page"
      >
        <KeyboardArrowLeft />
      </IconButton>
      <IconButton
        onClick={handlePage(currentPage + 1)}
        disabled={currentPage >= lastPage}
        aria-label="Next Page"
      >
        <KeyboardArrowRight />
      </IconButton>
      <IconButton
        onClick={handlePage(lastPage)}
        disabled={currentPage >= lastPage}
        aria-label="Last Page"
      >
        <LastPageIcon />
      </IconButton>
    </div>
  )
}

export const FIRST_PAGE = 1
export default withStyles(styles)(TableButtons)
