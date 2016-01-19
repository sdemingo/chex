


var testmaker = (function(){

    var readForm = function(form){
	var a = form.serializeObject()
	return a
    }


    var bindFunctions = function(){
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
