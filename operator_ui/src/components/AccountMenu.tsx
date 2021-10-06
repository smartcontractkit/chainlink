import React from 'react'
import { useDispatch } from 'react-redux'

import AccountCircleIcon from '@material-ui/icons/AccountCircle'
import IconButton from '@material-ui/core/IconButton'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import { Theme, withStyles, WithStyles } from '@material-ui/core/styles'

import { beginRegistration, submitSignOut } from 'actionCreators'

const styles = (theme: Theme) => {
  return {
    accountButton: {
      fontSize: 32,
      color: theme.palette.primary.main,
    },
    menuList: {
      paddingTop: 0,
      paddingBottom: 0,
    },
  }
}

interface Props extends WithStyles<typeof styles> {}

export const AccountMenu = withStyles(styles)(({ classes }: Props) => {
  const dispatch = useDispatch()
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleLogOut = () => {
    dispatch(submitSignOut())
    handleClose()
  }

  const handleRegisterMFA = () => {
    dispatch(beginRegistration())
    handleClose()
  }

  return (
    <React.Fragment>
      <IconButton disableRipple onClick={handleOpen}>
        <AccountCircleIcon className={classes.accountButton} />
      </IconButton>

      <Menu
        id="account-menu"
        anchorEl={anchorEl}
        getContentAnchorEl={null}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        transformOrigin={{ vertical: 8, horizontal: 64 }}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        disableAutoFocusItem
        MenuListProps={{
          className: classes.menuList,
        }}
      >
        <MenuItem onClick={handleRegisterMFA}>Register MFA Token</MenuItem>

        <MenuItem onClick={handleLogOut}>Log out</MenuItem>
      </Menu>
    </React.Fragment>
  )
})
