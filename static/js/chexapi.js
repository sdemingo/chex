

var DOMAIN = ""


function responseHandler(data,success,error){
    if (data.Error){
	console.log("Error response :"+data.Error)
	if (error){
	    error(data.Error)
	}
    }else{
	success(data)
    }
}


function addUser(u,success,error){
    $.ajax({
	url:DOMAIN+'/users/new',
	type: 'post',
	dataType: 'json',
	data: JSON.stringify(u),
	success: function(data){
	    responseHandler(data,success,error)
	},
	error: function(data){
            console.log("Server Internal Error:"+data);
        }
    });
}

