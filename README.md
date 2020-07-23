# golang-statsd

golang-statsd is a convenience wrapper around datadog-go/statsd functions

## Getting started
#### Import golang-statsd:
`import github.com/wistia/golang-statsd`
#### Configure golang-statsd:

golang-statsd is disabled by default. To enable it do something like the following:

```
statsd_host := "69.163.225.80"  
namespace := "dinos"
env := "prod"
component := "comics"  

statsd.Configure(statsd_host, namespace, env, component)  
```
`component` will be added as a tag under the key 'component'  

#### Disable golang-statsd:
`statsd.Disable()`  

Disabling golang-statsd initializes golang-statsd with a no-op writer; this is useful if you don't want to report metrics in a particular environment.  

#### Using golang-statsd:
Use golang-statsd as you would datadog-go/statsd but save yourself the trouble of dealing with error handling.  
See https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go