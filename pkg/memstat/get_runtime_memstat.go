package memstat

import "runtime"

func GetRuntimeMemstat() map[string]float64 {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	return map[string]float64{
		"Alloc":         float64(memStat.Alloc),
		"BuckHashSys":   float64(memStat.BuckHashSys),
		"Frees":         float64(memStat.Frees),
		"GCCPUFraction": memStat.GCCPUFraction,
		"GCSys":         float64(memStat.GCSys),
		"HeapAlloc":     float64(memStat.HeapAlloc),
		"HeapIdle":      float64(memStat.HeapIdle),
		"HeapInuse":     float64(memStat.HeapInuse),
		"HeapObjects":   float64(memStat.HeapObjects),
		"HeapReleased":  float64(memStat.HeapReleased),
		"HeapSys":       float64(memStat.HeapSys),
		"LastGC":        float64(memStat.LastGC),
		"Lookups":       float64(memStat.Lookups),
		"MCacheInuse":   float64(memStat.MCacheInuse),
		"MCacheSys":     float64(memStat.MCacheSys),
		"MSpanInuse":    float64(memStat.MSpanInuse),
		"MSpanSys":      float64(memStat.MSpanSys),
		"Mallocs":       float64(memStat.Mallocs),
		"NextGC":        float64(memStat.NextGC),
		"NumForcedGC":   float64(memStat.NumForcedGC),
		"NumGC":         float64(memStat.NumGC),
		"OtherSys":      float64(memStat.OtherSys),
		"PauseTotalNs":  float64(memStat.PauseTotalNs),
		"StackInuse":    float64(memStat.StackInuse),
		"StackSys":      float64(memStat.StackSys),
		"Sys":           float64(memStat.Sys),
	}
}
