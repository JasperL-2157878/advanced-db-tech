function getDirectionFromAngle(angle) {
    switch (true) {
        case angle < -35:
            return "left";
        case angle > 35:
            return "right";
        default:
            return "straight";
    }
}

function addStep(element, icon, text) {
    const stepElement = document.createElement('div');
    stepElement.className = "route-step";
    element.appendChild(stepElement);

    addIcon(stepElement, icon);

    const stepContentElement = document.createElement('div');
    stepContentElement.className = "step-content";
    stepElement.appendChild(stepContentElement);

    const stepInstructionElement = document.createElement('div');
    stepInstructionElement.className = "step-instruction";
    stepInstructionElement.innerText = text;
    stepContentElement.appendChild(stepInstructionElement);

    return stepContentElement
}

function setRouteSummary(distance, duration) {
    const element = document.getElementById('route-summary');
    element.innerText = `${(distance / 1000).toFixed(2)} km • ${Math.ceil(duration)} min`;
}

function addStepDetails(element, distance, duration) {
    const stepDetailsElement = document.createElement('div');
    stepDetailsElement.className = "step-details";
    stepDetailsElement.innerText = `${(distance / 1000).toFixed(2)} km • ${Math.ceil(duration)} min`;
    element.appendChild(stepDetailsElement);
}

function addIcon(element, icon) {
    const iconContainer = document.createElement('div');
    iconContainer.className="step-icon-containers"
    element.appendChild(iconContainer);

    const stepIconElement = document.createElement('div')
    stepIconElement.className = "step-icon";
    iconContainer.appendChild(stepIconElement);

    const iconElement = document.createElement('img');
    iconElement.className = "icon"
    iconElement.src = `/icons/${icon}.svg`;
    stepIconElement.appendChild(iconElement);
}

function addRoundabout(element, exit, street) {
    return addStep(
        element,
        "roundabout-" + Math.min(3, Math.max(1, exit)),
        `Take the ${exit} exit onto ${street}`
    )
}

function generateTurnByTurn(json) {
    const element = document.getElementById('route-steps');
    while (element.firstChild) {
        element.removeChild(element.firstChild);
    }

    let isStart = true;
    let onRoundabout = false;

    let totalMeters = 0;
    let totalMinutes = 0;
    let accumulatedMeters = 0;
    let accumulatedMinutes = 0;

    let currentStreet = null;
    let currentExit = null;
    let currentElement = null;

    json.features.forEach((feature, index) => {
        const streetChanged = feature.properties.street_name !== currentStreet;

        if (streetChanged) {
            if (currentStreet && currentElement) {
                addStepDetails(currentElement, accumulatedMeters, accumulatedMinutes);
            }

            currentStreet = feature.properties.street_name;

            accumulatedMeters = feature.properties.distance;
            accumulatedMinutes = feature.properties.duration;

            totalMeters += accumulatedMeters;
            totalMinutes += accumulatedMinutes;

            if (isStart) {
                isStart = false;
                currentElement = addStep(element, "start", `Start on ${currentStreet}`);
            }
        } else {
            accumulatedMeters += feature.properties.distance;
            accumulatedMinutes += feature.properties.duration;

            totalMeters += feature.properties.distance;
            totalMinutes += feature.properties.duration;
        }

        if (feature.properties.fow === 4) {
            onRoundabout = true;
            currentExit = (currentExit ?? 0) + 1;
        } else {
            if (onRoundabout) {
                currentElement = addRoundabout(element, currentExit, currentStreet);

                currentExit = null;
                onRoundabout = false;
            } else if (streetChanged) {
                switch (getDirectionFromAngle(feature.properties.angle_diff)) {
                    case 'left':
                        currentElement = addStep(element, "left", `Turn left onto ${currentStreet}`);
                        break;
                    case 'right':
                        currentElement = addStep(element, "right", `Turn right onto ${currentStreet}`);
                        break;
                    case 'straight':
                        currentElement = addStep(element, "straight", `Go straight onto ${currentStreet}`);
                        break;
                }
            }
        }
    });

    addStep(element, "finish", `Finish on ${currentStreet}`);
    addStepDetails(currentElement, accumulatedMeters, accumulatedMinutes);

    setRouteSummary(totalMeters, totalMinutes);
}