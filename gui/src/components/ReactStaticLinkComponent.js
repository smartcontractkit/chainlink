import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

export default ({ children, to, className }) => (
  <ReactStaticLink className={className} to={to}>{children}</ReactStaticLink>
)
