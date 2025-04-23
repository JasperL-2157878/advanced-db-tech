let loadedJSONLayers = [];
let currentLocationLayer = null;
let markers = []
let locationSet = false;
let currentRoutePoints = null;

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

    if (loadedJSONLayers.length === 0 && !locationSet) {
        map.setView(e.latlng, 15);
        locationSet = true;
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

    for (let i = 0; i < loadedJSONLayers.length; i++) {
        map.removeLayer(loadedJSONLayers[i]);
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

function getAlgorithmColor(value) {
    const colorMap = new Map([
        ["alg/dijkstra", "#1f77b4"],
        ["alg/astar", "#ff7f0e"],
        ["alg/bddijkstra", "#2ca02c"],
        ["alg/bdastar", "#d62728"],
        ["opt/none", "#9467bd"],
        ["opt/tnr", "#8c564b"],
        ["opt/ch", "#e377c2"],
        ["opt/chtnr", "#7f7f7f"]
    ]);

    return colorMap.get(value) || "#000000";
}

function loadJSON(map, json, from, to, algorithm) {
    const newRoutePoints = from + "_" + to

    if (currentRoutePoints !== newRoutePoints) {
        removeLayers(map)
        addMarkers(map, json)
        currentRoutePoints = newRoutePoints
    }

    generateTurnByTurn(json)

    const loadedJSONLayer = L.geoJSON(json, {
        onEachFeature: onEachFeature,
        style: {
            "color": getAlgorithmColor(algorithm),
        }
    });
    loadedJSONLayer.addTo(map)

    zoomToData(map, loadedJSONLayer);

    loadedJSONLayers.push(loadedJSONLayer)
}

function zoomToData(map, data) {
    const bounds = data.getBounds();
    map.fitBounds(bounds);
}

var map = L.map('map')
addOpenStreetMaps(map);
map.setView([50.9250399, 5.389659], 15);
registerLocationEvents(map)