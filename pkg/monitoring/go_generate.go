package monitoring

//go:generate mockery --name Metrics --inpackage  --structname MetricsMock --filename metrics_mock.go
//go:generate mockery --name Source --inpackage  --structname SourceMock --filename source_mock.go
//go:generate mockery --name SourceFactory --inpackage --structname SourceFactoryMock --filename source_factory_mock.go
//go:generate mockery --name Exporter --inpackage --structname ExporterMock --filename exporter_mock.go
//go:generate mockery --name ExporterFactory --inpackage --structname ExporterFactoryMock --filename exporter_factory_mock.go
