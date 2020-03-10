import React from 'react'

interface Props {
  width: number
  height: number
  fill: string
}

const ChainlinkLogo: React.FC<Props> = ({
  width = 37.8,
  height = 43.6,
  fill = '#375BD2',
}) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 37.8 43.6"
      width={width}
      height={height}
    >
      <title>Chainlink</title>
      <path
        d="M18.9 0l-4 2.3L4 8.6l-4 2.3V32.7L4 35l11 6.3 4 2.3 4-2.3L33.8 35l4-2.3V10.9l-4-2.3-10.9-6.3-4-2.3zM8 28.1V15.5l10.9-6.3 10.9 6.3v12.6l-10.9 6.3L8 28.1z"
        fill={fill}
      />
    </svg>
  )
}

export default ChainlinkLogo
