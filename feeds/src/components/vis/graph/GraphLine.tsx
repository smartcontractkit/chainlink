import React from 'react'

interface OwnProps {
  data: any
  position: any
}

interface Props extends OwnProps {}

const Line: React.FC<Props> = ({ data, position }) => {
  if (!data) return null

  return (
    <line
      className={`vis__line ${
        data.isFulfilled ? 'vis__line--fulfilled' : 'vis__line--pending'
      }`}
      {...position}
    ></line>
  )
}

export default Line
