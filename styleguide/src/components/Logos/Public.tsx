import React from 'react'
import Logo from './Logo'
import src from './explorer-public.svg'

interface Props {
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

export const PublicLogo: React.FC<Props> = props => {
  return <Logo src={src} {...props} />
}
