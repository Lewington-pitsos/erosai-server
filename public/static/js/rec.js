window.onload = function() {
    token = new URLSearchParams(window.location.search).get("token")

    var r = new XMLHttpRequest();
    r.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (this.status == 200) {
                var resp = JSON.parse(this.response);

                var allElements = ""
                resp.forEach(function(element) {
                    allElements += `<li class="recommendation"><p>Porn Confidence: ${element.Porn}</p><a href="${element.URL}">${element.URL}</a></li>`
                    console.log(element);
                });

                document.getElementById("rec-holder").innerHTML = allElements;


            } else if (this.status == 400) {
                alert("Oops... something went mildly wrong, maybe refresh and try again.")
            } else {
                alert("OOPS! Something went hideously wrong, contact an admin.")
            }
        }
    };
    
    r.open("GET", "/get-recommendations?token=" + token, true);
    r.send();
}