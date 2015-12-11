


/*

  Modulo de tests

*/

var tests = (function(){
    var settings={
	form:"",
	panel:""
    }

    /*

      Ajax Api

    */

    var addTest =  function(test,cb){
	$.ajax({
	    url:DOMAIN+'/tests/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(test),
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


    // Callback after the add user request
    var addTestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear test")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Test creado con éxito")
	    resetForm(settings.form)
	}
    }

    var buildExerciseList = function(){
	var exercises=[]
	var added=questionsList.getAdded()	
	var questIds=Object.keys(added).map(function(x){
	    return parseInt(x,10)
	})

	for (var i=0;i<questIds.length;i++){
	    var e={}
	    e.QuestId=questIds[i]
	    e.BadPoint=parseInt($("#testAddedQuestionPanel #"+questIds[i]+" .bad-points").first().val(),10)
	    if (!e.BadPoint){
		e.BadPoint="0"
	    }
	    e.GoodPoint=parseInt($("#testAddedQuestionPanel #"+questIds[i]+" .good-points").first().val(),10)
	    if (!e.GoodPoint){
		e.GoodPoint="0"
	    }
	    exercises.push(e)
	}

	return exercises
    }


    var readForm = function(){
	var tst = $("#testEditForm").serializeObject()
	tst.Tags = tst.Tags.split(",").map(function(e){
	    return e.trim()
	})
	tst.Tags.clean("")
	tst.State = 1
	tst.Exercises = buildExerciseList()

	tst.Ulist = Object.keys(usersList.getAdded()).map(function(x){
	    return parseInt(x,10)
	})
	
	return tst
    }
    
    var bindFunctions = function(){

	// Add test button
	$("#testNewSubmit").click(function(){
	    var tst = readForm()
	    if (!tst) {
		return
	    }
	    alert(JSON.stringify(tst))
	    //addTest(tst,addTestResponse)
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init
    }

})()




