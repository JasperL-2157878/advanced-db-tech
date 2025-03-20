const form = document.querySelector('#form');
const fromInput = document.querySelector("#from");
const fromDatalist = document.querySelector("#from-datalist");
const toInput = document.querySelector("#to");
const toDatalist = document.querySelector("#to-datalist");

var prevFromDatalistOptions;
var prevToDatalistOptions;

async function fetchJSON(url, def) {
    response = await fetch(url);
    if (!response.ok) {
        return def;
    }

    return await response.json();
}

function parseAddress(address) {
    return /^(?<street_name>[^0-9,]+)\s*(?<street_number>\d+)?[,\s]*(?<postal_code>\d{4})?\s*(?<city_name>\D+)?$/.exec(address);
}

function registerAutocomplete(input, datalist) {
    input.addEventListener('keyup', async function (e) {
        suggestions = await fetchJSON(`http://localhost:8080/api/v1/geocode?address=${input.value}`, []);
        
        if (datalist.id.includes('from')) {
            prevFromDatalistOptions = Array.from(datalist.children);
        } else {
            prevToDatalistOptions = Array.from(datalist.children);
        }

        datalist.innerHTML = ''
        suggestions.forEach(s => {
            min = Math.min(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
            max = Math.max(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
            
            address = parseAddress(input.value);
            if (!address) {
                return;
            }

            street = address.groups.street;
            number = address.groups.number;
            postal = address.groups.postal;
            city = address.groups.city;

            const option = document.createElement('option');

            if (number == undefined) {
                option.value = `${s.fullname}, ${s.l_pc} ${s.l_axon}`;
            } else if (min <= number && number <= max) {
                option.value = `${s.fullname} ${number}, ${s.l_pc} ${s.l_axon}`;
            }

            option.innerText = s.f_jnctid;
            datalist.appendChild(option);
        });
    });
}

registerAutocomplete(fromInput, fromDatalist);
registerAutocomplete(toInput, toDatalist);

form.addEventListener('submit', async function (e) {
    e.preventDefault();

    fromId = '';
    toId = '';

    prevFromDatalistOptions.forEach(option => {
        if (option.value == fromInput.value) {
            fromId = option.innerText;
        }
    });

    prevToDatalistOptions.forEach(option => {
        if (option.value == toInput.value) {
            toId = option.innerText;
        }
    });

    geojson = await fetchJSON(`http://localhost:8080/api/v1/route?from=${fromId}&to=${toId}`, {});
    console.log(`http://localhost:8080/api/v1/route?from=${fromId}&to=${toId}`);
    loadJSON(map, geojson);
});
