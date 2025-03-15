# Overview
Type `Workload` allows specification of a benchmark or test to run, allowing specification of metrics to be collected. A `Workload` can be executed against any connection that provides `io.Reader` (stdin) and `io.WriterCloser` (stdout/err) interfaces. There are useful predefined Workload entries.

Type `Metric` is used by `Workload` to describe how to parse (what string format to look for) and interpret a metric (larger or smaller better); there is predefined `Runtime` metric (smaller better).

Type `Tunable` allows specification of a filesystem (eg `/proc`or `/sys`) parameter that can be adjusted so as to minimise or maximise one of more `Workload` `Metrics`.

Type `Session` allows sweeping `Tunables`, measuring covariance in `Workload` `Metric`s.
