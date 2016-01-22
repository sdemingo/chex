


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

    var renderAnswer = function(form, a){

	if (a.BodyType == TYPE_TESTSINGLE){
	    console.log(a.Body.Solution)
	    var c=$(form)
		.find("[name=RawBody][value="+a.Body.Solution+"]")
		.first().prop("checked","true")
	}
	
    }
  
  
    var bindFunctions = function(){

	// preload answered in the database
	$(".answer-panel").each(function(){
	    var answerForm=$(this)
	    var a = readForm(answerForm)
	    if (a) {
		answers.list(a.ExerciseId,a.QuestId,function(response){
		    var answer 
		    if (Array.isArray(response) && response.length>0){
			answer=response[0]
			answers.render(answerForm,answer)
		    }
		})
	    }
	})


	$(".answer-panel").on( "change", function() {
	    changeSolution=true
	})

	// crea una nueva respuesta o actualiza la existente
	$(".answer-panel .submit").on("click",function(){
	    
	    var form = $(this).parent("form")
	    var a = readForm(form)
	    if (!a) {
		return
	    }  

	    if (changeSolution){
		var msg="La solución ha cambiado. ¿Desea introducir la nueva repuesta?"
		showConfirmMessage(msg,function(){
		    addSolutionAnswer(a,addAnswerResponse)
		})
	    }
	    addSolutionAnswer(a,addAnswerResponse)
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init,
	add: addAnswer,
	list: listAnswers,
	render: renderAnswer
    }

})()


