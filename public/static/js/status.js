function isValid(data) {
    return "Success" in data && data["Success"] == true
}

function changeStatus(ok) {
    if (ok) {
        updateStatus("Status:OK")
        updateDots()
    } else {
        updateStatus("Status:OOPS! Something's not right. Go yell at Louka.")
    }

}

function updateStatus(status) {
    updateDotCount(0)
    document.getElementById("status").innerHTML = status
}

function updateDots() {
    dotCount = updateDotCount(dotCount)
    updateVisibleDots(dotCount)
}

function updateDotCount(dotCount) {
    if (dotCount > 30) {
        return 0
    }

    return dotCount + 1
}

function updateVisibleDots(dotCount) {
    var dotHolder = document.getElementById("dots")
    dotHolder.innerHTML = ".".repeat(dotCount)
}