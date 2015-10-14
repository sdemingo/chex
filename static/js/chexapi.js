

var DOMAIN = ""


function addUser(u,success,error){
    $.ajax({
	url:DOMAIN+'/users/edit',
	type: 'post',
	dataType: 'json',
	data: JSON.stringify(u),
	success: success,
	error: function(data){
	    if (error) { error(data) }
            console.log("add user error:"+data);
        }
    });
}

