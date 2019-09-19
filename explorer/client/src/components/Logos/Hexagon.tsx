import React from 'react'
import Logo from '../Logo'
import src from './hexagon.svg'

interface Props {
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

export const HexagonLogo = (props: Props) => {
  return <Logo src={src} {...props} />
}
