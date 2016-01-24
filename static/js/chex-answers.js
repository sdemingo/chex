


/*
  
  Modulo de respuestas

*/

var answers = (function(){
    var settings={
	form:"#answerEditForm",
	panel:"#answerPanel"
    }

    var TYPE_TESTSINGLE   = 1
    var TYPE_TESTMULTIPLE = 2
    
    var changeSolution=false

    /*

      Ajax Api

    */


    var addAnswer =  function(a,cb){
	$.ajax({
	    url:DOMAIN+'/tests/exercises/do',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(a),
	    success: cb,
	    error: error
	});
    }

    var addSolutionAnswer =  function(a,cb){
	$.ajax({
	    url:DOMAIN+'/questions/solve',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(a),
	    success: cb,
	    error: error
	});
    }

    var listAnswers = function(ex,quest,cb){
	$.ajax({
	    url:DOMAIN+'/answers/list?exercise='+ex+'&quest='+quest,
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }

    var deleteAnswer = function(){

    }



    /*

      Dom functions 

    */

    // Callback after the add quest request
    var addAnswerResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear respuesta")
	    console.log(data.Error)
	}else{
	    window.setTimeout(function(){location.reload()}, 2000)
	}
    }

    var readForm = function(form){
	var a
	if (!form){
	    a = $(settings.form).serializeObject()
	}else{
	    a = $(form).serializeObject()
	}

	// RawSolution must be marshall into a simple string always
	if (Array.isArray(a.RawBody)){
	    a.RawSolution = a.RawBody.toString()
	}
	return a
    }

    var markAnswer = function(form, a, isSolution){

	if (a.BodyType == TYPE_TESTSINGLE){
	    console.log(a.Body.Solution)
	    var c=$(form)
		.find("[name=RawBody][value="+a.Body.Solution+"]")
		.first().prop("checked","true").addClass("marked")

	    if (isSolution){
		c.addClass("correct")
	    }
	}
	
	form.addClass("marked")
    }
  
  
    var bindFunctions = function(){

	$(".answer-panel").not(".marked").each(function(){
	    var answerForm=$(this)
	    var a = readForm(answerForm)
	    if (a) {
		listAnswers(a.ExerciseId,a.QuestId,function(response){
		    // if the exercise.Id in the request is 0 means that
		    // it has been request the teacher solution of the quest
		    var answer 
		    var isSolution = (a.ExerciseId == 0)
		    if (Array.isArray(response) && response.length>0){
			answer=response[0]   
			markAnswer(answerForm,answer,isSolution)
		    }
		})
	    }
	})

	// TODO:
	// request all teachers solution to the answer and remark the exercise
	// to compare with the student answer. Only in check-panels from check view 
	// of tests

	/*
	$(".check-panel.marked").each(function(){

	    var answerForm=$(this)
	    var a = readForm(answerForm)
	    if (a) {
		listAnswers(0,a.QuestId,function(response){
		    var answer 
		    var isSolution = true
		    if (Array.isArray(response) && response.length>0){
			answer=response[0]   
			markAnswer(answerForm,answer,isSolution)
		    }
		})

	})
	*/


	$(".answer-panel").on("change", function() {
	    if ($(this).hasClass("marked")){
		changeSolution=true
	    }
	})

	// crea una nueva respuesta o actualiza la existente
	$(".answer-panel .submit").on("click",function(){
	    var form = $(this).parent("form")
	    var a = readForm(form)
	    if (!a) {
		return
	    }  

	    if (changeSolution){
		var msg="Este ejercicio ya tenía una solución anterior. ¿Desea introducir la nueva repuesta?"
		showConfirmMessage(msg,function(){
		    addSolutionAnswer(a,addAnswerResponse)
		})
	    }else{
		addSolutionAnswer(a,addAnswerResponse)
	    }
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init,
	add: addAnswer,
	list: listAnswers
    }

})()


