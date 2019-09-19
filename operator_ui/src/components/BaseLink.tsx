import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

interface Props {
  children: React.ReactNode
  href: string
  id?: string
  className?: string
}

const BaseLink = ({ children, href, id, className }: Props) => (
  <ReactStaticLink to={href} id={id} className={className}>
    {children}
  </ReactStaticLink>
)

export default BaseLink
