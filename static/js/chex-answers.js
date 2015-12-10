


/*
  
  Modulo de respuestas

*/

var answers = (function(){
    var settings={
	form:"#answerEditForm",
	panel:"#answerPanel"
    }


    /*

      Ajax Api

    */


    var addAnswer =  function(a,cb){
	$.ajax({
	    url:DOMAIN+'/answers/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(a),
	    success: cb,
	    error: error
	});
    }

    var editAnswer = function(u){

    }


    var listAnswers = function(){

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


    
    var readForm = function(){
	var a = $(settings.form).serializeObject()

	// RawSolution must be marshall into a simple string always
	if (Array.isArray(a.RawBody)){
	    a.RawSolution = a.RawBody.toString()
	}
	return a
    }
    
    var bindFunctions = function(){
	// si existe la solución  la mostramos primeramente
	if ($(settings.panel+" #solvedPanel").length){
	    $(settings.panel+" #unSolvedPanel").hide()
	}

	// ocultar la solución y mostrar el formulario para editar
	$(settings.panel+" #answerUpdateButton").on("click",function(){
	    $(settings.panel+" #solvedPanel").hide()
	    $(settings.panel+" #unSolvedPanel").show()
	})

	// ocultar el formulario de respuesta y mostrar la solución
	$(settings.panel+" #answerNewCancel").on("click",function(){
	    $(settings.panel+" #unSolvedPanel").hide()
	    $(settings.panel+" #solvedPanel").show()
	})

	// crea una nueva respuesta o actualiza la existente
	$(settings.panel+" #answerNewSubmit").on("click",function(){
	    var a = readForm()
	    if (!a) {
		return
	    }
	    addAnswer(a,addAnswerResponse)
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init,	
    }

})()


