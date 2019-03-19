import React, { Component } from 'react'
import io from 'socket.io-client'

type Props = {}
type State = {
  count?: number
}

const showCount = (count?: number): string => {
  if (count === undefined) {
    return '...'
  }
  return count.toString()
}

class ConnectedNodes extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    const socket = io('/', { path: '/client' })

    socket.on('clnodeCount', (count: number) => {
      this.setState({ count: count })
    })
  }

  render() {
    return <div>Connected Nodes: {showCount(this.state.count)}</div>
  }
}

export default ConnectedNodes
