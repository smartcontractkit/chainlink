import React from 'react'

interface Props {
  width?: number
  height?: number
  className?: string
  'data-testid'?: string
}

const Success = ({ width, height, className, ...rest }: Props) => (
  <svg
    data-name="Layer 1"
    viewBox="0 0 48 48"
    width={width}
    height={height}
    className={className}
    data-testid={rest['data-testid']}
  >
    <defs>
      <mask
        id="prefix__a"
        x={13.75}
        y={16}
        width={21}
        height={15}
        maskUnits="userSpaceOnUse"
      >
        <path fillRule="evenodd" fill="#fff" d="M13.75 16h21v15h-21V16z" />
      </mask>
    </defs>
    <title>{'Artboard 1'}</title>
    <circle cx={24} cy={24} r={22} fill="none" stroke="#fff" strokeWidth={4} />
    <circle data-name="path-1" cx={24} cy={24} r={20} fill="#e8faf1" />
    <g>
      <g mask="url(#prefix__a)">
        <path
          d="M30.81 17a1.89 1.89 0 0 0-1.29.56c-2.94 2.9-5.2 5.33-7.89 8.05l-2.74-2.28a1.9 1.9 0 0 0-2.66.14 1.84 1.84 0 0 0 .14 2.62l.09.07 4.08 3.4a1.9 1.9 0 0 0 2.55-.11c3.38-3.34 5.79-6 9.09-9.27a1.82 1.82 0 0 0 0-2.62 1.92 1.92 0 0 0-1.4-.56"
          fill="#17c671"
          fillRule="evenodd"
        />
      </g>
    </g>
  </svg>
)

export default Success
