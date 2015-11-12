

/*
  
  Modulo de tests

*/

var tests = (function(){
    var settings={
	form:"",
	panel:"#testSelectQuestionPanel"
    }


    /*

      Ajax Api

    */

    var addTest =  function(test,cb){
	$.ajax({
	    url:DOMAIN+'/tests/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(q),
	    success: cb,
	    error: error
	});
    }

    var editTest = function(test,cb){

    }

    var listTests = function(tags,cb){
	
    }

    var deleteTest = function(test,cb){

    }
    

    /*

      Private and Dom functions 

    */

    // Listed the questions tags for search questions
    var listQuestionsTags = function(cb){
	questions.tags(cb)
    }

    // Callback after list questions request
    var listQuestionsResponse = function(response){
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

    // Callback after lists tags request
    var listQuestionsTagsResponse = function(response){
	if (response){
	    $(settings.panel+" .tags").empty()
	    $.each(response,function(i,e){
		$(settings.panel+" .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    var readForm = function(){

    }
    
    var bindFunctions = function(){

	// Add questions button
	$(settings.form+" #testNewSubmit").click(function(){
	    
	})

	// Show questions for select them
	$("#addMoreQuests").click(function(){
	    $("#testSelectedQuestionPanel").hide()
	    $("#testSelectQuestionPanel").show()
	    listQuestionsTags(listQuestionsTagsResponse)
	})

	// Add selected quests and show all
	$("#addSelectedQuests").click(function(){
	    $("#testSelectedQuestionPanel").show()
	    $("#testSelectQuestionPanel").hide()
	})


	// List Tests
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
		    questions.list(tags,listQuestionsResponse)
		}
	})

    }


    var init = function() {
	$("#testSelectQuestionPanel").hide()
	$("#testSelectedQuestionPanel ul").empty()
	bindFunctions()
    }

    return{
	init: init,
    }

})()


