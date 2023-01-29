package port

import "github.com/hugosrc/shortlink/internal/core/domain"

type MetricsProducer interface {
	Produce(metrics *domain.LinkMetrics) error
}
