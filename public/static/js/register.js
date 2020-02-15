window.onload = function() {
    document.getElementById("body").addEventListener('click', function(event) {
        event.preventDefault();
        if (event.target.id == "submit") {

            username = document.getElementById("username").value;
            password = document.getElementById("password").value;
            console.log(username);
            var r = new XMLHttpRequest();
            r.onreadystatechange = function() {
                if (this.readyState == 4) {
                    if (this.status == 200) {
                        var resp = JSON.parse(this.response);

                        console.log(resp)

                        if (resp.Outcome != "Success") {
                          flashStatus(false, resp.Outcome + " (Registration failed)");
                        } else {
                            flashStatus(true, "Account registered successfully");
                        }
                    } else if (this.status == 400) {
                        document.location.reload(true)
                        alert("Oops... something went mildly wrong, maybe refresh and try again.")
                    } else {
                        alert("OOPS! Something went hideously wrong, contact an admin.")
                    }
                }
            };
            
            r.open("POST", "/register-attempt", true);
            r.setRequestHeader('Content-Type', 'application/json');
            r.send(JSON.stringify({
                Username: username,
                Password: password,
            }));
        } 
        // alert("hello");
        //validation code to see State field is mandatory.  
    });
}