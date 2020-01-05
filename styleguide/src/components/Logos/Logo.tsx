import React from 'react'

interface Props {
  src: string
  href: string
  className?: string
  width?: number
  height?: number
  alt?: string
}

const Logo = ({
  src,
  href,
  className,
  width,
  height,
  alt = 'Chainlink Explorer',
}: Props) => {
  return (
    <a href={href} className={className}>
      <img src={src} width={width} height={height} alt={alt} />
    </a>
  )
}

export default Logo
