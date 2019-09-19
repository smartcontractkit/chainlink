import React from 'react'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import ReactResizeDetector from 'react-resize-detector'

const styles = (theme: Theme) =>
  createStyles({
    appBar: {
      backgroundColor: theme.palette.common.white,
      zIndex: theme.zIndex.modal + 1,
    },
    toolbar: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
    },
  })

interface Props extends WithStyles<typeof styles> {
  children: any
  onResize: (width: number, height: number) => void
}

const Header = ({ children, classes, onResize }: Props) => {
  return (
    <AppBar className={classes.appBar} color="default">
      <ReactResizeDetector
        refreshMode="debounce"
        refreshRate={200}
        handleWidth
        onResize={onResize}
      >
        <Toolbar className={classes.toolbar}>{children}</Toolbar>
      </ReactResizeDetector>
    </AppBar>
  )
}

export default withStyles(styles)(Header)
