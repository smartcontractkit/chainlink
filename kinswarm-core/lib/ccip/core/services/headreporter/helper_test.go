package headreporter

func (p *prometheusReporter) SetBackend(b PrometheusBackend) {
	p.backend = b
}
