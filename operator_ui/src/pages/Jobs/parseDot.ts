import graphlibDot from 'graphlib-dot'

export type Stratify = {
  id: string
  parentIds: string[]
  attributes?: { [key: string]: string }
}

type Edge = {
  v: string
  w: string
}

export function parseDot(dot: string): Stratify[] {
  // We want to permit the use of angle brackets to make the
  // specs more readable for multi-line task attributes. The backend
  // dot parsing library supports angle brackets which do not contain
  // valid HTML inside them, but the frontend dot library graphlibDot does not.
  // Since the dot parsing on the frontend is merely used to display the graph nodes
  // its fine to omit the angle bracket attributes.
  const dotNoAngleBrackets = dot.replace(/\w+\s*=\s*<([^>]|[\r\n])*>/g, '')
  const digraph = graphlibDot.read(dotNoAngleBrackets)
  const edges = digraph.edges()

  return digraph.nodes().map((id: string) => {
    const nodeInformation: Stratify = {
      id,
      parentIds: edges
        .filter((edge: Edge) => edge.w === id)
        .map((edge: Edge) => edge.v),
    }

    if (Object.keys(digraph.node(id)).length > 0) {
      nodeInformation.attributes = digraph.node(id)
    }

    return nodeInformation
  })
}
