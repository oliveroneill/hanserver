// han server address
const serverAddress = "localhost:8080";

window.onload = function () {
  'use strict';
  /**
   * Initialise viewer.js
   */
  var Viewer = window.Viewer;
  var console = window.console || { log: function () {} };
  var pictures = document.querySelector('.docs-pictures');
  var toggles = document.querySelector('.docs-toggles');
  var buttons = document.querySelector('.docs-buttons');
  var options = {
        movable: false,
        zoomable: false,
        rotatable: false,
        scalable: false,
        tooltip: false,
        // manually close the image when clicking outside of the image
        viewed: function () {
          document.getElementsByClassName('viewer-canvas')[0].onclick = function(e){
              if(e.target.className == 'viewer-canvas') {
                window.currentViewer.hide();
              }
          };
        }
  };

  var viewer = new Viewer(pictures, options);
  window.currentViewer = viewer;
};

// region size is used to draw blue circles on the map for each region
const regionSize = 5000;
/**
 * Arrays used to keep track of marker on the map
 */
var markersArray = [];
var imagesArray = [];
// marker for image location
var imageMarkerIcon = {
  url: 'http://maps.google.com/mapfiles/ms/icons/blue-dot.png',
};

function initMap() {
  // initialise position is Uluru, for no particular reason
  var lat = -25.363
  var lng = 131.044
  var uluru = {lat: lat, lng: lng};
  var map = new google.maps.Map(document.getElementById('map'), {
    zoom: 4,
    center: uluru
  });
  placeMarker(uluru, map);

  // on click a new marker will be placed
  google.maps.event.addListener(map, 'click', function(event) {
     placeMarker({lat:event.latLng.lat(), lng:event.latLng.lng()}, map);
  });
}
/**
 * Place a new marker on the map
 */
function placeMarker(location, map) {
  clearOverlays(markersArray);
    var marker = new google.maps.Marker({
        position: location,
        map: map
    });
    updateLoc(map, location.lat, location.lng);
    markersArray.push(marker);
}

/**
 * Delete all other markers
 */
function clearOverlays(markers) {
  for (var i = 0; i < markers.length; i++ ) {
    markers[i].setMap(null);
  }
  markers.length = 0;
}

function updateLoc(map, lat, lng) {
  // let the user know the results are loading
  document.getElementById("header").innerHTML = "Loading...";
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      var response = JSON.parse(this.responseText);
      var images = response.images;
      for (var i = 0; i < images.length; i++) {
        addImage(images[i], map);
      }
      // finished... this won't wait for regions
      document.getElementById("header").innerHTML = "Results";
    }
  };
  // clear the previously set images
  var imageList = document.getElementById('pics');
  imageList.innerHTML = '';
  // remove previous overlays
  clearOverlays(imagesArray);
  // query the server
  xhttp.open("GET", "http://"+serverAddress+"/api/image-search?lat="+lat+"&lng="+lng+"", true);
  xhttp.send();

  // we also need to get region data
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      var regions = JSON.parse(this.responseText);
      // will draw a blue circle for each region
      for (var i = 0; i < regions.length; i++) {
        var cityCircle = new google.maps.Circle({
          strokeColor: '#0000FF',
          strokeOpacity: 0.8,
          strokeWeight: 0,
          fillColor: '#0000FF',
          fillOpacity: 0.35,
          map: map,
          center: regions[i],
          radius: regionSize,
          clickable: false
        });
      }
    }
  };
  xhttp.open("GET", "http://"+serverAddress+"/api/get-regions", true);
  xhttp.send();
}

function addImage(image, map) {
  // create a list element in the pics list
  var ul = document.getElementById("pics");
  var li = document.createElement("li");
  var div = document.createElement("div");
  // needed for viewer.js
  div.setAttribute("class", "image");
  var img = document.createElement("img");
  img.setAttribute("src", image.url);
  img.setAttribute("alt", image.caption);
  // distance and recency info is set in an h2 element
  var h2 = document.createElement("h2");

  var distance = convertDistance(Math.round(image.distance));
  // set the subtitle
  h2.innerHTML = distance+" - "+timeSince(new Date(image.createdTime*1000))+" ago";
  // add the element
  div.appendChild(img);
  div.appendChild(h2);
  li.appendChild(div);
  ul.appendChild(li);

  // add the marker
  var marker = new google.maps.Marker({
    position: {lat: image.location.lat, lng: image.location.lng},
    map: map,
    // this will show a thumbnail when the marker is pressed
    html: "<img width=\"50\" height=\"50\" src=\""+image.url+"\"/>",
    icon: imageMarkerIcon
  });
  infowindow = new google.maps.InfoWindow({
    content: "..."
  });
  // show the image when it's pressed
  google.maps.event.addListener(marker, 'click', function () {
    infowindow.setContent(this.html);
    infowindow.open(map, this);
  });
  imagesArray.push(marker);
}
