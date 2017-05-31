# han server
[![Build Status](https://travis-ci.org/oliveroneill/hanserver.svg?branch=master)](https://travis-ci.org/oliveroneill/hanserver)

A server that stores images from arbitrary sources and can return them based
on location and recency. Han is short for 'here and now'.

## Dependencies
* [Docker](https://docs.docker.com/engine/installation/)
* [docker-compose](https://docs.docker.com/compose/install/)

## Usage
There are two separate components that are part of han:
* Image population (hancollector) - this will use a variety of collector
implementations to retrieve images from different sources and store them
in a unified format.
* Web server (hanhttpserver) - this retrieves the images for a client using HTTP.

Both of these are started through `hanhttpserver`, which can be started by
simply calling `docker-compose build && docker-compose up` from the base
directory. Configuration is required for image population to work, see the
Configuration section below.
Alternatively they can be started individually by calling that same
command from within either `hanhttpserver` or `hancollector`.
Note that `hanhttpserver` automatically starts `hancollector` within the same
process, this is used to keep track of API calls between the server and the
collector. `hanhttpserver` can be started without `hancollector` by using the
`--no-collection` option.

### Slack logging
Errors can be logged through Slack by passing in the `--slacktoken` argument
into `hanhttpserver`. This is logged to the "hanserver" channel but can be
changed in `hanapi/reporting/reporting.go`.

### Configuration
Before calling `docker-compose up` you will need to copy `default_config.json`
and set the required fields to configure the collectors. Copying this json
and calling it `config.json` is recommended since the `gitignore` includes
this file.
`hanhttpserver` will throw an error if no collectors are enabled.
If you don't want to use the implemented collectors, just set `enabled` to
`false`. You must then implement your own collector, see
`hancollector/README.md` for more info.

The `hanapi` directory contains common classes between these two components.

There's an additional README in both `hanhttpserver` and `hancollector` that
discusses their development.

## Testing
All tests can be run using the command `go test ./...`, as you can see there
are only two sets of tests at the moment, this will be worked on in the future.

## TODO
This is a list of features or issues I'd like to work on in the future.
* Deployment - the two Dockerfiles contain the same dependencies and should use
the same base image
* Regioning - to make this project scalable, locations are broken up into
regions, these regions are used to avoid populating the whole world with images.
This should aim to keep the database size down by choosing the most recently
used locations.
* Cleaning up images - images that have been deleted from their original
source need to be taken down, there needs to be a neat way of doing this
without periodically downloading the images to check the response code