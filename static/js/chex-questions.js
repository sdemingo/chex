

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
	    success: function(response){
		cb(response)
	    },
	    error: error
	});
    }

    var editQuest = function(q,cb){

    }

    var listTags = function(cb){
	$.ajax({
	    url:DOMAIN+'/questions/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: function(response){
		cb(response)
	    },
	    error: error
	})
    }

    var listQuests = function(tags,cb){
	$.ajax({
	    url:DOMAIN+'/questions/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success:  function(response){
		cb(response)
	    },
	    error: error
	})
    }

    var deleteQuest = function(q,cb){

    }


    /*

      Dom functions 

    */


    var addQuestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear pregunta")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Pregunta creada con éxito")
	    resetForm(settings.form)
	}
    }

    var listQuestResponse = function(response){
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

    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		$(settings.panel+" .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
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
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(settings.panel+" .tags").on("click",function(e){
	    e.preventDefault()
	    tags=[]
	    $(settings.panel+" .results").empty()
	    $(settings.panel+" .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    listQuests(tags,listQuestResponse)
		}
	})
    }





    /*
      
      Public interface of the module

    */


    var init = function() {
	listTags(listTagsResponse)
	bindFunctions()
    }


    return{
	init: init,
	list: listQuests,
	tags: listTags,
	add: addQuest,
	del: deleteQuest
    }

})()


