import * as d3 from 'd3'
import { humanizeUnixTimestamp } from 'utils'

export default class DeviationGraph {
  margin = { top: 30, right: 30, bottom: 30, left: 50 }
  width = 1200
  height = 200
  svg
  path
  tooltip
  info
  x
  y
  line
  overlay
  tooltipPrice
  tooltipTimestamp
  tooltipPercentage
  options = {}
  colorRange = ['#c51515', '#c3c3c3']

  constructor(options) {
    this.options = options
  }

  bisectDate = d3.bisector(d => d.timestamp).left

  nearest(n) {
    return n < 0 ? Math.floor(n) : Math.round(n)
  }

  build() {
    this.svg = d3
      .select('.deviation-history-graph')
      .append('svg')
      .attr('viewBox', `0 0 ${1300} ${400}`)

    this.path = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')'
      )
      .append('path')
      .attr('class', 'line')
      .attr('stroke', '#a0a0a0')
      .attr('fill', 'none')

    this.tooltip = this.svg
      .append('g')
      .attr('class', 'tooltip')
      .style('display', 'none')

    this.tooltip
      .append('circle')
      .attr('class', 'y')
      .style('fill', '#2a59da')
      .style('stroke', '#2a59da')
      .attr('r', 4)

    this.info = this.svg
      .append('g')
      .attr('class', 'info')
      .attr(
        'transform',
        'translate(' + this.width + ',' + this.margin.top + ')'
      )

    this.tooltipPercentage = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--percentage')
      .attr('x', '10')
      .attr('y', '0')

    this.tooltipPrice = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--price')
      .attr('x', '10')
      .attr('y', '15')

    this.tooltipTimestamp = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--timestamp')
      .attr('x', '10')
      .attr('y', '30')

    this.overlay = this.svg
      .append('rect')
      .attr('width', this.width - this.margin.left)
      .attr('height', this.height)
      .style('fill', 'none')
      .style('pointer-events', 'all')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')'
      )
      .on('mouseout', () => {
        this.tooltip.style('display', 'none')
        this.info.style('display', 'none')
      })

    const color = d3
      .scaleLinear()
      .range(this.colorRange)
      .domain([1, 2])

    this.linearGradient = this.svg
      .append('defs')
      .append('linearGradient')
      .attr('id', 'linear-gradient')
      .attr('gradientTransform', 'rotate(90)')

    this.linearGradient
      .append('stop')
      .attr('offset', '0%')
      .attr('stop-color', color(1))

    this.linearGradient
      .append('stop')
      .attr('offset', '40%')
      .attr('stop-color', color(2))

    this.linearGradient
      .append('stop')
      .attr('offset', '60%')
      .attr('stop-color', color(2))

    this.linearGradient
      .append('stop')
      .attr('offset', '100%')
      .attr('stop-color', color(1))
  }

  addDeviation(data) {
    const reducedData = data.reduce((memo, current, index) => {
      const slice = memo.slice(index - 1, index)
      const average = d3.mean(slice, d => d.response) || current.response

      const deviation =
        (100 * (current.response - average)) /
        ((current.response + average) / 2)

      return [
        ...memo,
        ...[
          {
            ...current,
            ...{ deviation: Number(deviation.toFixed(3)) }
          }
        ]
      ]
    }, [])
    return reducedData
  }

  update(updatedData) {
    if (!updatedData) {
      return
    }

    const parsedData = JSON.parse(JSON.stringify(updatedData))
    const data = this.addDeviation(parsedData)

    this.x = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.timestamp))
      .range([0, this.width - this.margin.left])

    this.y = d3
      .scaleLinear()
      .domain(d3.extent(data, d => this.nearest(d.deviation)))
      .range([this.height, 0])

    const y_axis = d3
      .axisLeft()
      .scale(this.y)
      .ticks(3)
      .tickFormat(f => `${f}%`)

    this.svg
      .append('g')
      .attr('class', 'y-axis')
      .attr(
        'transform',
        `translate(${this.margin.left - 10}, ${this.margin.top})`
      )
      .call(y_axis)

    const x_axis = d3
      .axisBottom()
      .scale(this.x)
      .ticks(7)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.svg
      .append('g')
      .attr('class', 'x-axis')
      .attr(
        'transform',
        `translate(${this.margin.left}, ${this.height + this.margin.top + 10})`
      )
      .call(x_axis)

    this.line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.deviation)))
      .curve(d3.curveMonotoneX)

    this.path
      .attr('d', this.line(data))
      .attr('stroke-width', 1)
      .attr('stroke', 'url(#linear-gradient)')
      .attr('fill', 'none')

    const totalLength = this.path.node().getTotalLength()

    this.path
      .attr('stroke-dasharray', totalLength + ' ' + totalLength)
      .attr('stroke-dashoffset', totalLength)
      .transition()
      .duration(2000)
      .attr('stroke-dashoffset', 0)

    this.overlay.on('mousemove', () => this.mousemove(data))
  }

  mousemove(data) {
    const x0 = this.x.invert(d3.mouse(this.overlay.node())[0])
    const i = this.bisectDate(data, x0, 1)
    const d0 = data[i - 1]
    const d1 = data[i]
    if (!d1) {
      return
    }
    const d = x0 - d0.timestamp > d1.timestamp - x0 ? d1 : d0
    this.tooltip
      .style('display', 'block')
      .attr(
        'transform',
        'translate(' +
          (this.x(d.timestamp) + this.margin.left) +
          ',' +
          (this.y(d.deviation) + this.margin.top) +
          ')'
      )

    this.info.style('display', 'block')
    this.tooltipTimestamp.text(() => humanizeUnixTimestamp(d.timestamp))
    this.tooltipPrice.text(
      () => `${this.options.valuePrefix} ${d.responseFormatted}`
    )
    this.tooltipPercentage.text(() => `${d.deviation.toFixed(2)}%`)
  }
}
