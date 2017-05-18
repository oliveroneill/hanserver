/**
 * Puts a date into a relative measure for human readability
 */
function timeSince(date) {
  var seconds = Math.floor((new Date() - date) / 1000);
  var interval = Math.floor(seconds / 31536000);
  if (interval > 1) {
	return interval + " years";
  }
  interval = Math.floor(seconds / 2592000);
  if (interval > 1) {
	return interval + " months";
  }
  interval = Math.floor(seconds / 86400);
  if (interval > 1) {
	return interval + " days";
  }
  interval = Math.floor(seconds / 3600);
  if (interval > 1) {
	return interval + " hours";
  }
  interval = Math.floor(seconds / 60);
  if (interval > 1) {
	return interval + " minutes";
  }
  return Math.floor(seconds) + " seconds";
}

/**
 * Converts a distance in meters into a more readable distance
 * This is capped at 100km
 */
function convertDistance(meters) {
  if (meters > 1000 && meters < 100000) {
	var distance = meters / 1000;
	return distance+"km";
  } else if (meters >= 100000) {
	return ">100km";
  }
  return meters+"m";
}