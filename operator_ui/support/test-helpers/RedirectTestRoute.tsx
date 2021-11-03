import React from 'react'

import { Route } from 'react-router-dom'

// RedirectTestRoute is used to define a route which can be used for testing
// that a redirect occurs.
//
// Example
//
// render(
//   <MemoryRouter>
//     <Route exact path="/">
//       <Primary />
//     </Route>
//
//     <RedirectTestRoute
//       path="/new"
//       message="Successful redirect"
//     />
//   </MemoryRouter>
// )
//
//  When a redirect occurs in the '/' Primary component you can then assert
//  that the new page contains the message.
//
// expect(screen.queryByText('Successful redirect')).toBeInTheDocument()
export const RedirectTestRoute = ({
  path,
  message,
}: {
  path: string
  message: string
}) => {
  return (
    <Route exact path={path}>
      {() => <div>{message}</div>}
    </Route>
  )
}
