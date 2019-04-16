import React from 'react'
import { Link as ReactStaticLink } from 'react-router-dom'

const ReactStaticLinkComponent = ({ children, to, className }) => (
  <ReactStaticLink className={className} to={to}>
    {children}
  </ReactStaticLink>
)

export default ReactStaticLinkComponent
