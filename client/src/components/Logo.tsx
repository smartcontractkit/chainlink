import React from 'react'
import { Link } from '@reach/router'
import logo from '../images/logo.svg'

interface IProps {
  className: string
}

const Logo = ({ className }: IProps) => {
  return (
    <Link to="/" className={className}>
      <img src={logo} />
    </Link>
  )
}

export default Logo
