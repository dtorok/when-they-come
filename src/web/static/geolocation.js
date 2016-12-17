function locateUser(onSuccess, onFailure) {
    var success = function(position) {
        onSuccess(position.coords.latitude, position.coords.longitude)
    };

    var failure = function(error) {
        onFailure("position.status", "Unable to retrieve your location: " + error.code + " - " + error.message)
    };

    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success, failure);
    } else {
        onFailure("Geolocation is not supported by your browser")
    }
}
