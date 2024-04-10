# Metrics server

## Supported API:

- `GET /metrics/all` - return all metric names from metrics package without values in xml format
- `POST metrics/filter` - return metric by filter from metrics package with values in xml format

## Examples:
Req `POST metrics/filter`:
```
<requested_metrics>
    <name>/cpu/classes/total:cpu-seconds</name>
    <name>/cpu/classes/user:cpu-seconds</name>
    <name>/sched/goroutines:goroutines</name>
    <name>/sched/latencies:seconds</name>
    <name>/sync/mutex/wait/total:seconds</name>
    <name>/cpu/classes/gc/pause:cpu-seconds</name>
    <name>/cpu/classes/gc/total:cpu-seconds</name>
    <name>/gc/heap/allocs:bytes</name>
    <name>/memory/classes/heap/free:bytes</name>
    <name>/memory/classes/heap/stacks:bytes</name>
</requested_metrics>
```
Res:
```
 <metrics>
  <metric>
   <name>/cpu/classes/total:cpu-seconds</name>
   <description>Estimated total available CPU time for user Go code or the Go runtime, as defined by GOMAXPROCS. In other words, GOMAXPROCS integrated over the wall-clock duration this process has been executing for. This metric is an overestimate, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics. Sum of all metrics in /cpu/classes.</description>
   <kind>KindFloat64</kind>
   <cumulative>true</cumulative>
   <value type="float64">0</value>
  </metric>
  <metric>
   <name>/cpu/classes/user:cpu-seconds</name>
   <description>Estimated total CPU time spent running user Go code. This may also include some small amount of time spent in the Go runtime. This metric is an overestimate, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics.</description>
   <kind>KindFloat64</kind>
   <cumulative>true</cumulative>
   <value type="float64">0</value>
  </metric>
  <metric>
   <name>/sched/goroutines:goroutines</name>
   <description>Count of live goroutines.</description>
   <kind>KindUint64</kind>
   <cumulative>false</cumulative>
   <value type="uint64">6</value>
  </metric>
  <metric>
   <name>/sched/latencies:seconds</name>
   <description>Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.</description>
   <kind>KindFloat64Histogram</kind>
   <cumulative>false</cumulative>
   <value type="histogram-median">0</value>
  </metric>
  <metric>
   <name>/sync/mutex/wait/total:seconds</name>
   <description>Approximate cumulative time goroutines have spent blocked on a sync.Mutex or sync.RWMutex. This metric is useful for identifying global changes in lock contention. Collect a mutex or block profile using the runtime/pprof package for more detailed contention data.</description>
   <kind>KindFloat64</kind>
   <cumulative>true</cumulative>
   <value type="float64">0</value>
  </metric>
  <metric>
   <name>/cpu/classes/gc/pause:cpu-seconds</name>
   <description>Estimated total CPU time spent with the application paused by the GC. Even if only one thread is running during the pause, this is computed as GOMAXPROCS times the pause latency because nothing else can be executing. This is the exact sum of samples in /gc/pause:seconds if each sample is multiplied by GOMAXPROCS at the time it is taken. This metric is an overestimate, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics.</description>
   <kind>KindFloat64</kind>
   <cumulative>true</cumulative>
   <value type="float64">0</value>
  </metric>
  <metric>
   <name>/cpu/classes/gc/total:cpu-seconds</name>
   <description>Estimated total CPU time spent performing GC tasks. This metric is an overestimate, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics. Sum of all metrics in /cpu/classes/gc.</description>
   <kind>KindFloat64</kind>
   <cumulative>true</cumulative>
   <value type="float64">0</value>
  </metric>
  <metric>
   <name>/gc/heap/allocs:bytes</name>
   <description>Cumulative sum of memory allocated to the heap by the application.</description>
   <kind>KindUint64</kind>
   <cumulative>true</cumulative>
   <value type="uint64">388968</value>
  </metric>
  <metric>
   <name>/memory/classes/heap/free:bytes</name>
   <description>Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime&#39;s estimate of free address space that is backed by physical memory.</description>
   <kind>KindUint64</kind>
   <cumulative>false</cumulative>
   <value type="uint64">8192</value>
  </metric>
  <metric>
   <name>/memory/classes/heap/stacks:bytes</name>
   <description>Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.</description>
   <kind>KindUint64</kind>
   <cumulative>false</cumulative>
   <value type="uint64">196608</value>
  </metric>
 </metrics>
```