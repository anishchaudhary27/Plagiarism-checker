function start() {
    var statusP = document.getElementById("status")
    var resultP = document.getElementById("result")
    var jobID = new URLSearchParams(window.location.search).get("jobId")
    if(jobID == null || jobID == "") {
        statusP.innerHTML = "Upload file"
        return
    }
    statusP.innerHTML = "Checking status"
    var looper = setInterval(() => {
        fetch("/status/" + jobID)
        .then(response => response.json())
        .then(data => {
            if(data.status == 1) {
                statusP.innerHTML = "Finished"
                resultP.innerHTML = data.result
                clearInterval(looper)
            } else if(data.status == -1) {
                statusP.innerHTML = "Error"
                resultP.innerHTML = data.result
                clearInterval(looper)
            } else {
                statusP.innerHTML = "Running"
                resultP.innerHTML = "..."
            }
        })
        .catch(err => console.log(err))
    },5000)
}

window.addEventListener('load', start)