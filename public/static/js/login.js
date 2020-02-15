window.onload = function () {
    var password = prompt("Please Enter your password")

    var r = new XMLHttpRequest();
    r.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (this.status == 200) {
                console.log(this.response)
                document.cookie = this.response
                alert("authentication successful, you can now access Dank Wagering Systems")
                document.location.pathname = "/details"
            } else if (this.status == 401) {
                alert("authentication failed, that password was incorrect")
                document.location.reload(true)
            } else if (this.status == 400) {
                document.location.reload(true)
                alert("Oops... something went mildly wrong, maybe refresh and try again.")
            } else {
                alert("OOPS! Something went hideously wrong, contact an admin.")
            }
        }
    };
    
    r.open("POST", "/login-attempt", true);
    r.setRequestHeader('Content-Type', 'application/json');
    r.send(JSON.stringify({
        Password: password,
    }));
}