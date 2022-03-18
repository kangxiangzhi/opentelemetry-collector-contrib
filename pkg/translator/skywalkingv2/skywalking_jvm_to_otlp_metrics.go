// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skywalkingv2

import (
	"strings"

	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
	v3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentV3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// SkywalkingToOtlpMetrics
type SkywalkingToOtlpMetrics struct {
	Logger *zap.Logger
}

//JVMToMetrics
func (s SkywalkingToOtlpMetrics) JVMToMetrics(authentication string, jvmMetricCollector *agentV3.JVMMetricCollection) pdata.Metrics {
	metricData := pdata.NewMetrics()
	rms := metricData.ResourceMetrics()
	rms.EnsureCapacity(len(jvmMetricCollector.Metrics))
	index := 0
	serviceInstance := s.getIp(jvmMetricCollector.ServiceInstance)
	for _, metric := range jvmMetricCollector.Metrics {
		s.jvmMetricToResourceMetrics(metric, rms.AppendEmpty(), jvmMetricCollector.Service)
		s.jvmToInternalResource(jvmMetricCollector.Service, serviceInstance,
			authentication, rms.At(index).Resource())
		index++
	}
	return metricData
}

func (s SkywalkingToOtlpMetrics) jvmToInternalResource(service string, serviceInstance string,
	authentication string, dest pdata.Resource) {
	attrs := dest.Attributes()
	attrs.EnsureCapacity(2)
	attrs.InsertString(attributeServiceName, service)
	attrs.InsertString(attributeServiceInstance, serviceInstance)
}

func (s SkywalkingToOtlpMetrics) jvmMetricToResourceMetrics(jvmMetric *agentV3.JVMMetric, out pdata.ResourceMetrics, service string) {
	ilms := out.InstrumentationLibraryMetrics()
	ilms.EnsureCapacity(JvmMetricLength)

	ilm := ilms.AppendEmpty()
	metrics := ilm.Metrics()
	metrics.EnsureCapacity(JvmMetricLength)

	metric := metrics.AppendEmpty()
	metric.SetDataType(pdata.MetricDataTypeGauge)
	metricDps := metric.Gauge().DataPoints()
	metricDps.EnsureCapacity(s.calculateDpsLenByMetrics(jvmMetric))

	index := 0
	//gc metric to otlp metric
	index = s.gcMetricToMetrics(jvmMetric.Time, jvmMetric.Gc, metricDps, index)
	//memory pool metric to otlp metric
	index = s.memoryPoolMetricToMetrics(jvmMetric.Time, jvmMetric.MemoryPool, metricDps, index)
	//memory area metric to optl metric
	index = s.memoryMetricToMetrics(jvmMetric.Time, jvmMetric.Memory, metricDps, index)
	//thread metric to otlp metric
	index = s.threadMetricToMetrics(jvmMetric.Time, jvmMetric.Thread, metricDps, index)
	//cpu metric to otpl metric
	index = s.cpuMetricToMetrics(jvmMetric.Time, jvmMetric.Cpu, metricDps, index)
}

// 根据传入参数metrics，计算出dps需要的长度
func (s SkywalkingToOtlpMetrics) calculateDpsLenByMetrics(jvmMetric *agentV3.JVMMetric) int {
	return len(jvmMetric.Gc)*GcMetricNumber +
		len(jvmMetric.MemoryPool)*MemoryPoolMetricNumber +
		len(jvmMetric.Memory)*MemoryMetricNumber +
		ThreadMetricNumber +
		CpuMetricNumber
}

func (s SkywalkingToOtlpMetrics) gcMetricToMetrics(time int64, gcList []*agentV3.GC,
	dps pdata.NumberDataPointSlice, index int) int {
	//metrics为 gc.count与gc.time
	for _, gc := range gcList {
		gcName := ""
		if gc.Phase == agentV3.GCPhase_NEW {
			gcName = New
		} else {
			gcName = Old
		}
		s.fillDoubleDataPoint(time, GCLabelKey, gcName+Count,
			float64(gc.Count), dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, GCLabelKey, gcName+Time,
			float64(gc.Time), dps.AppendEmpty())
		index++
	}
	return index
}

func (s SkywalkingToOtlpMetrics) memoryPoolMetricToMetrics(time int64, memoryPools []*agentV3.MemoryPool,
	dps pdata.NumberDataPointSlice, index int) int {
	for _, memoryPool := range memoryPools {
		poolTypeName := strings.ToLower(agentV3.PoolType_name[int32(memoryPool.Type)])
		s.fillDoubleDataPoint(time, MemoryPoolKey, poolTypeName+Commit,
			float64(memoryPool.Committed)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryPoolKey, poolTypeName+Init,
			float64(memoryPool.Init)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryPoolKey, poolTypeName+Used,
			float64(memoryPool.Used)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryPoolKey, poolTypeName+Max,
			float64(memoryPool.Max)/Megabytes, dps.AppendEmpty())
		index++
	}
	return index
}

func (s SkywalkingToOtlpMetrics) memoryMetricToMetrics(time int64, memoryList []*agentV3.Memory,
	dps pdata.NumberDataPointSlice, index int) int {
	for _, memory := range memoryList {
		heapName := ""
		if memory.IsHeap {
			heapName = HeapArea
		} else {
			heapName = NonHeapArea
		}
		s.fillDoubleDataPoint(time, MemoryKey, heapName+Commit,
			float64(memory.Committed)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryKey, heapName+Init,
			float64(memory.Init)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryKey, heapName+Used,
			float64(memory.Used)/Megabytes, dps.AppendEmpty())
		index++
		s.fillDoubleDataPoint(time, MemoryKey, heapName+Max,
			float64(memory.Max)/Megabytes, dps.AppendEmpty())
		index++
	}
	return index
}

func (s SkywalkingToOtlpMetrics) threadMetricToMetrics(time int64, thread *agentV3.Thread,
	dps pdata.NumberDataPointSlice, index int) int {

	if thread == nil { //空指针判断
		return index
	}
	s.fillDoubleDataPoint(time, ThreadKey, LiveCount,
		float64(thread.LiveCount), dps.AppendEmpty())
	index++
	s.fillDoubleDataPoint(time, ThreadKey, DaemonCount,
		float64(thread.DaemonCount), dps.AppendEmpty())
	index++
	s.fillDoubleDataPoint(time, ThreadKey, PeakCount,
		float64(thread.PeakCount), dps.AppendEmpty())
	index++

	return index
}

func (s SkywalkingToOtlpMetrics) cpuMetricToMetrics(time int64, cpu *v3.CPU,
	dps pdata.NumberDataPointSlice, index int) int {
	if cpu == nil { //空指针判断
		return index
	}
	s.fillDoubleDataPoint(time, CpuKey, UsagePercent, cpu.UsagePercent, dps.AppendEmpty())
	index++
	return index
}

func (s SkywalkingToOtlpMetrics) fillDoubleDataPoint(time int64, labelKey string, labelValue string,
	value float64, dp pdata.NumberDataPoint) {

	labelsMap := dp.Attributes()
	labelsMap.InsertString(labelKey, labelValue)
	dp.SetDoubleVal(value)
	dp.SetTimestamp(pdata.Timestamp(time))
}

func (s SkywalkingToOtlpMetrics) getIp(serviceInstance string) string {
	pos := strings.Index(serviceInstance, "@")
	if pos != -1 {
		return serviceInstance[pos+1:]
	}
	return serviceInstance
}
