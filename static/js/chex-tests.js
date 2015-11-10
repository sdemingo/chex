

/*
  
  Modulo de tests

*/

var questions = (function(){
    var settings={
	form:"#testsEditForm",
	panel:""
    }


    var addTest =  function(q){
	$.ajax({
	    url:DOMAIN+'/tests/add',
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

    var editTest = function(u){

    }

    var listTags = function(){

    }

    var listTests = function(tags){

    }

    var deleteTest = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
	    }
    
    var readForm = function(){
	/*var q = $(settings.form).serializeObject()
	q.Tags = q.Tags.split(",").map(function(e){
	    return e.trim()
	})
	q.Tags.clean("")
	
	return q*/
    }
    
    var bindFunctions = function(){
	//$("#testQuestionPanel")

	// Add questions button
	$(settings.form+" #testNewSubmit").click(function(){
	    
	})

	// Add users button
    }


    var init = function() {
	$("#testSelectQuestionPanel").hide()
	$("#testSelectedQuestionPanel ul").empty()
	listTags()
	bindFunctions()
    }

    return{
	init: init,
    }

})()


