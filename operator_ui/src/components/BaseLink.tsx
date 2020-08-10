import React from 'react'
import { Link } from 'react-router-dom'

interface Props {
  children: React.ReactNode
  href: string
  id?: string
  className?: string
}

const BaseLink = ({ children, href, id, className }: Props) => (
  <Link to={href} id={id} className={className}>
    {children}
  </Link>
)

export default BaseLink
