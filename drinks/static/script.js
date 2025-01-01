let dataOptions = {year: 'numeric', month: 'long', day: 'numeric'};

function downloadHistory() {
    let drinkLog = new Map(JSON.parse(localStorage.log));
    let data = JSON.stringify(Object.fromEntries(drinkLog), null, 2)
    const downloadLink = document.createElement('a');
    downloadLink.href = URL.createObjectURL(new Blob([data], {
        type: "application/json"
    }));
    downloadLink.setAttribute("download", "history.json");
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
}

function createCheckbox(drinkKey) {
    const checkbox = document.createElement("input")
    checkbox.setAttribute("type", "checkbox");
    checkbox.setAttribute("id", drinkKey);
    checkbox.setAttribute("class", "check-box");
    checkbox.setAttribute("inert", "true");
    return checkbox;
}

function showScore(id) {
    document.getElementById(id+".rating").style.visibility = "visible";
}

function hideScore(id) {
    document.getElementById(id+".rating").style.visibility = "hidden";
}

function markAsFinished(checkbox, drinkName, drinkDescription) {
    let date = new Date()
    if (localStorage.log) {
        let drinkLog = new Map(JSON.parse(localStorage.log));
        drinkLog.set(checkbox.id, {name: drinkName, date: date, description: drinkDescription});
        localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
    } else {
        let drinkLog = new Map();
        drinkLog.set(checkbox.id, {name: drinkName, date: date, description: drinkDescription});
        localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
    }
    checkbox.checked = true;
    addItemToHistory(checkbox.id, drinkName, drinkDescription, date);
    showScore(checkbox.id);
    console.log("finished " + checkbox.id);
}

function markToDo(checkbox) {
    if (localStorage.log) {
        let drinkLog = new Map(JSON.parse(localStorage.log));
        drinkLog.delete(checkbox.id);
        localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
    } else {
        let drinkLog = new Map();
        localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
    }
    checkbox.checked = false;
    removeItemFromHistory(checkbox.id);
    hideScore(checkbox.id);
    console.log("removed " + checkbox.id);
}

function addCheckbox(item) {
    let drink = item.firstElementChild;
    if (drink.firstChild) {
        let drinkName = drink.firstChild.textContent;
        let drinkDescription = item.lastElementChild.innerText;
        let drinkKey = drinkName.replace(/\s/g, "").replace(/\.+$/, "").toLowerCase();
        drink.insertBefore(createCheckbox(drinkKey), drink.firstChild);

        item.addEventListener("click", (event) => {
            if (event.target.tagName !== "SPAN") {
                let checkbox = event.target.parentElement.firstElementChild.firstElementChild;
                let checked = checkbox.checked;
                if (!checked) {
                    markAsFinished(checkbox, drinkName, drinkDescription);
                } else {
                    markToDo(checkbox);
                }
            }
        });
    }
}

function CreateHistoryBtn() {
    let historyLnk = document.createElement("a");
    historyLnk.textContent = "History";

    let historyBtn = document.createElement("div");
    historyBtn.appendChild(historyLnk);
    historyBtn.setAttribute("class", "btn btn-bling btn-gold history-btn");

    let historyDiv = document.createElement("div");
    historyDiv.setAttribute("class", "wrap cf");
    historyDiv.setAttribute("style", "justify-content: center;display: flex");
    historyDiv.appendChild(historyBtn);

    let historySec = document.createElement("section");
    historySec.setAttribute("class", "blk bg-black food-nav cf");
    historySec.setAttribute("style", "padding: 0");
    historySec.appendChild(historyDiv);

    return historySec;
}

function CreateHistoryDialog() {
    let closeLnk = document.createElement("a")
    closeLnk.setAttribute("inert", "true")
    closeLnk.innerText = "close"

    let closeBtn = document.createElement("div")
    closeBtn.setAttribute("class", "btn btn-bling btn-gold")
    closeBtn.appendChild(closeLnk)
    closeBtn.addEventListener("click", () => {
        document.querySelector("dialog").close()
    });

    let closeDiv = document.createElement("div")
    closeDiv.setAttribute("class", "wrap cf")
    closeDiv.setAttribute("style", "justify-content: center;display: flex");
    closeDiv.appendChild(closeBtn)

    let closeSec = document.createElement("section")
    closeSec.setAttribute("class", "blk food-nav cf")
    closeSec.appendChild(closeDiv)

    let dialogClose = document.createElement("form")
    dialogClose.setAttribute("method", "dialog")
    dialogClose.appendChild(closeSec)

    let historyList = document.createElement("ul")
    historyList.setAttribute("class", "history no-style no-top white")

    let historyHeader = document.createElement("h2")
    historyHeader.setAttribute("class", "sec-title white")
    historyHeader.innerText = "Your History"

    let saveHistoryBtn = document.createElement("div")
    saveHistoryBtn.setAttribute("class", "action-btn")
    saveHistoryBtn.innerText = "Save"
    saveHistoryBtn.addEventListener("click", () => {
        downloadHistory()
    });

    let shareHistoryBtn = document.createElement("div")
    shareHistoryBtn.setAttribute("class", "action-btn")
    shareHistoryBtn.innerText = "Share"
    shareHistoryBtn.addEventListener("click", () => {
        shareHistory()
    });

    let historyTitle = document.createElement("header")
    historyTitle.setAttribute("class", "generic-header")
    historyTitle.appendChild(historyHeader)
    historyTitle.appendChild(saveHistoryBtn)
    historyTitle.appendChild(shareHistoryBtn)

    let historySection = document.createElement("section")
    historySection.setAttribute("style", "overflow: auto; min-width: 100%; min-height: 100%; background-color: rgb(30, 29, 28)")
    historySection.appendChild(historyTitle)
    historySection.appendChild(historyList)
    historySection.appendChild(dialogClose)

    let historyDialog = document.createElement("dialog")
    historyDialog.setAttribute("class", "modal")
    historyDialog.setAttribute("style", "max-width: 90vw")
    historyDialog.appendChild(historySection)

    return historyDialog
}

