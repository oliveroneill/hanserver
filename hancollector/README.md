# hancollector
A set of collectors that store images in a unified format

## Usage
Run `hancollector` to start retrieving images from a set of regions defined in the database.
These regions are set based on requests to `hanhttpserver` but could also be set manually.

## Development
Adding new image sources requires implementing the `ImageCollector` interface found in `collectors/collector.go`.
This may include implementing `config.CollectorConfiguration` so that API keys etc. can be stored in a unified place. See the TODO section in the base README for ideas on how to do this better. Each configuration has an `Enabled` field to easily enable and disable the collectors you want.

The collection process is based on *regions*, which are commonly queried areas. These regions are periodically queried to retrieve the latest images. Regions should be chosen based on recent queries as an attempt to avoid relying on Instagram query latency, for example.