let username = document.getElementById("username");
let useremail = document.getElementById("useremail");
let userphoto = document.getElementById("userphoto");
let fieldTotBand = document.getElementById("total_band")
let filedRemBand= document.getElementById("con_band")
//let usedtranfertext = document.getElementById("usedtranfertext")


let usedtransfer = document.getElementById("usedtransfer")

function handleSignoutReq() {
    window.location.href = "https://downdrive.onrender.com/api/account/signout"
}


function handleDeleteReq(){
  let isOK =  window.confirm("Sure you want to  delete this account")
    if(isOK ){
        window.location.href = "https://downdrive.onrender.com/api/account/delete"
    }
}

async function getUserInfo(url = "/api/account/getuser") {
    const response = await fetch(url,
    {
         method: "GET",

    });

    return response.json();
}


(function(){
    getUserInfo().then((data) => {
        console.log(data)
        username.innerHTML = data.username;
        useremail.innerHTML = data.email;
        userphoto.src = data.user_image;

        fieldTotBand.innerText=data.data_allotted + ' MB';
        filedRemBand.innerText=data.data_allotted-data.data_remains + ' MB';
        //usedtransfertext.innerText = data.transferused / (1024 * 1024);
        usedtransfer.innerText = data.data_remains ;
        usedtransfer.style = "width: " + (data.data_remains / data.data_allotted) * 100  + "%";
        

    })
        .catch((err) => {
            console.log(err)
        });
})();