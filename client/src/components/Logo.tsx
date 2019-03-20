import React from 'react'
import Link from '@material-ui/core/Link';

interface IProps {
  className: string
}

const Logo = (props: IProps) => {
  return (
    <Link href="/" className={props.className}>
      Chainlink Stats
    </Link>
  )
}

export default Logo
