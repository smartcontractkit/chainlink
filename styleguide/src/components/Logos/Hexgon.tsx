import React from 'react'
import { Logo } from '../Logo'
import src from './icon-logo-blue.svg'

interface Props {
  width?: number
  height?: number
  alt?: string
}

const Hexagon = (props: Props) => {
  return <Logo src={src} {...props} />
}

export default Hexagon
