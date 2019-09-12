import React from 'react'
import logo from '../images/logo.svg'

interface Props {
  className?: string
  width?: number
  height?: number
}

const Logo = ({ className, width, height }: Props) => {
  return (
    <a href="/" className={className}>
      <img src={logo} width={width} height={height} alt="Chainlink Explorer" />
    </a>
  )
}

export default Logo
