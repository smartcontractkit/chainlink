import React from 'react'
import Button from 'components/Button'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({ data, buttonText, className }) => (
  <CopyToClipboard text={data}>
    <Button className={className}>{buttonText}</Button>
  </CopyToClipboard>
)

export default Copy
