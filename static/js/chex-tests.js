


/*
  
  Modulo de tests

*/

var tests = (function(){
    var settings={
	form:"",
	panel:""
    }

    var data={
	selectedQuestions:{},
	questionsCache:{}
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

    // Mark question as selected 
    var selectQuestion = function(element){
	if (element.hasClass("list-group-item-info")) {
            element.removeClass("list-group-item-info");
	    delete data.selectedQuestions[$(this).attr("id")]
        }else{
	    element.addClass("list-group-item-info");
	    data.selectedQuestions[element.attr("id")]=1
	}
    }

    // Listed the questions tags for search questions
    var listQuestionsTags = function(cb){
	$("#testSelectQuestionPanel .results").empty()
	questions.tags(cb)
    }

    // Callback after list questions by tag request
    var listQuestionsResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $("#testSelectQuestionPanel .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(q){
		$("#testSelectQuestionPanel .results")
		    .append(
			$('<li id='+q.Id+' class="list-group-item col-md-12">\
<div class="icons col-md-1">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
<a href="#" class="item-add glyphicon glyphicon-plus"></a>\
</div>\
<div class="text col-md-11">\
<a class="item-text" href="/questions/get?id='+q.Id+'" >'+resume(q.Text)+'</a>\
</div>\
</li>')
			    .on("click",".item-select",selectQuestionHandler)
			    .on("click",".item-add",addQuestionsHandler)
		    )
		
		data.questionsCache[q.Id]=q
		1    })
	}
    }

    // Callback after lists tags request
    var listQuestionsTagsResponse = function(response){
	if (response){
	    $("#testSelectQuestionPanel .tags").empty()
	    $.each(response,function(i,e){
		$("#testSelectQuestionPanel .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    // Event handler to select a question
    var selectQuestionHandler = function(event){
	event.preventDefault()
	var that = $(this).parent().parent()
	selectQuestion(that)
    }


    // Event handler to add a question or a selected questions to a
    // test
    var addQuestionsHandler = function(event){
	event.preventDefault()
	
	// mark this question as selected and add all of them
	var that = $(this).parent().parent()
	selectQuestion(that)

	listQuestionsSelected()
	$("#testSelectedQuestionPanel").show()
	$("#testSelectQuestionPanel").hide()
    }

/*
	      <li class="list-group-item col-md-12">
		<div class="icons col-md-2 text-center">
		  <input type="text" class="form-control item-input-value good-points"/>
		  <input type="text" class="form-control item-input-value bad-points"/>
		  <a href="" class="item-remove glyphicon glyphicon-remove"></a>
		</div>
		<div class="text col-md-10">
		  <a href="" class="item-link">Cras justo odioLorem ipsum dolor sit amet, consectetuer adipiscing elit. Donec hendrerit tempor tellus. Donec pretium posuere tellus. Proin quam nisl, tincidunt et, mattis eget, convallis nec, purus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nulla posuere. Donec vitae dolor.
		  </a>
		</div>
	      </li>
*/


    // List every questions selected
    var listQuestionsSelected = function(){
	$("#testSelectedQuestionPanel ul").empty()

	for (var id in data.selectedQuestions) {
	    q = data.questionsCache[id]
	    if (!q){
		return
	    }
	    
	    
	    $("#testSelectedQuestionPanel .results")
		.append(
		    $('<li id='+q.Id+' class="list-group-item col-md-12">')
			.append('<div class="icons col-md-2 text-center">\
		                   <input type="text" class="form-control item-input-value good-points"/>\
                                   <input type="text" class="form-control item-input-value bad-points"/>\
                       		  <a href="" class="item-remove glyphicon glyphicon-remove"></a>\
                      		</div>')
			.append(
			    $('<div class="col-md-10">')
				.append('<a href="/questions/get?id='+q.Id+'" class="item-link">'+resume(q.Text)+'</a>')
			)
		)
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
	
	// Cancel the selection question form and return to selected
	// question panel
	$("#cancelSelectedQuests").click(function(){
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


