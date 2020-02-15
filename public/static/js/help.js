var audioCtx 
const dingURL = "/assets/ding2.wav"
const highDingURL = "/assets/high-ding.wav"
var dingBuffer
var highDingBuffer

function playBuffer(audioBuffer) {
    const source = audioCtx.createBufferSource();
    source.buffer = audioBuffer;
    source.connect(audioCtx.destination);
    source.start();
  }

function contains(array, obj) {
    var i = array.length;
    while (i--) {
        if (array[i] == obj) {
            return true;
        }
    }
    return false;
}

function currentDatePlus(days) {
    today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth()+1; //As January is 0.
    var yyyy = today.getFullYear();
    return pad(mm, 2)+"/"+pad(dd + days, 2)+"/"+yyyy;
}

function getNotificationsFunc(path, rowRenderCallback, key) {
    return function() {
        var r = new XMLHttpRequest();
        r.onreadystatechange = function() {
            if (this.readyState == 4) {
                if (this.status == 200) {
                    data = JSON.parse(this.responseText);
                    if (isValid(data)) {
                        changeStatus(true)
                        rowRenderCallback(data[key])
                    } else {
                        changeStatus(false)
                    }
                } else {
                    changeStatus(false)
                }
            }
        };
        
        r.open("GET", path, true)
        r.send()
    }
}

function rowRenderFunc(rowStringCallback) {
    return function renderRows(data) {
        var table = document.getElementById("table-body");
        for (var index in data) {
            var notification = data[index]
            if (rowIds.indexOf(notification["ID"]) == -1) {
                addNewRow(notification, table, rowStringCallback);
                rowIds.push(notification["ID"])
            }
        }
    }
}

function addNewRow(data, table, rowStringCallback) {
    var newRow = document.createElement("tr");
    newRow.setAttribute("id", `row-${data["ID"]}`);
    newRow.innerHTML = rowStringCallback(data, newRow);
    table.appendChild(newRow);
}

function playSound(buffer) {
    playBuffer(buffer)
}

function tableCell(value) {
    return `<td>${value}</td>`
}

function getRow(rowID) {
    return document.getElementById(`row-${rowID}`);
}

function flashStatus(success, text) {
    var alertBox = document.getElementById("alert-box");
    if (success) {
        makeSuccessful(alertBox)
    } else {
        makeError(alertBox)
    }
    alertBox.firstElementChild.innerHTML = text
    alertBox.classList.remove("hidden")
    setTimeout(function() {
        alertBox.classList.add("hidden")
    }, 3000)
}

function makeError(alertBox) {
    alertBox.classList.remove("success")
    alertBox.classList.add("error")
}

function makeSuccessful(alertBox) {
    alertBox.classList.remove("error")
    alertBox.classList.add("success")
}

function removeRow(rowID) {
    var row = getRow(rowID)
    row.parentNode.removeChild(row);
}

function defaultValue(value) {
    if (value === "") {
        return "0"
    }

    return value
}

function clearElement(element) {
    element.innerHTML = ""
}

function clearTable(element) {
    clearElement(document.getElementById("table-body"))
    rowIds = []
}

function pad(num, size) {
    var s = num+"";
    while (s.length < size) s = "0" + s;
    return s;
}