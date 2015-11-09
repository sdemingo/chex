

/*
  
  Modulo de preguntas

*/

var questions = (function(){
    var settings={
	form:"#questEditForm",
	panel:"#questList"
    }


    var addQuest =  function(q){
	$.ajax({
	    url:DOMAIN+'/questions/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(q),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al crear pregunta")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Pregunta creada con éxito")
		    resetForm()
		}
	    },
	    error: error
	});
    }

    var editQuest = function(u){

    }

    var listTags = function(){
	$.ajax({
	    url:DOMAIN+'/questions/tags/list',
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

    var listQuests = function(tags){
	$.ajax({
	    url:DOMAIN+'/questions/list',
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
			    .append("<a href=\"/questions/get?id="+e.Id+"\" class=\"list-group-item\">"+resume(e.Text)+"</a>")
		    })
		}
	    },
	    error: error
	})
    }

    var deleteQuest = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
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
	    addQuest(q)
	})

	// List Quests
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
		    listQuests(tags)
		}
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


