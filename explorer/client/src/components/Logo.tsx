import React from 'react'
import logo from '../images/logo.svg'

interface IProps {
  className?: string
  width?: number
  height?: number
}

const Logo = ({ className, width, height }: IProps) => {
  return (
    <a href="/" className={className}>
      <img src={logo} width={width} height={height} />
    </a>
  )
}

export default Logo
