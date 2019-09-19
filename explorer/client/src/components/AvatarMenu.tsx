import React, { useState, useRef } from 'react'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import Popper from '@material-ui/core/Popper'
import Grow from '@material-ui/core/Grow'
import Fab from '@material-ui/core/Fab'
import Paper from '@material-ui/core/Paper'
import ClickAwayListener from '@material-ui/core/ClickAwayListener'
import MenuList from '@material-ui/core/MenuList'
import Avatar from '@material-ui/core/Avatar'
import classNames from 'classnames'
import face from './face.svg'

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    button: {
      marginTop: -4,
    },
    avatar: {
      width: 30,
      height: 30,
    },
    menuListGrow: {
      marginTop: 10,
      borderRadius: spacing.unit * 2,
      backgroundColor: palette.primary.main,
    },
  })

interface Item {
  text: string
}

interface Props extends WithStyles<typeof styles> {
  children: any
  className?: any
}

const AvatarMenu = ({ classes, className, children }: Props) => {
  const anchorEl = useRef<HTMLElement>(null)
  const [open, setOpenState] = useState(false)
  const handleToggle = () => setOpenState(!open)

  const handleClose = () => {
    setOpenState(false)
  }

  return (
    <>
      <Fab
        size="medium"
        color="primary"
        aria-label="Profile"
        className={classNames(classes.button, className)}
        buttonRef={anchorEl}
        aria-owns={open ? 'menu-list-grow' : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <Avatar alt="Profile" src={face} className={classes.avatar} />
      </Fab>
      <Popper open={open} anchorEl={anchorEl.current} transition disablePortal>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper square={false}>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList className={classes.menuListGrow}>{children}</MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  )
}

export default withStyles(styles)(AvatarMenu)
