import React from 'react'
import { Link as ReactStaticLink } from 'react-static'

export default ({ children, to, className }) => (
  <ReactStaticLink className={className} to={to}>{children}</ReactStaticLink>
)
