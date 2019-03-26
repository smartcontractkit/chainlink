import React from 'react'
import Link from '@material-ui/core/Link'

interface IProps {
  className: string
}

const Logo = ({ className }: IProps) => {
  return (
    <Link href="/" className={className}>
      Chainlink Stats
    </Link>
  )
}

export default Logo
