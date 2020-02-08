import Avatar from '@material-ui/core/Avatar'
import ClickAwayListener from '@material-ui/core/ClickAwayListener'
import Fab from '@material-ui/core/Fab'
import Grow from '@material-ui/core/Grow'
import MenuList from '@material-ui/core/MenuList'
import Paper from '@material-ui/core/Paper'
import Popper from '@material-ui/core/Popper'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import classNames from 'classnames'
import React, { useRef, useState } from 'react'
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

interface Props extends WithStyles<typeof styles> {
  className?: string
}

const AvatarMenu: React.FC<Props> = ({ classes, className, children }) => {
  const anchorEl = useRef<HTMLElement>(null)
  const [open, setOpenState] = useState(false)
  const handleToggle = () => setOpenState(!open)
  const handleClose = () => setOpenState(false)

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
