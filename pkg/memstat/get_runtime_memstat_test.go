package memstat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRuntimeMemstat(t *testing.T) {
	memStat := GetRuntimeMemstat()

	assert.Contains(t, memStat, "Alloc")
	assert.Contains(t, memStat, "BuckHashSys")
	assert.Contains(t, memStat, "Frees")
	assert.Contains(t, memStat, "GCCPUFraction")
	assert.Contains(t, memStat, "GCSys")
	assert.Contains(t, memStat, "HeapAlloc")
	assert.Contains(t, memStat, "HeapIdle")
	assert.Contains(t, memStat, "HeapInuse")
	assert.Contains(t, memStat, "HeapObjects")
	assert.Contains(t, memStat, "HeapReleased")
	assert.Contains(t, memStat, "HeapSys")
	assert.Contains(t, memStat, "LastGC")
	assert.Contains(t, memStat, "Lookups")
	assert.Contains(t, memStat, "MCacheInuse")
	assert.Contains(t, memStat, "MCacheSys")
	assert.Contains(t, memStat, "MSpanInuse")
	assert.Contains(t, memStat, "MSpanSys")
	assert.Contains(t, memStat, "Mallocs")
	assert.Contains(t, memStat, "NextGC")
	assert.Contains(t, memStat, "NumForcedGC")
	assert.Contains(t, memStat, "NumGC")
	assert.Contains(t, memStat, "OtherSys")
	assert.Contains(t, memStat, "PauseTotalNs")
	assert.Contains(t, memStat, "StackInuse")
	assert.Contains(t, memStat, "StackSys")
	assert.Contains(t, memStat, "Sys")
}
