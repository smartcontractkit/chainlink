import React from 'react'

interface Props {
  className?: string
  children: React.ReactNode
}

const SearchBox = ({ className, children }: Props) => {
  return (
    <form method="GET" action="/job-runs" className={className}>
      {children}
    </form>
  )
}

export default SearchBox
