package http_server

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"runtime/metrics"
)

type ApplicationMetricsRequest struct {
	XMLName xml.Name `xml:"requested_metrics"`
	Metrics []string `xml:"name"`
}

type ApplicationMetrics struct {
	XMLName xml.Name     `xml:"metrics"`
	Metrics []*MetricDto `xml:"metric"`
}

type MetricDto struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	KindString  string   `xml:"kind"`
	Cumulative  bool     `xml:"cumulative"`
	Value       ValueDto `xml:"value,omitempty"`
}

type ValueDto struct {
	Type  string `xml:"type,attr"`
	Value any    `xml:",chardata"`
}

type MetricsHandler struct {
}

// GetAllMetrics Get all metrics from metrics.All() in xml format
func (m *MetricsHandler) GetAllMetrics(w http.ResponseWriter, _ *http.Request) {
	var err error
	var result []byte
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result = []byte(err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_, _ = w.Write(result)
	}()

	allMetrics := metrics.All()
	metricsDto := &ApplicationMetrics{}
	for _, metric := range allMetrics {
		metricsDto.Metrics = append(metricsDto.Metrics, &MetricDto{
			Name:        metric.Name,
			Description: metric.Description,
			KindString:  getValueKingString(metric.Kind),
			Cumulative:  metric.Cumulative,
		})
	}

	result, err = xml.MarshalIndent(metricsDto, " ", " ")
}

/*
GetMetricsValues Get metrics values by filter
Most useful:
CPU:
/cpu/classes/total:cpu-seconds - Estimated total available CPU time for user Go code or the Go runtime, as defined by GOMAXPROCS
/cpu/classes/user:cpu-seconds - Estimated total CPU time spent running user Go code
For goroutines (scheduler and sync):
/sched/goroutines:goroutines - Count of live goroutines
/sched/latencies:seconds - Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running
/sync/mutex/wait/total:seconds - Approximate cumulative time goroutines have spent blocked on a sync.Mutex or sync.RWMutex. This metric is useful for identifying global changes in lock contention. Collect a mutex or block profile using the runtime/pprof package for more detailed contention data.
GC:
/cpu/classes/gc/pause:cpu-seconds - Estimated total CPU time spent with the application paused by the GC
/cpu/classes/gc/total:cpu-seconds - Estimated total CPU time spent performing GC tasks
/gc/heap/allocs:bytes - Cumulative sum of memory allocated to the heap by the application
Memory:
/memory/classes/heap/free:bytes - Memory that is completely free and eligible to be returned to the underlying system, but has not been
/memory/classes/heap/stacks:bytes - Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use
*/
func (m *MetricsHandler) GetMetricsValues(w http.ResponseWriter, r *http.Request) {
	var err error
	var result []byte
	defer func() {
		if recErr := recover(); recErr != nil {
			err = fmt.Errorf("%v", recErr)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result = []byte(err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_, _ = w.Write(result)
	}()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	request := &ApplicationMetricsRequest{}
	err = xml.Unmarshal(body, request)
	if err != nil {
		return
	}

	metricsDto := &ApplicationMetrics{}
	for _, reqMetric := range request.Metrics {
		isFound := false
		for _, metric := range metrics.All() {
			if metric.Name != reqMetric {
				continue
			}
			if metric.Kind == metrics.KindBad {
				err = fmt.Errorf("unsupported metric type (%s) for requested metric %s",
					getValueKingString(metric.Kind), reqMetric)
				return
			}
			readMetrics := make([]metrics.Sample, 1)
			readMetrics[0].Name = reqMetric
			metrics.Read(readMetrics)
			sampleValue := readMetrics[0]
			if metric.Kind != sampleValue.Value.Kind() {
				err = fmt.Errorf("unknown metric value %s after read from metrics %s",
					getValueKingString(sampleValue.Value.Kind()), metric.Name)
				return
			}
			metricDto := &MetricDto{
				Name:        metric.Name,
				Description: metric.Description,
				KindString:  getValueKingString(metric.Kind),
				Cumulative:  metric.Cumulative,
			}
			switch metric.Kind {
			case metrics.KindUint64:
				metricDto.Value.Type = "uint64"
				metricDto.Value.Value = sampleValue.Value.Uint64()
			case metrics.KindFloat64:
				metricDto.Value.Type = "float64"
				metricDto.Value.Value = sampleValue.Value.Float64()
			case metrics.KindFloat64Histogram:
				metricDto.Value.Value = medianBucket(sampleValue.Value.Float64Histogram())
				metricDto.Value.Type = "histogram-median"
			default:
				panic("unknown metric kind")
			}
			metricsDto.Metrics = append(metricsDto.Metrics, metricDto)
			isFound = true
			break
		}
		if !isFound {
			err = fmt.Errorf("unknown metric %s", reqMetric)
			return
		}
	}
	result, err = xml.MarshalIndent(metricsDto, " ", " ")
}

func getValueKingString(kind metrics.ValueKind) string {
	switch kind {
	case metrics.KindBad:
		return "KindBad"
	case metrics.KindUint64:
		return "KindUint64"
	case metrics.KindFloat64:
		return "KindFloat64"
	case metrics.KindFloat64Histogram:
		return "KindFloat64Histogram"
	default:
		return "KindUnknown"
	}
}

func medianBucket(h *metrics.Float64Histogram) float64 {
	total := uint64(0)
	for _, count := range h.Counts {
		total += count
	}
	thresh := total / 2
	total = 0
	for i, count := range h.Counts {
		total += count
		if total >= thresh {
			return h.Buckets[i]
		}
	}
	panic("should not happen")
}
