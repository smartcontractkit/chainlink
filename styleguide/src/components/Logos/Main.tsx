import React from 'react'
import src from './operator-public.svg'
import Logo from './Logo'

interface Props {
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

export const Main: React.FC<Props> = props => {
  return <Logo src={src} alt="Chainlink Operator" {...props} />
}
