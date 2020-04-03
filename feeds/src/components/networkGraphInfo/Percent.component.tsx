import React from 'react'

export interface Props {
  value: number
}

const Percent: React.FC<Props> = ({ value }) => {
  return <>{value}%</>
}

export default Percent
