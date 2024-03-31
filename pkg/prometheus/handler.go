package prometheus

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FiberPromClient struct {
	histogramVec *prometheus.HistogramVec
}

func NewFiberPromClient() FiberPromClient {
	httpRequestProm := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.ExponentialBuckets(1, 2, 20),
	}, []string{"path", "method", "status"})

	return FiberPromClient{
		histogramVec: httpRequestProm,
	}
}

func (m *FiberPromClient) Register(app fiber.Router) {
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

// Echo example
func (m *FiberPromClient) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Route().Path == "/metrics" {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		elapsed := float64(time.Since(start).Nanoseconds()) / 1e6

		status := fiber.StatusInternalServerError
		if err != nil {
			var e *fiber.Error
			if errors.As(err, &e) {
				status = e.Code
			}
		} else {
			status = c.Response().StatusCode()
		}

		path := c.Route().Path
		method := c.Route().Method

		m.histogramVec.WithLabelValues(path, method, strconv.Itoa(status)).Observe(elapsed)

		return err
	}
}
