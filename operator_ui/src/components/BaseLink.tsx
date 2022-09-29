import React from 'react'
import { Link } from 'react-router-dom'

interface Props {
  children: React.ReactNode
  href: string
  id?: string
  className?: string
  onClick?: (event: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => void
}

const BaseLink = ({ children, href, id, className, onClick }: Props) => (
  <Link to={href} id={id} className={className} onClick={onClick}>
    {children}
  </Link>
)

export default BaseLink
