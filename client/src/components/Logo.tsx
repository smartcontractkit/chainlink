import React from 'react'
import Link from '@material-ui/core/Link'
import logo from '../images/logo.svg'

interface IProps {
  className: string
}

const Logo = ({ className }: IProps) => {
  return (
    <Link href="/" className={className}>
      <img src={logo} />
    </Link>
  )
}

export default Logo
