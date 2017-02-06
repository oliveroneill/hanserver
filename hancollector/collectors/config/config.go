package config

// Each collector should specify a configuration
// This should determine whether a collector is enabled, as well as other
// details
//
// Unfortunately this isn't as effective as I would have liked, because
// you cannot cast a generic interface to a specific one. So there's no
// way to refer to the implemented configuration's additional fields from
// within the collector. However this still enforces the use of a configuration
//
// These configurations can be specified through json or through code by
// implementing the configuration struct however needed. This wasn't necessary
// for this task

// CollectorConfiguration is the base configuration
type CollectorConfiguration interface {
    IsEnabled() bool
}