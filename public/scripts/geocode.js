const routeForm = document.getElementById('route-form');
const routeSubmit = document.getElementById('route-submit');
const fromInput = document.getElementById('from');
const fromDatalist = document.getElementById('from-datalist');
const toInput = document.getElementById('to');
const toDatalist = document.getElementById('to-datalist');
const algorithmSelect = document.getElementById('algorithm-select');

let from = '';
let to = '';
let fromTimeoutID;
let toTimeoutID;

async function fetchJSON(url, def) {
    response = await fetch(url);
    if (!response.ok) {
        return def;
    }

    return await response.json();
}

function addAddressOption(datalist, value, text) {
    const option = document.createElement('option');
    option.value = value;
    option.innerText = text;
    datalist.appendChild(option);
}

async function getGeocodeResults(input) {
    const geocodes = await fetchJSON(window.location.origin + `/api/v1/geocode?address=${input.value}`, []);
    const geocodeResults = new Map()

    geocodes.forEach(s => {
        const address = /^(?<street>[^0-9,]+)\s*(?<number>\d+)?[,\s]*(?<postal>\d{4})?\s*(?<city>\D+)?$/.exec(input.value);
        if (!address) {
            return;
        }

        const number = address.groups.number?.trim();
        if (number === undefined) {
            geocodeResults.set(`${s.fullname}, ${s.l_pc} ${s.l_axon}`, s.f_jnctid);
        } else {
            const min = Math.min(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
            const max = Math.max(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);

            for (i = min; i <= max; i++) {
                if (`${i}`.startsWith(`${number}`)) {
                    geocodeResults.set(`${s.fullname} ${i}, ${s.l_pc} ${s.l_axon}`, s.f_jnctid);
                }
            }
        }
    });

    return geocodeResults;
}

async function getPlaceResults(input) {
    const places = await fetchJSON(window.location.origin + `/api/v1/places?input=${input.value}`, []);
    const placeResults = new Map()

    places.forEach(s => {
        placeResults.set(`${s.fullname}, ${s.l_pc} ${s.l_axon}`, s.f_jnctid);
    });

    return placeResults;
}

async function addAddressSuggestions(input, datalist) {
    const [geocodeResults, placeResults] = await Promise.all([
        getGeocodeResults(input),
        getPlaceResults(input)
    ]);
    datalist.innerHTML = ''

    const maxResults = 10;
    const maxPerSource = 5;

    let geocodeArray = Array.from(geocodeResults.entries()).slice(0, maxPerSource);
    let placeArray = Array.from(placeResults.entries()).slice(0, maxPerSource);

    if (geocodeArray.length < maxPerSource) {
        const remaining = maxPerSource - geocodeArray.length;
        placeArray = Array.from(placeResults.entries()).slice(0, maxPerSource + remaining);
    } else if (placeArray.length < maxPerSource) {
        const remaining = maxPerSource - placeArray.length;
        geocodeArray = Array.from(geocodeResults.entries()).slice(0, maxPerSource + remaining);
    }

    const combinedResults = [...placeArray, ...geocodeArray].slice(0, maxResults);

    combinedResults.forEach(([key, value]) => {
        addAddressOption(datalist, key, value);
    });

    if (combinedResults.length === 0) {
        addAddressOption(datalist, `${input.value} − No results found`, '');
    }
}

fromInput.addEventListener('input', function () {
    clearTimeout(fromTimeoutID);
    fromTimeoutID = setTimeout(addAddressSuggestions, 500, fromInput, fromDatalist);
});

toInput.addEventListener('input', function () {
    clearTimeout(toTimeoutID);
    toTimeoutID = setTimeout(addAddressSuggestions, 500, toInput, toDatalist);
});

routeForm.addEventListener('submit', async function (e) {
    e.preventDefault();

    let params = new URLSearchParams(window.location.search);

    Array.from(fromDatalist.children).forEach(option => {
        if (option.value == fromInput.value) {
            from = option.innerText;
            params.set("from", from);
        }
    });

    Array.from(toDatalist.children).forEach(option => {
        if (option.value == toInput.value) {
            to = option.innerText;
            params.set("to", to)
        }
    });

    params.set("algorithm", algorithmSelect.value);

    window.history.pushState('', '', window.location.origin + '?' + params.toString())

    routeSubmit.innerText = 'Loading';
    routeSubmit.toggleAttribute('disabled');

    geojson = await fetchJSON(window.location.origin + `/api/v1/route/${params.get("algorithm")}?from=${from}&to=${to}`, {});

    routeSubmit.innerText = 'Route';
    routeSubmit.toggleAttribute('disabled');

    loadJSON(map, geojson, from, to, params.get("algorithm"));
});
