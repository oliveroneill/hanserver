# hanhttpserver
A web server for serving up han images, this also includes a simple demo webpage.
This also starts `hancollector` in the background for collecting images.

## Usage
Call `hanhttpserver` with `-nocollection` to disable image collection.
When starting this program there is a demo webpage that can be used to click on different places on the map and observe the feed from that location.
Just open `demo/index.html` in the browser. The server location needs to be set via `demo/src/js/main.js` line 2.
Individual regions are marked in blue, these are areas that will be automatically populated and queried based on the amount of queries made in that location. Tapping will change your user location to that area and reload the feed.

Be careful when using this demo page in production, as making new image queries changes `hancollector`'s priorities of where to populate.