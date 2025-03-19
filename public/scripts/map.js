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
}

function registerLocationEvents(map) {
    map.on('locationfound', function (e) {
        onLocationFound(map, e);
    })
    map.locate({watch: true});
}

function addMarker(map, lat, long, start) {
    const icon = L.icon({
        iconUrl: `icons/${start ? 'start' : 'finish'}.svg`,
        iconSize: [35, 35],
        iconAnchor: start ? [35 / 2, 35] : [9, 30],
    });

    markers.push(
        L.marker([lat, long], {icon}).addTo(map)
    )
}

function removeLayers(map) {
    if (currentLocationLayer) {
        map.removeLayer(currentLocationLayer);
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
    registerLocationEvents(map)
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
loadJSON(map,
// 20250319112128
// http://localhost:8080/api

    {
        "type": "FeatureCollection",
        "features": [
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4882036,
                                51.2308615
                            ],
                            [
                                5.4865927,
                                51.2299434
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404250,
                    "street_name": "Erkestraat",
                    "fow": 3,
                    "angle_diff": null,
                    "distance": 152,
                    "duration": 0.1824
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4865927,
                                51.2299434
                            ],
                            [
                                5.4865162,
                                51.2299764
                            ],
                            [
                                5.486401,
                                51.2300393
                            ],
                            [
                                5.4860378,
                                51.2302531
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 421416,
                    "street_name": "Broekkant",
                    "fow": 3,
                    "angle_diff": 56.257111779873,
                    "distance": 51.9,
                    "duration": 0.0623
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4860378,
                                51.2302531
                            ],
                            [
                                5.4856764,
                                51.2305447
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 414137,
                    "street_name": "Broekkant",
                    "fow": 3,
                    "angle_diff": 8.854248395539,
                    "distance": 41.1,
                    "duration": 0.0493
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4856764,
                                51.2305447
                            ],
                            [
                                5.4849458,
                                51.2311097
                            ],
                            [
                                5.4846777,
                                51.2312733
                            ],
                            [
                                5.484561,
                                51.2313262
                            ],
                            [
                                5.4844794,
                                51.2313386
                            ],
                            [
                                5.4843706,
                                51.2313571
                            ],
                            [
                                5.4840408,
                                51.2313613
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404394,
                    "street_name": "Broekkant",
                    "fow": 3,
                    "angle_diff": -2.786186663366,
                    "distance": 153.9,
                    "duration": 0.1848
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4840408,
                                51.2313613
                            ],
                            [
                                5.4840557,
                                51.2314868
                            ],
                            [
                                5.4840229,
                                51.2319172
                            ],
                            [
                                5.4839725,
                                51.232061
                            ],
                            [
                                5.4838786,
                                51.2322685
                            ],
                            [
                                5.4837654,
                                51.2324474
                            ],
                            [
                                5.4834207,
                                51.2328001
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404552,
                    "street_name": "Venderstraat",
                    "fow": 3,
                    "angle_diff": 85.192977911194,
                    "distance": 169.8,
                    "duration": 0.1358
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4834207,
                                51.2328001
                            ],
                            [
                                5.4831246,
                                51.233103
                            ],
                            [
                                5.4830277,
                                51.2332532
                            ],
                            [
                                5.4829306,
                                51.23346
                            ],
                            [
                                5.4828524,
                                51.2336919
                            ],
                            [
                                5.482726,
                                51.2341988
                            ],
                            [
                                5.482662,
                                51.234419
                            ],
                            [
                                5.4825517,
                                51.2351342
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404551,
                    "street_name": "Venderstraat",
                    "fow": 3,
                    "angle_diff": -0.196644692283,
                    "distance": 269.8,
                    "duration": 0.2159
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4825517,
                                51.2351342
                            ],
                            [
                                5.4825382,
                                51.2352191
                            ],
                            [
                                5.4819629,
                                51.2389
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404505,
                    "street_name": "Lillerheidestraat",
                    "fow": 3,
                    "angle_diff": 1.668784640007,
                    "distance": 421,
                    "duration": 0.3368
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4819629,
                                51.2389
                            ],
                            [
                                5.4817385,
                                51.2403371
                            ],
                            [
                                5.4817181,
                                51.2404673
                            ],
                            [
                                5.4817057,
                                51.2405608
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 616857,
                    "street_name": "Lillerheidestraat",
                    "fow": 3,
                    "angle_diff": 0.009081998584,
                    "distance": 185.6,
                    "duration": 0.1485
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4817057,
                                51.2405608
                            ],
                            [
                                5.4816013,
                                51.2412976
                            ],
                            [
                                5.481575,
                                51.2414825
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 617491,
                    "street_name": "Lillerheidestraat",
                    "fow": 3,
                    "angle_diff": 0.270637065490,
                    "distance": 102.9,
                    "duration": 0.0824
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.481575,
                                51.2414825
                            ],
                            [
                                5.4743667,
                                51.2400936
                            ],
                            [
                                5.4740292,
                                51.2400234
                            ],
                            [
                                5.473611,
                                51.2399582
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 406997,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -92.873066792551,
                    "distance": 581.5,
                    "duration": 0.4652
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.473611,
                                51.2399582
                            ],
                            [
                                5.4734614,
                                51.2399369
                            ],
                            [
                                5.4734402,
                                51.2399339
                            ],
                            [
                                5.4732784,
                                51.239911
                            ],
                            [
                                5.472788,
                                51.2398559
                            ],
                            [
                                5.47138,
                                51.2397265
                            ],
                            [
                                5.4702742,
                                51.2396243
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 616525,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 2.060769405087,
                    "distance": 236,
                    "duration": 0.1888
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4702742,
                                51.2396243
                            ],
                            [
                                5.4701017,
                                51.2395296
                            ],
                            [
                                5.47001,
                                51.239432
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404702,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -30.785405964758,
                    "distance": 28.6,
                    "duration": 0.0229
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.47001,
                                51.239432
                            ],
                            [
                                5.4700019,
                                51.2394234
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404701,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -10.665641937559,
                    "distance": 1.1,
                    "duration": 0.0009
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4700019,
                                51.2394234
                            ],
                            [
                                5.4699332,
                                51.2393503
                            ],
                            [
                                5.4697815,
                                51.2392669
                            ],
                            [
                                5.4695531,
                                51.2391728
                            ],
                            [
                                5.4693508,
                                51.2391105
                            ],
                            [
                                5.4689804,
                                51.2390337
                            ],
                            [
                                5.4679527,
                                51.238843
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404411,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 11.337413064027,
                    "distance": 160.4,
                    "duration": 0.1283
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4679527,
                                51.238843
                            ],
                            [
                                5.4678148,
                                51.2388174
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404410,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 0.314789869543,
                    "distance": 10,
                    "duration": 0.008
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4678148,
                                51.2388174
                            ],
                            [
                                5.4639351,
                                51.2380789
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404409,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -0.260539997890,
                    "distance": 283.1,
                    "duration": 0.2265
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4639351,
                                51.2380789
                            ],
                            [
                                5.4637166,
                                51.2380382
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404228,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 0.225746859033,
                    "distance": 15.9,
                    "duration": 0.0127
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4637166,
                                51.2380382
                            ],
                            [
                                5.4634246,
                                51.2379839
                            ],
                            [
                                5.4612683,
                                51.2375814
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 409739,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -0.017077609154,
                    "distance": 178.4,
                    "duration": 0.1427
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4612683,
                                51.2375814
                            ],
                            [
                                5.4605516,
                                51.2374448
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404545,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -0.222271494184,
                    "distance": 52.3,
                    "duration": 0.0418
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4605516,
                                51.2374448
                            ],
                            [
                                5.4545277,
                                51.2363057
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 427754,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 0.082899762438,
                    "distance": 439.4,
                    "duration": 0.3515
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4545277,
                                51.2363057
                            ],
                            [
                                5.4513514,
                                51.235732
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 408361,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 0.469689592527,
                    "distance": 230.8,
                    "duration": 0.1847
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4513514,
                                51.235732
                            ],
                            [
                                5.4510565,
                                51.2356741
                            ],
                            [
                                5.4510373,
                                51.2356703
                            ],
                            [
                                5.4509274,
                                51.2356488
                            ],
                            [
                                5.4508736,
                                51.2356391
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 621596,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -0.875033376730,
                    "distance": 34.9,
                    "duration": 0.0279
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4508736,
                                51.2356391
                            ],
                            [
                                5.4507828,
                                51.2356236
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 417494,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 1.103439982007,
                    "distance": 6.6,
                    "duration": 0.0053
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4507828,
                                51.2356236
                            ],
                            [
                                5.4506981,
                                51.2356077
                            ],
                            [
                                5.4502894,
                                51.2355438
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404056,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 0.500119191737,
                    "distance": 35.6,
                    "duration": 0.0285
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4502894,
                                51.2355438
                            ],
                            [
                                5.4501937,
                                51.2355263
                            ],
                            [
                                5.4500514,
                                51.2355003
                            ],
                            [
                                5.4496094,
                                51.2353968
                            ],
                            [
                                5.4494514,
                                51.2353459
                            ],
                            [
                                5.4491534,
                                51.2352414
                            ],
                            [
                                5.448822,
                                51.2351097
                            ],
                            [
                                5.4485088,
                                51.2349756
                            ],
                            [
                                5.4481178,
                                51.2347988
                            ],
                            [
                                5.4479961,
                                51.2347451
                            ],
                            [
                                5.4478446,
                                51.2346711
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 615212,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -1.170628135190,
                    "distance": 197.8,
                    "duration": 0.1583
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4478446,
                                51.2346711
                            ],
                            [
                                5.4476595,
                                51.2345631
                            ],
                            [
                                5.4475154,
                                51.2344868
                            ],
                            [
                                5.4475092,
                                51.2344825
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 407055,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -4.189461631179,
                    "distance": 31.5,
                    "duration": 0.0252
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4475092,
                                51.2344825
                            ],
                            [
                                5.4473855,
                                51.2343989
                            ],
                            [
                                5.4472366,
                                51.2342882
                            ],
                            [
                                5.446935,
                                51.2340486
                            ],
                            [
                                5.446566,
                                51.233748
                            ],
                            [
                                5.446003,
                                51.2332753
                            ],
                            [
                                5.4454497,
                                51.2327906
                            ],
                            [
                                5.4453526,
                                51.2327057
                            ],
                            [
                                5.4452094,
                                51.2325799
                            ],
                            [
                                5.4451584,
                                51.2325352
                            ],
                            [
                                5.4450692,
                                51.2324569
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 616740,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -7.277062647921,
                    "distance": 282.7,
                    "duration": 0.2261
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4450692,
                                51.2324569
                            ],
                            [
                                5.4432693,
                                51.2308447
                            ],
                            [
                                5.4431332,
                                51.230727
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 407869,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": -0.521106891857,
                    "distance": 235.2,
                    "duration": 0.2352
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4431332,
                                51.230727
                            ],
                            [
                                5.442916,
                                51.2305969
                            ],
                            [
                                5.4428427,
                                51.2305476
                            ],
                            [
                                5.4426433,
                                51.2304623
                            ],
                            [
                                5.4424923,
                                51.2303855
                            ],
                            [
                                5.4424248,
                                51.2303685
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 615810,
                    "street_name": "Hamonterweg",
                    "fow": 3,
                    "angle_diff": 10.084519852819,
                    "distance": 64,
                    "duration": 0.064
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4424248,
                                51.2303685
                            ],
                            [
                                5.4422569,
                                51.2302726
                            ],
                            [
                                5.4419713,
                                51.2301301
                            ],
                            [
                                5.4416691,
                                51.230005
                            ],
                            [
                                5.4413932,
                                51.2299052
                            ],
                            [
                                5.4413309,
                                51.229885
                            ],
                            [
                                5.441177,
                                51.2298426
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403962,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": -4.497021019337,
                    "distance": 105.5,
                    "duration": 0.1055
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.441177,
                                51.2298426
                            ],
                            [
                                5.4410709,
                                51.2298227
                            ],
                            [
                                5.4408478,
                                51.229783
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403958,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": 5.886243767330,
                    "distance": 23.9,
                    "duration": 0.0239
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4408478,
                                51.229783
                            ],
                            [
                                5.4406682,
                                51.2297565
                            ],
                            [
                                5.4404566,
                                51.2297384
                            ],
                            [
                                5.4401462,
                                51.229707
                            ],
                            [
                                5.4398687,
                                51.2296823
                            ],
                            [
                                5.4394998,
                                51.2296526
                            ],
                            [
                                5.4392224,
                                51.2296343
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403936,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": 3.757851097759,
                    "distance": 114.8,
                    "duration": 0.1148
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4392224,
                                51.2296343
                            ],
                            [
                                5.4387224,
                                51.2295901
                            ],
                            [
                                5.4386403,
                                51.2295828
                            ],
                            [
                                5.438375,
                                51.2295624
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 616649,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": -0.808454794053,
                    "distance": 59.7,
                    "duration": 0.0597
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.438375,
                                51.2295624
                            ],
                            [
                                5.4375003,
                                51.2294928
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 404100,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": 0.009396323711,
                    "distance": 61.6,
                    "duration": 0.0616
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4375003,
                                51.2294928
                            ],
                            [
                                5.4374056,
                                51.2294833
                            ],
                            [
                                5.4372479,
                                51.2294618
                            ],
                            [
                                5.43715,
                                51.2294433
                            ],
                            [
                                5.4370834,
                                51.2294354
                            ],
                            [
                                5.4369845,
                                51.2294132
                            ],
                            [
                                5.4369223,
                                51.2293999
                            ],
                            [
                                5.4368197,
                                51.2293708
                            ],
                            [
                                5.4366114,
                                51.2292957
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 410534,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": -2.452604436443,
                    "distance": 66.4,
                    "duration": 0.0664
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4366114,
                                51.2292957
                            ],
                            [
                                5.4365514,
                                51.2292684
                            ],
                            [
                                5.4364892,
                                51.2292425
                            ],
                            [
                                5.4363274,
                                51.2291658
                            ],
                            [
                                5.4361711,
                                51.2290862
                            ],
                            [
                                5.4359609,
                                51.2289945
                            ],
                            [
                                5.4358099,
                                51.2289319
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 416820,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": -4.997106102128,
                    "distance": 69.1,
                    "duration": 0.0691
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4358099,
                                51.2289319
                            ],
                            [
                                5.435501,
                                51.2288207
                            ],
                            [
                                5.4353055,
                                51.2287576
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403917,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": 4.068495884613,
                    "distance": 40.2,
                    "duration": 0.0402
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4353055,
                                51.2287576
                            ],
                            [
                                5.4351887,
                                51.2287208
                            ],
                            [
                                5.4347886,
                                51.2286035
                            ],
                            [
                                5.4341913,
                                51.2284189
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 406571,
                    "street_name": "Heerstraat",
                    "fow": 3,
                    "angle_diff": 2.462503689549,
                    "distance": 86.5,
                    "duration": 0.0865
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4341913,
                                51.2284189
                            ],
                            [
                                5.4341246,
                                51.2283969
                            ],
                            [
                                5.4340932,
                                51.2283865
                            ],
                            [
                                5.4339108,
                                51.2283264
                            ],
                            [
                                5.4333958,
                                51.2281393
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 621283,
                    "street_name": "Koning Albertlaan",
                    "fow": 3,
                    "angle_diff": -1.436728783510,
                    "distance": 63.7,
                    "duration": 0.0637
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4333958,
                                51.2281393
                            ],
                            [
                                5.4333012,
                                51.2281603
                            ],
                            [
                                5.4332109,
                                51.2281998
                            ],
                            [
                                5.4327149,
                                51.2283796
                            ],
                            [
                                5.4325392,
                                51.2284371
                            ],
                            [
                                5.4324579,
                                51.2284638
                            ],
                            [
                                5.4324311,
                                51.2284755
                            ],
                            [
                                5.4323096,
                                51.228529
                            ],
                            [
                                5.4319699,
                                51.2286843
                            ],
                            [
                                5.4318984,
                                51.2287121
                            ],
                            [
                                5.4317992,
                                51.2287359
                            ],
                            [
                                5.431684,
                                51.2287558
                            ],
                            [
                                5.4315745,
                                51.2287767
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 621282,
                    "street_name": "Kerkstraat",
                    "fow": 3,
                    "angle_diff": 37.635649091577,
                    "distance": 146.6,
                    "duration": 0.2514
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4315745,
                                51.2287767
                            ],
                            [
                                5.4312695,
                                51.228836
                            ],
                            [
                                5.4310949,
                                51.2288601
                            ],
                            [
                                5.4308075,
                                51.2288804
                            ],
                            [
                                5.4306107,
                                51.228884
                            ],
                            [
                                5.4304504,
                                51.2288851
                            ],
                            [
                                5.4303481,
                                51.2288766
                            ],
                            [
                                5.4302787,
                                51.2288656
                            ],
                            [
                                5.4301906,
                                51.2288387
                            ],
                            [
                                5.4301009,
                                51.2288028
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403866,
                    "street_name": "Kerkstraat",
                    "fow": 3,
                    "angle_diff": -0.426583928217,
                    "distance": 106.4,
                    "duration": 0.1825
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4301009,
                                51.2288028
                            ],
                            [
                                5.4300998,
                                51.2288485
                            ],
                            [
                                5.4300927,
                                51.2288669
                            ],
                            [
                                5.4300796,
                                51.2289426
                            ],
                            [
                                5.4300592,
                                51.2290228
                            ],
                            [
                                5.4299938,
                                51.2292203
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403862,
                    "street_name": "Broesveldstraat",
                    "fow": 3,
                    "angle_diff": 102.163543619905,
                    "distance": 47.1,
                    "duration": 0.0629
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4299938,
                                51.2292203
                            ],
                            [
                                5.4298223,
                                51.2295948
                            ],
                            [
                                5.4294996,
                                51.2302634
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403854,
                    "street_name": "Broesveldstraat",
                    "fow": 3,
                    "angle_diff": -8.181264566006,
                    "distance": 121.1,
                    "duration": 0.1614
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4294996,
                                51.2302634
                            ],
                            [
                                5.4293368,
                                51.2305119
                            ],
                            [
                                5.4291659,
                                51.2307283
                            ],
                            [
                                5.4289661,
                                51.2309587
                            ],
                            [
                                5.4288772,
                                51.2310738
                            ],
                            [
                                5.4287775,
                                51.2311959
                            ],
                            [
                                5.428639,
                                51.2313892
                            ],
                            [
                                5.4285468,
                                51.2315296
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 403846,
                    "street_name": "Broesveldstraat",
                    "fow": 3,
                    "angle_diff": -10.319698550870,
                    "distance": 155.9,
                    "duration": 0.2079
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4285468,
                                51.2315296
                            ],
                            [
                                5.428476,
                                51.2316376
                            ],
                            [
                                5.4283988,
                                51.2317939
                            ],
                            [
                                5.4283677,
                                51.2318997
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 410505,
                    "street_name": "Broesveldstraat",
                    "fow": 3,
                    "angle_diff": 5.410131510683,
                    "distance": 43.2,
                    "duration": 0.0576
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4283677,
                                51.2318997
                            ],
                            [
                                5.4283645,
                                51.2319109
                            ],
                            [
                                5.4283446,
                                51.2320118
                            ],
                            [
                                5.428327,
                                51.2320878
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 410506,
                    "street_name": "Broesveldstraat",
                    "fow": 3,
                    "angle_diff": 10.806737604490,
                    "distance": 21.1,
                    "duration": 0.0282
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.428327,
                                51.2320878
                            ],
                            [
                                5.42816,
                                51.2320772
                            ],
                            [
                                5.4281104,
                                51.2320747
                            ],
                            [
                                5.4280593,
                                51.2320775
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 410504,
                    "street_name": "Sint-Martensstraat",
                    "fow": 3,
                    "angle_diff": -81.492422841775,
                    "distance": 18.8,
                    "duration": 0.0564
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4280593,
                                51.2320775
                            ],
                            [
                                5.4280081,
                                51.2320982
                            ],
                            [
                                5.4279556,
                                51.2321259
                            ],
                            [
                                5.4279089,
                                51.2321565
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 408941,
                    "street_name": "Sint-Martensstraat",
                    "fow": 3,
                    "angle_diff": 24.84922612913,
                    "distance": 13.8,
                    "duration": 0.0413
                }
            },
            {
                "type": "Feature",
                "geometry": {
                    "type": "MultiLineString",
                    "coordinates": [
                        [
                            [
                                5.4279089,
                                51.2321565
                            ],
                            [
                                5.4278415,
                                51.2322007
                            ],
                            [
                                5.4277489,
                                51.2322551
                            ]
                        ]
                    ]
                },
                "properties": {
                    "gid": 422601,
                    "street_name": "Sint-Martensstraat",
                    "fow": 3,
                    "angle_diff": 1.200629382367,
                    "distance": 15.7,
                    "duration": 0.047
                }
            }
        ]
    }
)