let loadedJSONLayer = null;
let currentLocationLayer = null;
let markers = []

function addOpenStreetMaps(map) {
    L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {maxZoom: 19}).addTo(map);
}

function onEachFeature(feature, layer) {
    if (feature.properties && feature.properties.angle_diff) {
        layer.bindPopup(getDirectionFromAngle(feature.properties.angle_diff) + ' on ' + feature.properties.street_name + " (" + feature.properties.angle_diff + "°)");
    }
}

function onLocationFound(map, e) {
    if (currentLocationLayer) {
        map.removeLayer(currentLocationLayer);
    }

    currentLocationLayer = L.circle(e.latlng, e.accuracy).addTo(map);

    if (loadedJSONLayer === null) {
        map.setView(e.latlng, 15);
    }
}

function registerLocationEvents(map) {
    map.on('locationfound', function (e) {
        onLocationFound(map, e);
    })
    map.locate({watch: true});
}

function addMarker(map, lat, long, start) {
    const icon = L.icon({
        iconUrl: `icons/${start ? 'start' : 'finish'}-marker.svg`,
        iconSize: start ? [35, 35] : [30, 30],
        iconAnchor: start ? [30 / 2, 35] : [2, 31],
    });

    markers.push(
        L.marker([lat, long], {icon}).addTo(map)
    )
}

function removeLayers(map) {
    if (currentLocationLayer) {
        map.removeLayer(currentLocationLayer);
    }

    if (loadedJSONLayer) {
        map.removeLayer(loadedJSONLayer);
    }

    markers.forEach(marker => {
        map.removeLayer(marker);
    });

    markers = []
}

function addMarkers(map, json) {
    if ((json.features ?? []).length > 0) {
        const [long, lat] = json.features[0].geometry.coordinates[0][0]
        addMarker(map, lat, long, true)

        const lastCoordinatesGroup = json.features[json.features.length - 1].geometry.coordinates
        const lastCoordinates = lastCoordinatesGroup[lastCoordinatesGroup.length - 1]
        const [long2, lat2] = lastCoordinates[lastCoordinates.length - 1]
        addMarker(map, lat2, long2, false)
    }
}

function loadJSON(map, json) {
    removeLayers(map)

    generateTurnByTurn(json)
    addMarkers(map, json)

    loadedJSONLayer = L.geoJSON(json, {
        onEachFeature: onEachFeature
    });
    loadedJSONLayer.addTo(map)

    zoomToData(map, loadedJSONLayer);
}

function zoomToData(map, data) {
    const bounds = data.getBounds();
    map.fitBounds(bounds);
}

var map = L.map('map')
addOpenStreetMaps(map);
registerLocationEvents(map)