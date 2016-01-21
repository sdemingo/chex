


var testmaker = (function(){

    var readForm = function(form){
	var a = form.serializeObject()
	return a
    }


    var bindFunctions = function(){

	// preload answered in the database of earlier sessions
	$(".answer-panel").each(function(){
	    var answerForm=$(this)
	    var a = readForm(answerForm)
	    if (a) {
		// pedimos resupestas para este ejercicio
		// a /answers/list?exercise=a.ExerciseId&quest=a.QuestId

		answers.list(a.ExerciseId,a.QuestId,function(response){
		    var answer 
		    if (Array.isArray(response) && response.length>0){
			answer=response[0]
			answers.render(answerForm,answer)
		    }
		})
	    }
	})

	// add an answer for an exercise of a test in the server
	$(".submit-answer").click(function(e){
	    e.preventDefault()
	    var form=$(this).parent("form")	    
	    var a = readForm(form)
	    if (!a) {
		return
	    }

	    answers.add(a,function(response){
		if (response.Error){
		    showErrorMessage("Error al enviar la respuesta")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Pregunta contestada!")
		}
	    })
	})
    }
    
    var init = function(){
	bindFunctions()
    }

    return{
	init: init
    }
})()
