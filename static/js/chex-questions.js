

/*
  
  Modulo de preguntas

*/

var questions = (function(){
    var settings={
	form:"#questEditForm",
	panel:"#questList"
    }


    /*

      Ajax Api

    */

    var addQuest =  function(q,cb){
	$.ajax({
	    url:DOMAIN+'/questions/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(q),
	    success: cb,
	    error: error
	});
    }

    var getQuest =  function(id,cb){
	$.ajax({
	    url:DOMAIN+'/questions/get?id='+id,
	    type: 'post',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }

    var doQuest =  function(id,cb){
	$.ajax({
	    url:DOMAIN+'/questions/do?id='+id,
	    type: 'post',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }


    var editQuest = function(q,cb){

    }

    var listTags = function(cb){
	if (!cb){
	    cb=listTagsResponse
	}
	$.ajax({
	    url:DOMAIN+'/questions/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	})
    }

    var listQuests = function(tags,cb){
	$.ajax({
	    url:DOMAIN+'/questions/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: cb,
	    error: error
	})
    }

    var deleteQuest = function(q,cb){

    }


    /*

      Dom functions 

    */


    // Callback after the add quest request
    var addQuestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear pregunta")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Pregunta creada con éxito")
	    resetForm(settings.form)
	}
    }

    // Callback after the list quest request
    var listQuestsResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(settings.panel+" .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(e){
		$(settings.panel+" .results")
		    .append("<li class=\"list-group-item\"><a href=\"/questions/get?id="+e.Id+"\" >"+resume(e.Text)+"</a></li>")
	    })
	}
    }


    // Callback after the list quests tags request
    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		if (e.trim().length > 0){
		    $(settings.panel+" .tags")
			.append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
			.on("click",selectTag)
		    CHEX.questionTags[e]=1
		}
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
    var launchSearchByTag = function(event){
	var tags=[]
	$(settings.panel+" .results").empty()
	$(settings.panel+" .tags").find(".label-primary").each(function(){
	    tags.push($(this).html())
	})

	    if (tags.length>0){
		listQuests(tags,listQuestsResponse)
	    }
    }

    
    var readForm = function(){
	var q = $(settings.form).serializeObject()
	q.Tags = q.Tags.split(",").map(function(e){
	    return e.trim()
	})
	q.Tags.clean("")
	
	return q
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


	// Add Quest
	$(settings.form+" #questNewSubmit").click(function(){
	    var q = readForm()
	    if (!q) {
		return
	    }
	    addQuest(q,addQuestResponse)
	})

	// List Quests
	$(settings.panel+" .tags").on("click",launchSearchByTag)
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})
    }





    /*
      
      Public interface of the module

    */


    var init = function() {
	//listTags(listTagsResponse)
	bindFunctions()
    }


    return{
	init: init,
	list: listQuests,
	tags: listTags,
	add: addQuest,
	get: getQuest,
	del: deleteQuest
    }

})()


