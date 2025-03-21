const routeFrom = document.getElementById('route-form');
const routeSubmit = document.getElementById('route-submit');
const fromInput = document.getElementById('from');
const fromDatalist = document.getElementById('from-datalist');
const toInput = document.getElementById('to');
const toDatalist = document.getElementById('to-datalist');

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

async function addAddressSuggestions(input, datalist) {
    suggestions = await fetchJSON(window.location.origin + `/api/v1/geocode?address=${input.value}`, []);
    datalist.innerHTML = ''

    suggestions.forEach(s => {
        min = Math.min(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
        max = Math.max(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
        
        address = /^(?<street>[^0-9,]+)\s*(?<number>\d+)?[,\s]*(?<postal>\d{4})?\s*(?<city>\D+)?$/.exec(input.value);
        if (!address) {
            return;
        }

        street = address.groups.street.trim();
        number = address.groups.number?.trim();
        postal = address.groups.postal?.trim();
        city = address.groups.city?.trim();

        if (number == undefined) {
            addAddressOption(datalist, `${s.fullname}, ${s.l_pc} ${s.l_axon}`, s.f_jnctid);
        } else {
            for (i = min; i <= max; i++) {
                if (`${i}`.startsWith(`${number}`)) {
                    addAddressOption(datalist, `${s.fullname} ${i}, ${s.l_pc} ${s.l_axon}`, s.f_jnctid);
                }
            }
        }
    });

    if (suggestions.length == 0 && input.value.length != 0) {
        addAddressOption(datalist, `${input.value} − No results found`, '');
    }
}

fromInput.addEventListener('input', function() {
    clearTimeout(fromTimeoutID);
    fromTimeoutID = setTimeout(addAddressSuggestions, 500, fromInput, fromDatalist);
});

toInput.addEventListener('input', function() {
    clearTimeout(toTimeoutID);
    toTimeoutID = setTimeout(addAddressSuggestions, 500, toInput, toDatalist);
});

routeFrom.addEventListener('submit', async function (e) {
    e.preventDefault();

    fromId = '';
    toId = '';

    Array.from(fromDatalist.children).forEach(option => {
        if (option.value == fromInput.value) {
            fromId = option.innerText;
        }
    });

    Array.from(toDatalist.children).forEach(option => {
        if (option.value == toInput.value) {
            toId = option.innerText;
        }
    });

    routeSubmit.innerText = 'Loading';
    routeSubmit.toggleAttribute('disabled');

    geojson = await fetchJSON(window.location.origin + `/api/v1/route?from=${fromId}&to=${toId}`, {});

    routeSubmit.innerText = 'Route';
    routeSubmit.toggleAttribute('disabled');

    loadJSON(map, geojson);
});
