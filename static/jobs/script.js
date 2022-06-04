function start() {
    var statusSpan = document.getElementById("status")
    var container = document.getElementById("container")
    var jobId = new URLSearchParams(window.location.search).get("id")
    if(jobId == null || jobId == "") {
        window.location.replace('/')
    }
    statusSpan.innerHTML = "Checking status"
    var looper = setInterval(() => {
        fetch("/status/" + jobId)
        .then(response => response.json())
        .then(data => {
            if(data.status == 1) {
                if(data.result) {
                    container.innerHTML = ""
                    statusSpan.innerHTML = "Finished"
                    for(let i in data.result) {
                        let ref = data.result[i]
                        container.innerHTML = container.innerHTML + 
                        `
                            <div>
                                <h3>${ref.title}</h3>
                                <h4>Similarity: ${ref.similarity}%</h4>
                                <p>${ref.description}</p>
                                <a href="${ref.link}">${ref.link}</a>
                            </div><br>
                        `
                    }
                }
                else {
                    statusSpan.innerHTML = "Error"
                }
                clearInterval(looper)
            } else if(data.status == -1) {
                statusSpan.innerHTML = "Error"
                clearInterval(looper)
            } else {
                statusSpan.innerHTML = "Running"
            }
        })
        .catch(err => console.log(err))
    },5000)
}

window.addEventListener('load', start)