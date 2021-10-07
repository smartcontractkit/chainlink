import React from 'react'
import { Link } from 'react-router-dom'
import MenuItem, { MenuItemProps } from '@material-ui/core/MenuItem'

interface MenuItemLinkProps extends MenuItemProps {
  to: string
  replace?: boolean
}

export const MenuItemLink = (props: MenuItemLinkProps) => {
  return (
    <MenuItem button {...props} component={Link as any}>
      {props.children}
    </MenuItem>
  )
}
