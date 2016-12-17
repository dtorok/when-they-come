function Almafa() {
    console.log("hahaha")
}


function Map(div, getStops, getVehicles) {
    var map
    var currPos;
    var visibleStops = {};
    var init = function() {
        LONDON = this.coord2Pos(51.515037, -0.072384)
        currPos = LONDON;

        initGoogleMaps()
        refreshStops()
    };

    var initGoogleMaps = function() {
        map = new google.maps.Map(div, {
            zoom: 16,
            center: currPos
        });

        map.addListener("dragend", refreshStops)
        map.addListener("zoom_changed", refreshStops)
    };

    var refreshStops = function() {
        var center = map.getCenter()

        getStops(center, showStops)
    };

    var showStops = function(data) {
        for (i in data) {
            var d = data[i]

            if (d.id in visibleStops) {
                continue
            }

            var stop = new Stop(d, map, getVehicles)
            visibleStops[d.id] = stop
        }
    };

    // ===
    // public meethods
    // ===
    this.setPosition = function(position) {

    };

    this.coord2Pos = function(lat, lng) {
        return {lat: lat, lng: lng}
    }

    init()
}

function Stop(data, map, getVehicles) {
    var id = data.id
    var name = data.name
    var distance = data.distance
    var lat = data.lat
    var lon = data.lon

    marker = new google.maps.Marker({
        position: {
            lat: lat,
            lng: lon
        },
        map: map,
        icon: {
            path: google.maps.SymbolPath.CIRCLE,
            fillColor: "white",
            fillOpacity: 1,
            strokeColor: "blue",
            strokeWeight: 2,
            scale: 18
        }
    });

    var stopInfo = new StopInfo(name, distance);

    var info = new google.maps.InfoWindow({
        content: stopInfo.getContent()
    });

    info.addListener('domready', function() {
        getVehicles(id, function(data) {
            stopInfo.updateVehicles(data)
        })
    });

    marker.addListener("click", function() {
        info.open(map, this)
    });
}

function StopInfo(name, distance) {
    var content = $('<div/>').addClass('stopinfo')

    var title = $('<div/>').addClass("name").text(name + " (" + distance + ")")
    content.append(title)

    var vehicles = $('<div/>').addClass("vehicles").text("Loading...")
    content.append(vehicles)

    this.updateVehicles = function(data) {
        var strong = function(t) { return $('<strong/>').text(t)}
        vehicles.text("")

        var ul = $('<ul/>')

        data.sort(function(a, b) { return a.arrivesIn - b.arrivesIn })

        for (i in data) {
            var d = data[i]
            var li = $('<li/>')

            var lineName = strong(d.lineName)
            var towards = strong(d.towards)
            var arrivesAt = strong(d.arrivesAt)
            var arrivesIn = strong(Math.round(d.arrivesIn / 60))

            li
                .append(lineName)
                .append(" line towards ")
                .append(towards)
                .append(" arrivesin ")
                .append(arrivesIn)
                .append(" mins")

            ul.append(li)
        }

        vehicles.append(ul)
        console.log(data)
    };

    this.getContent = function() {
        return content[0]
    };

}