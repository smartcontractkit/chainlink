import React from 'react'
import src from './no-activity-icon.svg'
import Logo from './Logo'

interface Props {
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

export const NoContent: React.FC<Props> = props => {
  return <Logo src={src} alt="No Content" {...props} />
}
