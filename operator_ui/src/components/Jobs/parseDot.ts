import graphlibDot from 'graphlib-dot'

export type Stratify = {
  id: string
  parentIds: string[]
  attributes: { [key: string]: string }
}

type Edge = {
  v: string
  w: string
}

export function parseDot(dot: string): Stratify[] {
  const digraph = graphlibDot.read(dot)
  const edges = digraph.edges()

  return digraph.nodes().map((node: string) => ({
    id: node,
    parentIds: edges
      .filter((edge: Edge) => edge.w === node)
      .map((edge: Edge) => edge.v),
    attributes: digraph.node(node),
  }))
}
