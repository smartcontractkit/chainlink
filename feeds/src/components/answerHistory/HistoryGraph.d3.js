import * as d3 from 'd3'
import { formatAnswer } from 'contracts/utils'
import { humanizeUnixTimestamp } from 'utils'

export default class HistoryGraph {
  margin = { top: 30, right: 30, bottom: 30, left: 50 }
  width = 1300
  height = 250
  svg
  path
  tooltip
  x
  y
  line
  overlay
  tooltipPrice
  tooltipTimestamp
  config = {}

  constructor(config) {
    this.config = config
  }

  bisectDate = d3.bisector(d => d.timestamp).left

  build() {
    this.svg = d3
      .select('.answer-history-graph')
      .append('svg')
      .attr(
        'viewBox',
        `0 0 ${this.width} ${this.height +
          this.margin.top +
          this.margin.bottom}`,
      )

    this.bollinger = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')',
      )
      .style('opacity', this.config.bollinger ? 1 : 0)
      .attr('class', 'bollinger')

    this.bollingerArea = this.bollinger
      .append('path')
      .attr('class', 'bollinger-area')

    this.bollingerMa = this.bollinger
      .append('path')
      .attr('class', 'bollinger-ma')

    this.clip = this.svg
      .append('defs')
      .append('svg:clipPath')
      .attr('id', 'clip')
      .append('svg:rect')
      .attr('width', this.width - this.margin.left)
      .attr('height', this.height)
      .attr('x', 0)
      .attr('y', 0)

    this.path = this.svg
      .append('g')
      .attr('class', 'chart-line')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')',
      )
      .append('path')
      .attr('class', 'line')
      .attr('stroke', '#a0a0a0')
      .attr('fill', 'none')
      .attr('clip-path', 'url(#clip)')

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

    this.tooltipPrice = this.tooltip
      .append('text')
      .attr('class', 'answer-history-graph--price')
      .attr('x', '10')
      .attr('y', '-5')

    this.tooltipTimestamp = this.tooltip
      .append('text')
      .attr('class', 'answer-history-graph--timestamp')
      .attr('x', '10')
      .attr('y', '10')

    this.brushX = d3.brushX()

    this.brush = this.svg
      .append('g')
      .attr('class', 'brush')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')',
      )
      .on('mouseout', () => this.tooltip.style('display', 'none'))

    this.zoomOutBtn = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + (this.width - 70) + ',' + this.height + ')',
      )
      .style('opacity', 0)
      .style('cursor', 'pointer')

    this.zoomOutBtn
      .append('rect')
      .attr('class', 'y')
      .style('fill', '#375bd2')
      .style('stroke', '#375bd2')
      .attr('width', '70')
      .attr('height', '25')

    this.zoomOutBtn
      .append('text')
      .attr('fill', '#fff')
      .text('Zoom out')
      .attr('x', '10')
      .attr('y', '16')
      .style('font-size', 12)
  }

  updateBrushed() {
    if (!this.x) {
      return
    }

    const extent = d3.event.selection
    if (extent) {
      this.x.domain([this.x.invert(extent[0]), this.x.invert(extent[1])])
      this.brush.call(this.brushX.move, null)

      this.zoomOutBtn
        .transition()
        .duration(300)
        .style('opacity', 1)
    }

    this.xAxis
      .transition()
      .duration(300)
      .call(
        d3
          .axisBottom()
          .scale(this.x)
          .ticks(7)
          .tickFormat(f => humanizeUnixTimestamp(f)),
      )

    const line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.answer)))

    this.path
      .transition()
      .duration(300)
      .attr('d', line)
  }

  zoomOut(data) {
    this.x.domain(d3.extent(data, d => d.timestamp))

    const xAxis = d3
      .axisBottom()
      .scale(this.x)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.xAxis
      .transition()
      .duration(300)
      .call(xAxis)

    const line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.answer)))

    this.path
      .transition()
      .duration(300)
      .attr('d', line)

    this.zoomOutBtn
      .transition()
      .duration(300)
      .style('opacity', 0)
  }

  update(data) {
    if (!data) {
      return
    }

    this.x = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.timestamp))
      .range([0, this.width - this.margin.left])

    this.y = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.answer))
      .range([this.height, 0])

    const yAxis = d3
      .axisLeft()
      .scale(this.y)
      .ticks(4)
      .tickFormat(f =>
        formatAnswer(
          f,
          this.config.multiply,
          this.config.decimalPlaces,
          this.config.formatDecimalPlaces,
        ),
      )

    this.svg
      .append('g')
      .attr('class', 'y-axis')
      .attr(
        'transform',
        `translate(${this.margin.left - 10}, ${this.margin.top})`,
      )
      .call(yAxis)

    const xAxis = d3
      .axisBottom()
      .scale(this.x)
      .ticks(7)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.xAxis = this.svg
      .append('g')
      .attr('class', 'x-axis')
      .attr(
        'transform',
        `translate(${this.margin.left}, ${this.height + this.margin.top + 10})`,
      )
      .call(xAxis)

    this.line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.answer)))

    this.path.datum(data).attr('d', this.line)

    this.brush.on('mousemove', () => this.mousemove(data))

    this.updateMa(data)

    this.svg.on('dblclick', () => this.zoomOut(data))
    this.zoomOutBtn.on('click', () => this.zoomOut(data))

    this.brush.call(
      this.brushX
        .extent([
          [0, 0],
          [this.width - this.margin.left, this.height],
        ])
        .on('end', () => this.updateBrushed()),
    )
  }

  mousemove(data) {
    const x0 = this.x.invert(d3.mouse(this.brush.node())[0])
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
          (this.y(d.answer) + this.margin.top) +
          ')',
      )
    this.tooltipTimestamp.text(() =>
      humanizeUnixTimestamp(d.timestamp, 'LL LTS'),
    )
    this.tooltipPrice.text(
      () => `${this.config.valuePrefix} ${d.answerFormatted}`,
    )
  }

  getBollingerBands(n, k, data) {
    const bands = []
    for (let i = n - 1, len = data.length; i < len; i++) {
      const slice = data.slice(i + 1 - n, i)
      const mean = d3.mean(slice, d => d.answer)
      const stdDev = Math.sqrt(
        d3.mean(slice.map(d => Math.pow(d.answer - mean, 2))),
      )
      bands.push({
        timestamp: data[i].timestamp,
        answerId: data[i].answerId,
        ma: mean,
        low: mean - k * stdDev,
        high: mean + k * stdDev,
      })
    }

    return bands
  }

  updateMa(data) {
    const n = 20 // n-period of moving average
    const k = 2 // k times n-period standard deviation above/below moving average
    const bandsData = this.getBollingerBands(n, k, data)
    const x = d3.scaleTime().range([0, this.width - this.margin.left])
    const y = d3.scaleLinear().range([this.height, 0])

    x.domain(d3.extent(data, d => d.timestamp))
    y.domain(d3.extent(data, d => d.answer))

    const ma = d3
      .line()
      .x(d => x(d.timestamp))
      .y(d => y(d.ma))

    const bandsArea = d3
      .area()
      .x(d => x(d.timestamp))
      .y0(d => y(d.low))
      .y1(d => y(d.high))

    this.bollingerArea.datum(bandsData).attr('d', bandsArea)

    this.bollingerMa
      .datum(bandsData)
      .style('opacity', 0)
      .attr('d', ma)
  }

  toggleBollinger(toggle) {
    this.bollinger.style('opacity', toggle ? 1 : 0)
  }
}
