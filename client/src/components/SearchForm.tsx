import React from 'react'

interface IProps {
  className?: string
  children: React.ReactNode
}

const SearchBox = ({ className, children }: IProps) => {
  return (
    <form method="GET" action="/job-runs" className={className}>
      {children}
    </form>
  )
}

export default SearchBox
