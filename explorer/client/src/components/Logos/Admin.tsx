import React from 'react'
import logo from './admin.svg'

interface Props {
  className?: string
  width?: number
  height?: number
}

export const AdminLogo = ({ className, width, height }: Props) => {
  return (
    <a href="/admin" className={className}>
      <img
        src={logo}
        width={width}
        height={height}
        alt="Chainlink Explorer Admin"
      />
    </a>
  )
}
