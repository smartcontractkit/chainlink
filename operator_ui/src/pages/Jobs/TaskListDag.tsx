import React from 'react'
import { theme } from 'theme'
import Typography from '@material-ui/core/Typography'
import * as d3dag from 'd3-dag'
import * as d3 from 'd3'
import { Stratify } from './parseDot'

type Node = {
  x: number
  y: number
}

type NodeElement = {
  data: Stratify
  x: number
  y: number
}

function createDag({
  stratify,
  ref,
  setTooltip,
}: {
  stratify: Stratify[]
  ref: HTMLInputElement
  setTooltip: Function
}): void {
  const nodeRadius = 18
  const width = ref.offsetWidth
  const height = stratify.length * 60

  // Clean up
  d3.select(ref).select('svg').remove()

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

  const dag = d3dag.dagStratify<Stratify>()(stratify)

  d3dag
    .sugiyama()
    .size([width - 150, height])
    .layering(d3dag.layeringSimplex())
    .decross(d3dag.decrossOpt())
    .coord(d3dag.coordVert())(dag)

  const line = d3
    .line<Node>()
    .curve(d3.curveCatmullRom)
    .x((node) => node.x)
    .y((node) => node.y)

  // Styling links
  groupSelection
    .append('g')
    .selectAll('path')
    .data(dag.links())
    .enter()
    .append('path')
    .attr('d', ({ points }) => line(points))
    .attr('fill', 'none')
    .attr('stroke-width', 2)
    .attr('stroke', theme.palette.grey['300'])

  const nodes = groupSelection
    .append('g')
    .selectAll('g')
    .data(dag.descendants())
    .enter()
    .append('g')
    .attr('style', 'cursor: default')
    .attr('id', (node) => node.id)
    .attr('transform', ({ x, y }: any) => `translate(${x}, ${y})`)
    .on('mouseover', (_, node) => {
      setTooltip(node)
      d3.select<d3.BaseType, NodeElement>(`#circle-${node.data.id}`)
        .transition()
        .attr('r', nodeRadius + 3)
        .duration(50)
    })
    .on('mouseout', (_, node) => {
      setTooltip(null)
      d3.select(`#circle-${node.data.id}`)
        .transition()
        .attr('r', nodeRadius)
        .duration(50)
    })

  // Styling dots
  nodes
    .append('circle')
    .attr('id', (node) => {
      return `circle-${node.data.id}`
    })
    .attr('r', nodeRadius)
    .attr('fill', 'black')
    .attr('stroke', 'white')
    .attr('stroke-width', 6)
    .attr('fill', (node) => {
      if (node.data.attributes) {
        switch (node.data.attributes.status) {
          case 'in_progress':
            // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
            // @ts-ignore because material UI doesn't update theme types with options
            return theme.palette.warning.main
          case 'completed':
            // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
            // @ts-ignore because material UI doesn't update theme types with options
            return theme.palette.success.main
          case 'errored':
            return theme.palette.error.main
          case 'aborted':
            return theme.palette.grey['500']
          default:
            return theme.palette.primary.main
        }
      }
      return theme.palette.primary.main
    })

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
  stratify: Stratify[]
}

export const TaskList = ({ stratify }: Props) => {
  const [tooltip, setTooltip] = React.useState<NodeElement>()
  const graph = React.useRef<HTMLInputElement>(null)

  React.useEffect(() => {
    if (graph.current) {
      createDag({ stratify, ref: graph.current, setTooltip })
    }
  }, [stratify])

  React.useEffect(() => {
    function handleResize() {
      if (graph.current) {
        createDag({ stratify, ref: graph.current, setTooltip })
      }
    }

    window.addEventListener('resize', handleResize)

    return () => window.removeEventListener('resize', handleResize)
  }, [stratify, graph])

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
            width: '300px',
            transform: `translate(${tooltip.x}px, ${tooltip.y}px)`,
          }}
        >
          <Typography variant="body1" color="textPrimary">
            <b>{tooltip.data.id}</b>
          </Typography>
          {tooltip.data?.attributes &&
            Object.entries(tooltip.data.attributes)
              .filter(([key]) => key !== 'error') // We want to filter errors out as they can get quite long
              .map(([key, value]) => (
                <div key={key}>
                  <Typography
                    variant="body1"
                    color="textSecondary"
                    component="div"
                  >
                    <b>{key}:</b> {value || '-'}
                  </Typography>
                </div>
              ))}
        </div>
      )}
      <div
        id="graph"
        style={{
          display: 'flex',
          justifyContent: 'center',
          marginLeft: theme.spacing.unit * 3,
          marginRight: theme.spacing.unit * 3,
          paddingTop: theme.spacing.unit * 3,
          paddingBottom: theme.spacing.unit * 3,
        }}
        ref={graph}
      />
    </div>
  )
}

export default TaskList
