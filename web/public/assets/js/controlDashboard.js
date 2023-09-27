const container = document.getElementById("containerClass")
const firstProgressClass = document.getElementById("progressClass0")

let globalCount = 0

async function getProgressData(url = "/api/progress") {
    const response = await fetch(url)
    return response.json();
}


(function(){


    getProgressData().then((data) => {
        console.log(globalCount)
        if (globalCount < data.length) {
            // let rem = data.length - globalCount
            for (let i = globalCount; i < data.length ; i++){
                if (i === 0) {
                    let firstProgressName = document.getElementById("progressName0")
                    firstProgressName.innerText = data[i].filename
                    console.log("name of first file is " + firstProgressName.innerHtml)
                }else{
                    let temProgressClass = firstProgressClass.cloneNode(true)
                    temProgressClass.setAttribute("id","progressClass" + i)

                    let temProgressName = temProgressClass.querySelector(".progressName")
                    temProgressName.setAttribute("id","progressName" + i )
                    temProgressName.textContent = data[i].filename


                    let temProgressBar = temProgressClass.querySelector(".progress .progress-bar")
                    temProgressBar.setAttribute("id","progress" + i )
                    container.append(temProgressClass)
                }

            }
            globalCount = data.length
            console.log(globalCount)

        }
        for (let i = 0; i < data.length; i++){
            let temProgressBar = document.getElementById("progress"+i)
            let temCompleted = (data[i].done / data[i].filesize) * 100;
            temProgressBar.innerHTML = math.round(temCompleted);
            temProgressBar.style = "width:" + temCompleted +"%"
        }



    setTimeout(arguments.callee, 1000);

    })
        .catch((err) =>{
            console.log(err)
            console.log("some error happend with data")
        })
})();

