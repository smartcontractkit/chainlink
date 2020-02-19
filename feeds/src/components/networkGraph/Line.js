import React from 'react'

const Line = ({ data, position }) => {
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
