import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

interface IProps {
  children: React.ReactNode
  to: string
  className?: string
}

const Link = ({ children, to, className }: IProps) => {
  return (
    <ReactStaticLink className={className} to={to}>
      {children}
    </ReactStaticLink>
  )
}

export default Link
