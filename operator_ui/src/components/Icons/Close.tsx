import React, { FC, SVGProps } from 'react'

const Close: FC<SVGProps<SVGSVGElement>> = (props) => (
  <svg
    style={{ cursor: 'pointer' }}
    width="46"
    height="46"
    viewBox="0 0 46 46"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <circle r="22" transform="matrix(1 0 0 -1 23 23)" stroke="#364864" />
    <rect
      width="1.4605"
      height="18.6214"
      transform="matrix(0.707107 0.707107 0.707107 -0.707107 16.1333 28.834)"
      fill="#686F84"
    />
    <rect
      x="30.3333"
      y="28.834"
      width="1.4605"
      height="18.6214"
      transform="rotate(135 30.3333 28.834)"
      fill="#686F84"
    />
  </svg>
)

export default Close
