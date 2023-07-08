var progress1 = document.getElementById("progress1");
var progress1Text = document.getElementById("progressName")

async function getProgressData(url = "/api/progress") {
    const response = await fetch(url)
        //, {
   //  method: "GET",
   //  mode: "cors",
   // // cache: "no",
   //  credentials: "same-origin",
  //  });
    
    return response.json();
}

(function(){
    getProgressData().then((data) => {
    progress1Text.innerHTML = data.filename;
    let completed = (data.done / data.filesize) * 100;
    progress1.innerHTML = completed;
    progress1.style = "width: "+ completed + "%";
    setTimeout(arguments.callee, 500);

    })
})();

