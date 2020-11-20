import React from 'react'
import { theme } from 'index'
import Typography from '@material-ui/core/Typography'
import * as d3dag from 'd3-dag'
import * as d3 from 'd3'
import { parseDot, Stratify } from './parseDot'

type Node = {
  x: number
  y: number
} & Stratify

type NodeElement = {
  data: Stratify
  x: number
  y: number
}

function createDag({
  dotSource,
  ref,
  setTooltip,
}: {
  dotSource: string
  ref: HTMLInputElement
  setTooltip: Function
}): void {
  const nodeRadius = 18
  const stratify = parseDot(`digraph {${dotSource}}`)
  const width = ref.offsetWidth
  const height = stratify.length * 60

  const svgSelection = d3
    .select(ref)
    .append('svg')
    .attr('width', width)
    .attr('height', height + nodeRadius * 2)
    .attr('style', 'overflow: visible')
    .attr('viewBox', `${0} ${0} ${width} ${height}`)

  const groupSelection = svgSelection
    .append('g')
    .attr('transform', `translate(${nodeRadius * 2}, 0)`)

  const dag = d3dag.dagStratify()(stratify)

  d3dag
    .sugiyama()
    .size([width - 150, height])
    .layering(d3dag.layeringSimplex())
    .decross(d3dag.decrossOpt())
    .coord(d3dag.coordVert())(dag)

  const line = d3
    .line<NodeElement>()
    .curve(d3.curveCatmullRom)
    .x((node) => node.x)
    .y((node) => node.y)

  groupSelection
    .append('g')
    .selectAll('path')
    .data(dag.links())
    .enter()
    .append('path')
    .attr('d', ({ points }) => line(points))
    .attr('fill', 'none')
    .attr('stroke-width', 2)
    .attr('stroke', '#3c40c6')

  const nodes = groupSelection
    .append('g')
    .selectAll('g')
    .data(dag.descendants())
    .enter()
    .append('g')
    .attr('style', 'cursor: default')
    .attr('id', (node) => node.id)
    .attr('transform', ({ x, y }) => `translate(${x}, ${y})`)
    .on('mouseover', (_, node) => {
      setTooltip(node)
      d3.select(`#circle-${node.data.id}`)
        .transition()
        .attr('r', nodeRadius + 3)
        .duration('50')
    })
    .on('mouseout', (_, node) => {
      setTooltip(null)
      d3.select(`#circle-${node.data.id}`)
        .transition()
        .attr('r', nodeRadius)
        .duration('50')
    })

  nodes
    .append('circle')
    .attr('id', (node) => {
      return `circle-${node.data.id}`
    })
    .attr('r', nodeRadius)
    .attr('fill', 'black')
    .attr('stroke', 'white')
    .attr('stroke-width', 6)
    .attr('fill', '#3c40c6')

  nodes
    .append('text')
    .text((node) => node.data.id)
    .attr('x', 30)
    .attr('font-weight', 'normal')
    .attr('font-family', 'sans-serif')
    .attr('text-anchor', 'start')
    .attr('font-size', '1em')
    .attr('alignment-baseline', 'middle')
    .attr('fill', 'black')
}

interface Props {
  dotSource: string
}

const TaskList = ({ dotSource }: Props) => {
  const [tooltip, setTooltip] = React.useState<NodeElement>()
  const graph = React.useRef<HTMLInputElement>(null)

  React.useEffect(() => {
    if (graph.current) {
      createDag({ dotSource, ref: graph.current, setTooltip })
    }
  }, [dotSource])

  return (
    <div style={{ position: 'relative' }}>
      {tooltip && (
        <div
          style={{
            position: 'absolute',
            left: '-305px',
            border: '1px solid rgba(0, 0, 0, 0.1)',
            padding: theme.spacing.unit,
            background: 'white',
            borderRadius: 5,
            minWidth: '300px',
            transform: `translate(${tooltip.x}px, ${tooltip.y}px)`,
          }}
        >
          <Typography variant="body1" color="textPrimary">
            <b>{tooltip.data.id}</b>
          </Typography>
          {tooltip.data?.attributes &&
            Object.entries(tooltip.data.attributes).map(([key, value]) => (
              <div key={key}>
                <Typography
                  variant="body1"
                  color="textSecondary"
                  component="div"
                >
                  <b>{key}:</b> {value}
                </Typography>
              </div>
            ))}
        </div>
      )}
      <div
        id="graph"
        style={{ padding: `${theme.spacing.unit * 3}px 0px` }}
        ref={graph}
      />
    </div>
  )
}

export default TaskList
