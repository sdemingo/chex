

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

    var listTags = function(panel,results){
	$.ajax({
	    url:DOMAIN+'/questions/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: function(data){
		if (data){
		    $.each(data,function(i,e){
			$(panel)
			    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
		    })
			}
		if (results){
		    results.allTags=data
		}
	    },
	    error: error
	})
    }

    var listQuests = function(panel,results){
	$.ajax({
	    url:DOMAIN+'/questions/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:results.seletedTags.join(",")},
	    success: function(data){
		if ((!data) || (data.length==0)){
		    $(panel+" .results")
			.append("<span class=\"list-group-item\">No hubo resultados</span>")
		}else{
		    data.forEach(function(e){
			$(panel+" .results")
			    .append("<li class=\"list-group-item\"><a href=\"/questions/get?id="+e.Id+"\" >"+resume(e.Text)+"</a></li>")
		    })
		}
		if (results){
		    results.quests=data
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
	initTagPanel(settings.panel,{})
    }





    // Public methods of module

    var initTagPanel = function(panel,results){

	listTags(panel+" .tags",results)

	$(panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(panel+" .tags").on("click",function(e){
	    e.preventDefault()
	    results.seletedTags=[]
	    $(panel+" .results").empty()
	    $(panel+" .tags").find(".label-primary").each(function(){
		results.seletedTags.push($(this).html())
	    })
		if (results.seletedTags.length>0){
		    listQuests(panel,results)
		}
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init,
	tags: initTagPanel
    }

})()


