import React from 'react'
import Logo from '../Logo'
import src from './public.svg'

interface Props {
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

export const PublicLogo = (props: Props) => {
  return <Logo src={src} {...props} />
}
