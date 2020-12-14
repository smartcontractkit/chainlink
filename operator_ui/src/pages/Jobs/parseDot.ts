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
  const digraph = graphlibDot.read(dot)
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
