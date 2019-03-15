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
    const socket = io('/')
    socket.on('connectionCount', (count: number) => {
      this.setState({ count: count })
    })
  }

  render() {
    return (
      <div>
        Connected Nodes: {this.state.count || '...'}
      </div>
    )
  }
}

export default ConnectedNodes
