package handler

import (
	"fmt"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/memstorage"
	"log"
	"net/http"
	"strconv"
)

type MetricsListRequestHandler struct {
	gaugeMetricRepository   *memstorage.GaugeMetricRepository
	counterMetricRepository *memstorage.CounterMetricRepository
}

func NewMetricsListRequestHandler(gaugeMetricRepository *memstorage.GaugeMetricRepository, counterMetricRepository *memstorage.CounterMetricRepository) *MetricsListRequestHandler {
	return &MetricsListRequestHandler{gaugeMetricRepository: gaugeMetricRepository, counterMetricRepository: counterMetricRepository}
}

func (obj *MetricsListRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: MetricsListRequestHandler")
	body := `
<!doctype html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<title>Metrics list</title>
	</head>
	<body>
		<div class="container">
			<h1>Metrics list</h1>
			<table class="table">`
	body += `<tr><td colspan="2"><b>Gauge metrics</b></td></tr>`
	for _, v := range obj.gaugeMetricRepository.FindAll() {
		body += fmt.Sprintf("<tr><td>%s</td><td>%v</td><tr>", v.Name(), strconv.FormatFloat(v.Value(), 'f', -1, 64))
	}
	body += `<tr><td colspan="2"><b>Counter metrics</b></td></tr>`
	for _, v := range obj.counterMetricRepository.FindAll() {
		body += fmt.Sprintf("<tr><td>%s</td><td>%v</td><tr>", v.Name(), strconv.FormatInt(v.Value(), 10))
	}
	body += `
			</table>
		</div>
	</body>
</html>`
	response.Write([]byte(body))
}
