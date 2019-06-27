import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

interface IProps {
  children: React.ReactNode
  to: string
  className?: string
}

const BaseLink = ({ children, to, className }: IProps) => (
  <ReactStaticLink className={className} to={to}>
    {children}
  </ReactStaticLink>
)

export default BaseLink
