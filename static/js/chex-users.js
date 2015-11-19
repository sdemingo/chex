

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
	//panel:"#usersList"
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
    }


    var init = function() {
	bindFunctions()
    }


    return{
	init: init,
	list: listUsers,
	tags: listTags,
	add: addUser,
	del: deleteUser
    }
})()






var usersFinder = (function(){

    var settings={}

    // Callback after the list user tags request
    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		$(settings.panel+" .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
		    .on("click",selectTag)
	    })
		}
    }

    // Callback after the list user request 
    var listUsersResponse = function(response){
	if ((!response) || (response.length==0)){
	    $(settings.panel+" .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(u){
		var item = $('<li id='+u.Id+' class="list-group-item col-md-12">')
		    //.append('<div class="icons col-md-2">')
		    .append('<div class="text col-md-10">\
<a class="item-text" href="/users/get?id='+u.Id+'" >'+u.Name+'</a>\
</div>')
		
		if (settings.itemSelectHandler 
		    || settings.itemAddHandler 
		    || settings.itemRemoveHandler 
		    || settings.itemEditHandler){

		    item.prepend('<div class="icons col-md-2">')
		}

		if (settings.itemSelectHandler){
		    item.find(".icons").prepend('<a href="#" class="item-select glyphicon glyphicon-ok"></a>')
		    item.on("dblclick",".item-select",function(){alert("Seleccionas todos")})
		    item.on("click",".item-select",settings.itemSelectHandler)
		}

		if (settings.itemAddHandler){
		    item.find(".icons").prepend('<a href="#" class="item-add glyphicon glyphicon-plus"></a>')
		    item.on("click",".item-add",settings.itemAddHandler)
		}

		if (settings.itemRemoveHandler){
		    item.find(".icons").prepend('<a href="#" class="item-remove glyphicon glyphicon-remove"></a>')
		    item.on("click",".item-remove",settings.itemRemoveHandler)
		}

		if (settings.itemEditHandler){
		    item.find(".icons").prepend('<a href="#" class="item-edit glyphicon glyphicon-edit"></a>')
		    item.on("click",".item-edit",settings.itemEditHandler)
		}

		$(settings.panel+" .results").append(item)
	    })
	}
    }
    

    // Mark tag as selected 
    var selectTag = function(event){
	event.preventDefault()

	var element = $(this)
	if (!element.is("li")){
	    return
	}
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
		users.list(tags,listUsersResponse)
	    }
    }


    var bindFunctions = function(){
	// List Users
	$(settings.panel+" .tags").on("click",launchSearchByTag)
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})
    }


    var buildComponent = function(){
	$(settings.panel)
	    .append(
		$('<ul class="col-md-12 tags">')
	    )
	    .append(
		$('<ul class="col-md-12 list-group results">')
	    )

	if (settings.closeButton){
	    $(settings.panel).prepend(
		$('<button type="button" class="btn btn-default pull-right">Cerrar</button>')
		.on("click",settings.closeButton)
	    )
	}
    }


    var init = function(options){
	settings=options
	buildComponent()

	$(settings.panel+" .tags").empty()
	users.tags(listTagsResponse)
	bindFunctions()
    }


    return{
	init: init
    }
})()