function addRating(item) {
    let drink = item.firstElementChild;
    if (drink.firstChild) {
        let drinkName = item.querySelector("input").id;
        let drinkKey = drinkName.replace(/\s/g, "").toLowerCase();
        drink.appendChild(createRatingComponent(drinkKey))
    }
}

function setStars(rating, score) {
    let stars = rating.querySelectorAll("span");
    if (score.className.includes("checked")) {
        for (let i = 4; i >= score.id - 1; i--) {
            stars[i].setAttribute("class", "fa fa-star")
        }
    } else {
        for (let i = 0; i < score.id; i++) {
            stars[i].setAttribute("class", "fa fa-star checked")
        }
    }
}

function saveRating(key, score) {
    if (localStorage.log) {
        let drinkLog = new Map(JSON.parse(localStorage.log));
        drinkLog.get(key).score = score
        localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
    }
}

function createRatingComponent(name) {
    let rating = document.createElement("span")
    rating.setAttribute("id", name+".rating")
    rating.setAttribute("class", "rating")
    rating.style.visibility = "hidden"

    for (let i = 0; i < 5; i++) {
        let star = document.createElement("span")
        star.setAttribute("id", i + 1 + "")
        star.setAttribute("class", "fa fa-star")
        star.addEventListener("click", event => {
            let key = event.target.parentElement.id.split(".")[0]
            let score = event.target.id
            console.log("rating " + key + ":" + score)
            setStars(event.target.parentElement, event.target)
            saveRating(key, score)
        })
        rating.appendChild(star)
    }

    return rating
}


function addItemToHistory(id, drink_name, drink_description, finish_date) {
    let historyList = document.querySelector(".history")

    let item = document.createElement("li");
    item.setAttribute("id", id + ".history")

    let drink = document.createElement("h4");
    drink.setAttribute("class", "white");
    drink.innerText = drink_name;

    let date = document.createElement("span");
    date.setAttribute("class", "finished-date gold");
    date.innerText = new Date(finish_date).toLocaleDateString("en-US", dataOptions)

    drink.appendChild(date)

    let description = document.createElement("div");
    description.setAttribute("class", "item-details top white");
    description.innerText = drink_description;

    item.appendChild(drink)
    item.appendChild(description)
    historyList.appendChild(item)
}

function removeItemFromHistory(key) {
    document.getElementById(key + ".history").remove();
}

function shareHistory() {
    let drinkLog = new Map(JSON.parse(localStorage.log));
    let data = JSON.stringify(Object.fromEntries(drinkLog))
    fetch("/wg-drinks/history", {
        credentials: "same-origin",
        mode: "same-origin",
        method: "post",
        headers: { "Content-Type": "application/json" },
        body: data
    })
        .then(resp => {
            if (resp.status === 200) {
                return resp.json()
            } else {
                console.log("Status: " + resp.status)
                return Promise.reject("server")
            }
        })
        .then(respJson => {
            navigator.clipboard.writeText(location.protocol + '//' + location.host + "/wg-drinks/history/"+respJson.id);
            alert("Copied the text: " + location.protocol + '//' + location.host + "/wg-drinks/history/"+respJson.id);
        })
        .catch(err => {
            if (err === "server") return
            console.log(err)
        })
}

window.onload = function () {
    //Add history button to page
    document.querySelector(".generic-header").appendChild(CreateHistoryBtn())

    // Add display action
    document.querySelector(".history-btn").addEventListener("click", () => {
        document.querySelector("dialog").show()
    })

    // Create User History
    let body = document.querySelector("body");
    body.insertBefore(CreateHistoryDialog(), body.firstElementChild)

    // Add checkboxes
    document.querySelectorAll(".menu-content").forEach((menu) => {
        menu.querySelectorAll("li")
            .forEach(item => {
                addCheckbox(item)
            });
    })

    // Add star rating
    document.querySelectorAll("#house-cocktails").forEach((menu) => {
        menu.querySelectorAll("li")
            .forEach(item => {
                addRating(item)
            });
    })

    if (localStorage.log) {
        let drinkLog = new Map(JSON.parse(localStorage.log));
        drinkLog.forEach(function (value, key) {
            if (value.date) {
                let checkbox = document.getElementById(key);
                let rating = document.getElementById(key+".rating");

                // Populate checkboxes
                if (checkbox) {
                    checkbox.checked = true;
                }

                // Show star rating
                if (rating) {
                    rating.style.visibility = "visible";
                }

                // Populate star rating
                if (rating && value.score) {
                    let stars = rating.querySelectorAll("span");
                    for (let i = 0; i < value.score; i++) {
                        stars[i].setAttribute("class", "fa fa-star checked")
                    }
                }

                // Populate history Dialog
                addItemToHistory(key, value.name, value.description, value.date);
            }
        })
    }
}