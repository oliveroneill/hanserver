# han server
A server that stores images from arbitrary sources and can return them based on location and recency. Han is short for 'here and now'.

## Dependencies
* [Docker](https://docs.docker.com/engine/installation/)
* [docker-compose](https://docs.docker.com/compose/install/)

## Usage
There are two separate components that are part of han:
* Image population (hancollector) - this will use a variety of collector implementations to retrieve images from different sources and store them in a unified format.
* Web server (hanhttpserver) - this retrieves the images for a client using HTTP.

Both of these can be started by simply calling `docker-compose build && docker-compose up` from the base directory. Alternatively they can be started individually by calling that same command from within either `hanhttpserver` or `hancollector`.
Before calling `docker-compose up` you will need to open `hancollector/collectors/config/instagram.go` and set the values labelled `SETME`.
If you don't want to use instagram, just set `Enabled` to `false` and the `APIKey` to an empty string.
The `hanapi` directory contains common classes between the two.

There's an additional README in both `hanhttpserver` and `hancollector` that discusses their development.

## Testing
All tests can be run using the command `go test ./...`, as you can see there are only two sets of tests at the moment, this will be worked on in the future.

## TODO
This is a list of features or issues I'd like to work on in the future.
* Deployment - the two Dockerfiles contain the same dependencies and should use the same base image
* Regioning - to make this project scalable, locations are broken up into regions, these regions are used to avoid populating the whole world with images. This should aim to keep the database size down by choosing the most recently used locations.
* Configuration - I wanted to keep the configuration between each collector separate but also not bog the user down with having to store a set of configuration files within some common path. For now, that meant storing the configuration for each collector within static code (see `hancollector/collectors/config`). Files (or a single master file) is definitely a better way to go and I'll move to this when I have the chance.
* Cleaning up images - images that have been deleted from their original source need to be taken down, there needs to be a neat way of doing this without periodically downloading the images to check the response code