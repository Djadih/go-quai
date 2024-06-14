package quaistats

import (
	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/core/types"
	"github.com/dominant-strategies/go-quai/metrics_config"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Metrics
	headerGauges  = metrics_config.NewGaugeVec("headerchain", "HeaderChain metrics")
	headerPrime   = headerGauges.WithLabelValues("prime")
	headerRegions = make([]prometheus.Gauge, common.MaxRegions)
	headerZones   = make([][]prometheus.Gauge, common.MaxRegions)
)

func init() {
	for region := 0; region < common.MaxRegions; region++ {
		loc := common.Location{byte(region)}
		headerRegions[region] = headerGauges.WithLabelValues(loc.Name())

		headerZones[region] = make([]prometheus.Gauge, common.MaxZones)
		for zone := 0; zone < common.MaxZones; zone++ {
			// headerZones = append(headerZones, headerGauges.WithLabelValues("region", "zone"))
			loc, err := common.NewLocation(region, zone)
			if err != nil {
				panic(err)
			}
			headerZones[region][zone] = headerGauges.WithLabelValues(loc.Name())
		}
	}
}

func updateHeaderMetrics(block *types.WorkObject) {
	// Update prime
	headerPrime.Set(float64(block.NumberU64(common.PRIME_CTX)))
	headerRegions[block.Location().Region()].Set(float64(block.NumberU64(common.REGION_CTX)))
	headerZones[block.Location().Region()][block.Location().Zone()].Set(float64(block.NumberU64(common.ZONE_CTX)))
}
