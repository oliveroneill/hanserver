# hancollector
A set of collectors that store images in a unified format

## Usage
Run `hancollector` with the first argument being the path of a json file that
specifies collector configuration (see `default_config.json` in the parent
directory as an example). This will start retrieving images from a set of
regions defined in the database. If no regions are in the database,
`hancollector` will create a region based in San Francisco. Regions can be
viewed in the `regions` collections in the `han` mongo database.
These regions are set based on requests to `hanhttpserver` but could also be
set manually.
NOTE: `hanhttpserver` starts this itself, so this does not need to be run at
the same time.

## Development
Adding new image sources requires implementing the `ImageCollector` interface
found in `collectors/collector.go`.
This may include implementing `config.CollectorConfiguration` so that API keys
etc. can be stored in a unified place.

The collection process is based on *regions*, which are commonly queried areas.
These regions are periodically queried to retrieve the latest images. Regions
should be chosen based on recent queries as an attempt to avoid relying on
Instagram query latency, for example.