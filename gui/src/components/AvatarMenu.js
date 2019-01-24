import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import { useHooks, useState, useRef } from 'use-react-hooks'
import { Link } from 'react-router-dom'
import Popper from '@material-ui/core/Popper'
import Grow from '@material-ui/core/Grow'
import Fab from '@material-ui/core/Fab'
import Paper from '@material-ui/core/Paper'
import MenuItem from '@material-ui/core/MenuItem'
import MenuList from '@material-ui/core/MenuList'
import ClickAwayListener from '@material-ui/core/ClickAwayListener'
import Avatar from '@material-ui/core/Avatar'
import Typography from '@material-ui/core/Typography'
import { grey } from '@material-ui/core/colors'
import face from 'images/face.svg'

const styles = theme => {
  return {
    button: {
      marginTop: -4
    },
    avatar: {
      width: 30,
      height: 30
    },
    menuListGrow: {
      marginTop: 10,
      borderRadius: theme.spacing.unit * 2,
      backgroundColor: theme.palette.primary.main
    },
    menuItem: {
      '&:hover': {
        backgroundColor: 'transparent'
      }
    },
    link: {
      color: theme.palette.common.white,
      textDecoration: 'none',
      '&:hover': {
        color: grey[200]
      }
    }
  }
}

const AvatarMenu = useHooks(({ classes }) => {
  const anchorEl = useRef(null)
  const [open, setOpenState] = useState(false)
  const handleToggle = () => setOpenState(!open)

  const handleClose = event => {
    if (anchorEl.current.contains(event.target)) {
      return
    }

    setOpenState(false)
  }

  return (
    <React.Fragment>
      <Fab
        size='medium'
        color='primary'
        aria-label='Profile'
        className={classes.button}
        buttonRef={anchorEl}
        aria-owns={open ? 'menu-list-grow' : undefined}
        aria-haspopup='true'
        onClick={handleToggle}
      >
        <Avatar alt='Profile' src={face} className={classes.avatar} />
      </Fab>
      <Popper open={open} anchorEl={anchorEl.current} transition disablePortal>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            id='menu-list-grow'
            className={classes.menuListGrow}
            style={{ transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom' }}
          >
            <Paper square={false}>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList>
                  <MenuItem onClick={handleClose} className={classes.menuItem}>
                    <Link to='/signout' className={classes.link}>
                      <Typography
                        variant='body1'
                        className={classes.link}
                      >
                        Log out
                      </Typography>
                    </Link>
                  </MenuItem>
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </React.Fragment>
  )
})

export default withStyles(styles)(AvatarMenu)
