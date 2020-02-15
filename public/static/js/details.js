var rowIds = []
function getDetails(selectRenderCallback) {
    var promise = new Promise(function(resolve, reject) {
        var r = new XMLHttpRequest();
        r.onreadystatechange = function() {
            if (this.readyState == 4) {
                if (this.status == 200) {
                    data = JSON.parse(this.responseText);
                    selectRenderCallback(data)
                    resolve("Username data retrieved");
                } else {
                    reject("Unsuccessful request");
                }
            } 
        }
        r.open("GET", "/details-data", true);
        r.send();
    }) 

    return promise;
}

window.onload = function () {
    getDetails(rowRenderFunc(usernameRowString));

    document.getElementById("table-body").addEventListener("click", function(event){
        event.preventDefault();
        if (event.target.dataset.type === "submit" ){
            sendUpdateDetailsRequest(event.target.dataset.rowid);
        } 
    });

    document.getElementById("omnisave-button").addEventListener("click", function(event){
        event.preventDefault();
        var details = []
        rowIds.forEach(rowID => {
            var bookie = document.getElementsByClassName(`bookie-${rowID}`)[0].innerText;
            var username = document.getElementsByClassName(`username-${rowID}`)[0].value;
            var password = document.getElementsByClassName(`password-${rowID}`)[0].value;
            details.push({
                Bookie: bookie,
                Username: username,
                Password: password,
            })
        }); 
        updateRequest(JSON.stringify(details));
    });
}


function sendUpdateDetailsRequest(rowID) {
    var bookie = document.getElementsByClassName(`bookie-${rowID}`)[0].innerText;
    var username = document.getElementsByClassName(`username-${rowID}`)[0].value;
    var password = document.getElementsByClassName(`password-${rowID}`)[0].value;

    if (password === "") {
        flashStatus(false, "Cannot update details without a password. Please supply a password and try again.");
    } if (username === "") {
        flashStatus(false, "Cannot update details without a username. Please supply a username and try again.");
    } else {
        updateRequest(JSON.stringify([{
            Bookie: bookie,
            Username: username,
            Password: password,
        }]));
    }
}

function updateRequest(package) {
    var r = new XMLHttpRequest();
    r.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (this.status == 200) {
                flashStatus(true, "Successfully updated details");
                clearTable()
                getDetails(rowRenderFunc(usernameRowString));
            } else {
                flashStatus(false, "OOPS, Something went wrong when trying to update details :'(");
            }
        }
    };

    r.open("POST", "/update-details", true);
    r.setRequestHeader('Content-Type', 'application/json');
    r.send(package);
}

function usernameRowString(username) {
    content = "";
    content += `<td class="bookie-${username["ID"]}">${username["Bookie"]}</td>`;
    content += `<td><input type="text" class="username-${username["ID"]}" value="${username["Username"]}"></td>`;
    content += `<td><input type="text" class="password-${username["ID"]}" value=${username["Password"]}></td>`;
    return content + `<td><input class="confirm-balance" type="submit" value="update details" data-type="submit" data-rowid="${username["ID"]}"></td>`;
}