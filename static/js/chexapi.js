

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
	url:DOMAIN+'/users/add',
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


function editUser(u,success,error){
    $.ajax({
	url:DOMAIN+'/users/update',
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


function getTags(success,error){
    $.ajax({
	url:DOMAIN+'/users/tags',
	type: 'get',
	dataType: 'json',
	success: function(data){
	    responseHandler(data,success,error)
	},
	error: function(data){
            console.log("Server Internal Error:"+data);
        }
    });
}

function getUsers(filt,success,error){
    $.ajax({
	url:DOMAIN+'/users/list',
	type: 'get',
	dataType: 'json',
	data: {tags:filt.join(",")},
	success: function(data){
	    responseHandler(data,success,error)
	},
	error: function(data){
            console.log("Server Internal Error:"+data);
        }
    });
}
