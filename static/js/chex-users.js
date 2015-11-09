

/*
  
  Modulo de usuarios

*/

var users = (function(){

    var types={
	Name:"isNonEmpty",
	Mail:"isEmail",
	Tags:"isWordEnumeration"
    }

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
	})
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
	    error: error
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
	    error: error
	})
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

	validator.validate(u,types)
	if (validator.hasErrors()){
	    showErrorMessage("Existen campos mal formados o sin información")
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
	$(".alert").css("visibility", "hidden")
    }

    return{
	init: init,
    }
})()



