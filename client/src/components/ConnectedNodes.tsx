import React, { Component } from 'react'
import io from 'socket.io-client'

type Props = {}
type State = {
  count?: number
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
    let count: string | number | undefined = this.state.count
    if (count === undefined) {
      count = '...'
    }

    return <div>Connected Nodes: {count}</div>
  }
}

export default ConnectedNodes
