

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



    /*

      Ajax Api

    */


    var addUser =  function(u,cb){
	$.ajax({
	    url:DOMAIN+'/users/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: cb,
	    error: error
	});
    }

    var getUser =  function(id,cb){
	$.ajax({
	    url:DOMAIN+'/users/get?id='+id,
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }

    var editUser = function(u,cb){
	$.ajax({
	    url:DOMAIN+'/users/update',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: cb,
	    error: error
	})
    }

    var listTags = function(cb){
	$.ajax({
	    url:DOMAIN+'/users/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	})
    }

    var listUsers = function(tags,cb){
	$.ajax({
	    url:DOMAIN+'/users/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: cb,
	    error: error
	})
    }

    var deleteUser = function(u,cb){

    }




    /*

      Dom functions 

    */

    // Callback after the add user request
    var addUserResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear usuario")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Usuario creada con éxito")
	    resetForm(settings.form)
	}
    }

    // Callback after the edit user request
    var editUserResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al editar usuario")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Usuario editado con éxito")
	    resetForm(settings.form)
	}
    }

    // Callback after the list user tags request
    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		if (e.trim().length > 0){
		    $(settings.panel+" .tags")
			.append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
			.on("click",selectTag)
		    CHEX.userTags[e]=1
		}
	    })
		}
    }

    // Callback after the list user request 
    var listUsersResponse = function(response){
	if ((!response) || (response.length==0)){
	    $(settings.panel+" .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(e){
		$(settings.panel+" .results")
		    .append("<a href=\"/users/get?id="+e.Id+"\" class=\"list-group-item\">"+e.Name+"</a>")
	    })
	}
    }
    

    // Mark tag as selected 
    var selectTag = function(event){
	event.preventDefault()

	var element = $(this)
	if (element.hasClass("label-primary")) {
            element.removeClass("label-primary");
        }else{
	    element.addClass("label-primary");
	}
    }

    // Recover clicked tags and launch a search by these tags
    var launchSearchByTag = function(){
	tags=[]
	$(settings.panel+" .results").empty()
	$(settings.panel+" .tags").find(".label-primary").each(function(){
	    tags.push($(this).html())
	})

	    if (tags.length>0){
		listUsers(tags,listUsersResponse)
	    }
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
	    addUser(u,addUserResponse)
	})

	// Edit Users
	$(settings.form+" #userUpdateSubmit").click(function(){
	    var u = readForm()
	    if (!u){
		return
	    }
	    editUser(u,editUserResponse)
	})

	// List Users
	$(settings.panel+" .tags").on("click",launchSearchByTag)
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})
    }


    var init = function() {
	listTags(listTagsResponse)
	bindFunctions()
	$(".alert").css("visibility", "hidden")
    }

    return{
	init: init,
	list: listUsers,
	tags: listTags,
	add: addUser,
	get: getUser,
	del: deleteUser
    }
})()



