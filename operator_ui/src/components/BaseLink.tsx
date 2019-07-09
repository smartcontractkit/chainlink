import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

interface IProps {
  children: React.ReactNode
  href: string
  id?: string
  className?: string
}

const BaseLink = ({ children, href, id, className }: IProps) => (
  <ReactStaticLink to={href} id={id} className={className}>
    {children}
  </ReactStaticLink>
)

export default BaseLink
