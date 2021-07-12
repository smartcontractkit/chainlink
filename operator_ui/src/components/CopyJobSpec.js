import React from 'react'
import Copy from 'components/Copy'

const CopyJobSpec = ({ JobSpec, ...props }) => (
  <Copy buttonText="Copy JobSpec" data={JobSpec} {...props} />
)

export default CopyJobSpec
