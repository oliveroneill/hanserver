package config

import (
	"time"
)

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
	IsEnabled()		  bool
	GetCollectorName()   string
	// The frequency in seconds at which the collector should be updated for
	// all regions
	GetUpdateFrequency() time.Duration
	// Limit on query per GetQueryWindow seconds
	GetQueryLimit()	  int
	// in seconds
	GetQueryWindow()	 int64
}

// CollectorConfig is a type used for CollectorConfiguration interface
type CollectorConfig struct {
	CollectorConfiguration
	Enabled		 bool
	CollectorName   string
	UpdateFrequency time.Duration
	QueryLimit	  int
	QueryWindow	 int64
}

// IsEnabled if this collector should be used
func (c CollectorConfig) IsEnabled() bool {
	return c.Enabled
}

// GetCollectorName returns the name of the collector for logging purposes
func (c CollectorConfig) GetCollectorName() string {
	return c.CollectorName
}

// GetUpdateFrequency returns the frequency which the collector should be
// updated
func (c CollectorConfig) GetUpdateFrequency() time.Duration {
	return c.UpdateFrequency
}

// GetQueryLimit returns the limit of queries to be made per query window
func (c CollectorConfig) GetQueryLimit() int {
	return c.QueryLimit
}

// GetQueryWindow returns the window to be used with GetQueryLimit
func (c CollectorConfig) GetQueryWindow() int64 {
	return c.QueryWindow
}
