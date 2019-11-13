import React from 'react'
import { Image } from './Image'

interface Props {
  src: string
  width?: number
  height?: number
  alt?: string
}

export const Logo: React.FC<Props> = props => {
  return <Image {...props} />
}
