

$(document).ready(function () {
    $(".dropdown-toggle").dropdown();
    users.init()
    questions.init()
})




function showInfoMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-success").text(text)
    window.setTimeout(function() { $("#infoPanel").removeClass("alert-success").css("visibility", "hidden") }, 1500)
}

function showErrorMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-danger").text(text)
    window.setTimeout(function() { $("#infoPanel").removeClass("alert-danger").css("visibility", "hidden") }, 1500)
}



$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
	    if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
	    }
	    o[this.name].push(this.value || '');
        } else {
	    o[this.name] = this.value || '';
        }
    });
    return o;
};

Array.prototype.clean = function(deleteValue) {
    for (var i = 0; i < this.length; i++) {
	if (this[i] == deleteValue) {         
	    this.splice(i, 1);
	    i--;
	}
    }
    return this;
};

var DOMAIN=""

var error = function(data){
    console.log("Internal server error: "+data)
}

/*
  
  Modulo de usuarios

*/

var users = (function(){

    var settings={
	form:"#userEditForm",
	panel:"#usersList"
    }


    var addUser =  function(u){
	$.ajax({
	    url:DOMAIN+'/users/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al crear usuario")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Usuario creado con éxito")
		    resetForm()
		}
	    },
	    error: error
	});
    }

    var editUser = function(u){
	$.ajax({
	    url:DOMAIN+'/users/update',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al editar usuario")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Usuario editado con éxito")
		}
	    },
	    error: error
	});
    }

    var listTags = function(){
	$.ajax({
	    url:DOMAIN+'/users/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: function(data){
		if (data){
		    $.each(data,function(i,e){
			$(settings.panel+" .tags")
			    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
		    })
			}
	    },
	    error: function(data){
		console.log("Server Internal Error:"+data);
            }
	})
    }

    var listUsers = function(tags){
	$.ajax({
	    url:DOMAIN+'/users/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: function(data){
		if ((!data) || (data.length==0)){
		    $(settings.panel+" .results")
			.append("<span class=\"list-group-item\">No hubo resultados</span>")
		}else{
		    data.forEach(function(e){
			$(settings.panel+" .results")
			    .append("<a href=\"/users/get?id="+e.Id+"\" class=\"list-group-item\">"+e.Name+"</a>")
		    })
		}
	    },
	    error: function(data){
		console.log("Server Internal Error:"+data);
            }
	});
    }

    var deleteUser = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
	    }
    
    var readForm = function(){
	var u = $(settings.form).serializeObject()
	u.Tags = u.Tags.split(",").map(function(e){
	    return e.trim()
	})
	u.Tags.clean("")

	if (/\s+/.test(u.Tags.join())){
	    showErrorMessage("Las etiquetas no pueden contener espacios")
	    return
	}
	
	if ((u.Name=="") || (u.Maill=="")){
	    showErrorMessage("Existen campos sin información")
	    return
	}
	return u
    }
    
    var bindFunctions = function(){
	// Add User
	$(settings.form+" #userNewSubmit").click(function(){
	    var u = readForm()
	    if (!u) {
		return
	    }
	    addUser(u)
	})

	// Edit Users
	$(settings.form+" #userUpdateSubmit").click(function(){
	    var u = readForm()
	    if (!u){
		return
	    }
	    editUser(u)
	})

	// List Users
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(settings.panel+" .tags").on("click",function(e){
	    tags=[]
	    $(settings.panel+" .results").empty()
	    $(settings.panel+" .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    listUsers(tags)
		}
	})
    }


    var init = function() {
	listTags()
	bindFunctions()
	$(".alert").css("visibility", "hidden");
    }

    return{
	init: init,
	
    }
})()



/*
  
  Modulo de preguntas

*/

var questions = (function(){
    var settings={
	form:"#questEditForm",
	/*panel:"#usersList"*/
    }


    var addQuest =  function(u){
	
    }

    var editQuest = function(u){

    }

    var listTags = function(){

    }

    var listQuests = function(tags){

    }

    var deleteQuest = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
	    }
    
    var readForm = function(){

    }
    
    var bindFunctions = function(){

	// Edit Quest form
	$(settings.form+" .btn-add").on("click",function(){
	    $(settings.form+" .question-options")
		.append("<div class=\"input-group\">\
<input type=\"text\" class=\"form-control\" name=\"Options\">\
<span class=\"input-group-btn\">\
<button type=\"button\" class=\"btn btn-default btn-del\">-</button></div>\
</span>")
	})

	$(settings.form).on("click",".btn-del",function(){
	    $(this).closest("div.input-group").remove()
	})
    }


    var init = function() {
	listTags()
	bindFunctions()
    }

    return{
	init: init,
	
    }

})()
