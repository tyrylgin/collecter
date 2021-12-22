package metrics

import (
	"runtime"
)

var (
	runtimeMetrics struct {
		Alloc         Gauge
		BuckHashSys   Gauge
		Frees         Gauge
		GCCPUFraction Gauge
		GCSys         Gauge
		HeapAlloc     Gauge
		HeapIdle      Gauge
		HeapInuse     Gauge
		HeapObjects   Gauge
		HeapReleased  Gauge
		HeapSys       Gauge
		LastGC        Gauge
		Lookups       Gauge
		MCacheInuse   Gauge
		MCacheSys     Gauge
		MSpanInuse    Gauge
		MSpanSys      Gauge
		Mallocs       Gauge
		NextGC        Gauge
		NumForcedGC   Gauge
		NumGC         Gauge
		OtherSys      Gauge
		PauseTotalNs  Gauge
		StackInuse    Gauge
		StackSys      Gauge
		Sys           Gauge
	}
	isRuntimeMetricRegistered bool
)

// RegisterRuntimeMetrics register runtime.MemStats. Should be fired once
func RegisterRuntimeMetrics() {
	if isRuntimeMetricRegistered {
		return
	}

	runtimeMetrics.Alloc = NewGauge()
	runtimeMetrics.BuckHashSys = NewGauge()
	runtimeMetrics.Frees = NewGauge()
	runtimeMetrics.GCCPUFraction = NewGauge()
	runtimeMetrics.GCSys = NewGauge()
	runtimeMetrics.HeapAlloc = NewGauge()
	runtimeMetrics.HeapIdle = NewGauge()
	runtimeMetrics.HeapInuse = NewGauge()
	runtimeMetrics.HeapObjects = NewGauge()
	runtimeMetrics.HeapReleased = NewGauge()
	runtimeMetrics.HeapSys = NewGauge()
	runtimeMetrics.LastGC = NewGauge()
	runtimeMetrics.Lookups = NewGauge()
	runtimeMetrics.MCacheInuse = NewGauge()
	runtimeMetrics.MCacheSys = NewGauge()
	runtimeMetrics.MSpanInuse = NewGauge()
	runtimeMetrics.MSpanSys = NewGauge()
	runtimeMetrics.Mallocs = NewGauge()
	runtimeMetrics.NextGC = NewGauge()
	runtimeMetrics.NumForcedGC = NewGauge()
	runtimeMetrics.NumGC = NewGauge()
	runtimeMetrics.OtherSys = NewGauge()
	runtimeMetrics.PauseTotalNs = NewGauge()
	runtimeMetrics.StackInuse = NewGauge()
	runtimeMetrics.StackSys = NewGauge()
	runtimeMetrics.Sys = NewGauge()

	Register("Alloc", runtimeMetrics.Alloc)
	Register("BuckHashSys", runtimeMetrics.BuckHashSys)
	Register("Frees", runtimeMetrics.Frees)
	Register("GCCPUFraction", runtimeMetrics.GCCPUFraction)
	Register("GCSys", runtimeMetrics.GCSys)
	Register("HeapAlloc", runtimeMetrics.HeapAlloc)
	Register("HeapIdle", runtimeMetrics.HeapIdle)
	Register("HeapInuse", runtimeMetrics.HeapInuse)
	Register("HeapObjects", runtimeMetrics.HeapObjects)
	Register("HeapReleased", runtimeMetrics.HeapReleased)
	Register("HeapSys", runtimeMetrics.HeapSys)
	Register("LastGC", runtimeMetrics.LastGC)
	Register("Lookups", runtimeMetrics.Lookups)
	Register("MCacheInuse", runtimeMetrics.MCacheInuse)
	Register("MCacheSys", runtimeMetrics.MCacheSys)
	Register("MSpanInuse", runtimeMetrics.MSpanInuse)
	Register("MSpanSys", runtimeMetrics.MSpanSys)
	Register("Mallocs", runtimeMetrics.Mallocs)
	Register("NextGC", runtimeMetrics.NextGC)
	Register("NumForcedGC", runtimeMetrics.NumForcedGC)
	Register("NumGC", runtimeMetrics.NumGC)
	Register("OtherSys", runtimeMetrics.OtherSys)
	Register("PauseTotalNs", runtimeMetrics.PauseTotalNs)
	Register("StackInuse", runtimeMetrics.StackInuse)
	Register("StackSys", runtimeMetrics.StackSys)
	Register("Sys", runtimeMetrics.Sys)
}

func SnapshotRuntimeMetrics() {
	if !isRuntimeMetricRegistered {
		RegisterRuntimeMetrics()
	}

	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	runtimeMetrics.Alloc.Set(float64(memStat.Alloc))
	runtimeMetrics.BuckHashSys.Set(float64(memStat.BuckHashSys))
	runtimeMetrics.Frees.Set(float64(memStat.Frees))
	runtimeMetrics.GCCPUFraction.Set(memStat.GCCPUFraction)
	runtimeMetrics.GCSys.Set(float64(memStat.GCSys))
	runtimeMetrics.HeapAlloc.Set(float64(memStat.HeapAlloc))
	runtimeMetrics.HeapIdle.Set(float64(memStat.HeapIdle))
	runtimeMetrics.HeapInuse.Set(float64(memStat.HeapInuse))
	runtimeMetrics.HeapObjects.Set(float64(memStat.HeapObjects))
	runtimeMetrics.HeapReleased.Set(float64(memStat.HeapReleased))
	runtimeMetrics.HeapSys.Set(float64(memStat.HeapSys))
	runtimeMetrics.LastGC.Set(float64(memStat.LastGC))
	runtimeMetrics.Lookups.Set(float64(memStat.Lookups))
	runtimeMetrics.MCacheInuse.Set(float64(memStat.MCacheInuse))
	runtimeMetrics.MCacheSys.Set(float64(memStat.MCacheSys))
	runtimeMetrics.MSpanInuse.Set(float64(memStat.MSpanInuse))
	runtimeMetrics.MSpanSys.Set(float64(memStat.MSpanSys))
	runtimeMetrics.Mallocs.Set(float64(memStat.Mallocs))
	runtimeMetrics.NextGC.Set(float64(memStat.NextGC))
	runtimeMetrics.NumForcedGC.Set(float64(memStat.NumForcedGC))
	runtimeMetrics.NumGC.Set(float64(memStat.NumGC))
	runtimeMetrics.OtherSys.Set(float64(memStat.OtherSys))
	runtimeMetrics.PauseTotalNs.Set(float64(memStat.PauseTotalNs))
	runtimeMetrics.StackInuse.Set(float64(memStat.StackInuse))
	runtimeMetrics.StackSys.Set(float64(memStat.StackSys))
	runtimeMetrics.Sys.Set(float64(memStat.Sys))
}
