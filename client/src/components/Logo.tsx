import React from 'react'
import logo from '../images/logo.svg'

interface IProps {
  className?: string
}

const Logo = ({ className }: IProps) => {
  return (
    <a href="/" className={className}>
      <img src={logo} />
    </a>
  )
}

export default Logo
