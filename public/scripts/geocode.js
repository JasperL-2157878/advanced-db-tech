const fromInput = document.querySelector("#from");
const toInput = document.querySelector("#to");
const fromDatalist = document.querySelector("#from-datalist");
const toDatalist = document.querySelector("#to-datalist");

async function fetchSuggestions(input) {
    response = await fetch(`http://localhost:8080/api/v1/geocode?address=${input}`, { 
        method: 'GET'
    });
    
    if (!response.ok) {
        return [];
    }

    return await response.json();
}

function formatSuggestion(s, input) {
    if (input.match(/\d/)) {

    }


}

function registerAutocomplete(input, datalist) {
    input.addEventListener('keyup', async () => {
        suggestions = await fetchSuggestions(input.value);
        datalist.innerHTML = '';

        suggestions.forEach(s => {
            min = Math.min(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);
            max = Math.max(s.l_f_add, s.l_t_add, s.r_f_add, s.r_t_add);

            range = ` ${min}-${max}`
            if (min < 0 && max < 0) {
                range = ''
            } else if (min < 0 && max > 0) {
                range = ` ${max}`
            } else if (min > 0 && max < 0) {
                range = ` ${min}`
            }
            
            const option = document.createElement('option');
            option.innerText = `${s.fullname}${range}, ${s.l_pc} ${s.l_axon}`
            datalist.appendChild(option);
        });
    });
}

registerAutocomplete(fromInput, fromDatalist);
registerAutocomplete(toInput, toDatalist);
