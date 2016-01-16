


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
	$.ajax({
	    url:DOMAIN+'/tests/update',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(test),
	    success: cb,
	    error: error
	});
    }

    var listTests = function(tags,cb){
	$.ajax({
	    url:DOMAIN+'/tests/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: cb,
	    error: error
	});
    }

    var listTestsForUser = function(id,cb){
	$.ajax({
	    url:DOMAIN+'/tests/list',
	    type: 'get',
	    dataType: 'json',
	    data: {foruser:id},
	    success: cb,
	    error: error
	});
    }


    var listTags = function(cb){
	if (!cb){
	    cb=listTagsResponse
	}
	$.ajax({
	    url:DOMAIN+'/tests/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	})
    }

    var listUsersAllowed = function(test,cb){
	$.ajax({
	    url:DOMAIN+'/tests/users/list?id='+test.Id,
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }

    var listExercises = function(test,cb){
	$.ajax({
	    url:DOMAIN+'/tests/exercises/list?id='+test.Id,
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	});
    }

    var deleteTest = function(test,cb){

    }
    

    /*

      Private and Dom functions 

    */


    // Callback after the add test request
    var addTestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear test")
	    console.log(response.Error)
	}else{
	    showInfoMessage("Test creado con éxito")
	    resetForm(settings.form)
	}
    }

    // Callback after the edit test request
    var editTestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al editar test")
	    console.log(response.Error)
	}else{
	    showInfoMessage("Test editado con éxito")
	}
    }

    // Callback after the list user tags request
    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		if (e.trim().length > 0){
		    // $(settings.panel+" .tags")
		    // 	.append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
		    // 	.on("click",selectTag)
		    CHEX.testTags[e]=1
		}
	    })
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
	    e.Id=0
	    e.QuestId=questIds[i]
	    e.BadPoint=parseFloat($("#testAddedQuestionPanel #"+questIds[i]+" .bad-points").first().val())
	    if (!e.BadPoint){
		e.BadPoint=0
	    }
	    e.GoodPoint=parseFloat($("#testAddedQuestionPanel #"+questIds[i]+" .good-points").first().val())
	    if (!e.GoodPoint){
		e.GoodPoint=0
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
	    addTest(tst,addTestResponse)
	})

	// Update test button
	$("#testUpdateSubmit").click(function(){
	    var tst = readForm()
	    if (!tst) {
		return
	    }
	    editTest(tst,editTestResponse)
	})
    }


    var init = function() {
	//listTags(listTagsResponse)
	bindFunctions()
    }

    return{
	init: init,
	tags: listTags,
	list: listTests,
	listFor: listTestsForUser,
	listUsers: listUsersAllowed,
	listExercises: listExercises
    }

})()




