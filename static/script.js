window.onload = function() {
    if (localStorage.log) {
        drinkLog = new Map(JSON.parse(localStorage.log));
        drinkLog.forEach(function (value, key) {
            console.log(key);
            var checkbox = document.getElementById(key);
            if (checkbox) {
                checkbox.checked = true;
            }
        })
    }
}

document.addEventListener("click", (event) => {
    let drinkLog;
    if (event.target.checked) {
        if (localStorage.log) {
            drinkLog = new Map(JSON.parse(localStorage.log));
            drinkLog.set(event.target.id, new Date());
            localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
        } else {
            drinkLog = new Map();
            drinkLog.set(event.target.id, new Date());
            localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
        }
        console.log("finished " + event.target.id);
    } else {
        if (localStorage.log) {
            drinkLog = new Map(JSON.parse(localStorage.log));
            drinkLog.delete(event.target.id, new Date())
            localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
        } else {
            drinkLog = new Map()
            localStorage.log = JSON.stringify(Array.from(drinkLog.entries()));
        }
        console.log("removed" + event.target.id);
    }
});