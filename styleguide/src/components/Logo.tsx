import React from 'react'
import { Image } from './Image'

interface Props {
  src: string
  width?: number
  height?: number
  alt?: string
}

export const Logo = (props: Props) => {
  return <Image {...props} />
}
