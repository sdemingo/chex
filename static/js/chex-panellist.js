





var questionsList = (function(){

    var data={
	selectedQuestions:{},
	testsQuestions:{},
	questionsCache:{}
    }

    var settings={}

    /*

      Private and Dom functions 

    */


    // Mark question as selected 
    var selectItem = function(element,list){
	if (element.hasClass("list-group-item-info")) {
	    element.removeClass("list-group-item-info");
	    delete list[element.attr("id")]
	}else{
	    element.addClass("list-group-item-info");
	    list[element.attr("id")]=1
	}
    }


    // Listed the questions tags for search questions
    var listQuestionsTags = function(cb){
	$(".panel-select-questions .results").empty()
	questions.tags(cb)
    }


    // Callback after list questions by tag request
    var listQuestionsResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(".panel-select-questions .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(q){
		var li=$('<li id='+q.Id+' class="list-group-item col-md-12">\
<div class="icons row col-md-2">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
</div>\
<div class="text col-md-10">\
<a class="item-text" href="/questions/get?id='+q.Id+'" >'+resume(q.Text)+'</a>\
</div>\
</li>').on("click",".item-select",selectQuestionHandler)
		
		if (settings && settings.addItemIcon){
		    li.find(".icons").append('<a href="#" class="item-add glyphicon glyphicon-plus"></a>')
		    li.on("click",".item-add",addQuestionsHandler)
		}

		$(".panel-select-questions .results").append(li)
		
		data.questionsCache[q.Id]=q
		1    })
	}
    }

    // Callback after lists tags request
    var listQuestionsTagsResponse = function(response){
	if (response){
	    $(".panel-select-questions .tags").empty()
	    $.each(response,function(i,e){
		$(".panel-select-questions .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    // Event handler to select a question
    var selectQuestionHandler = function(event){
	event.preventDefault()
	var that = $(this).parents("li.list-group-item")

	selectItem(that,data.selectedQuestions)
    }


    // Event handler to add questions to the tests collection
    var addQuestionsHandler = function(event){
	event.preventDefault()
	
	// mark this question as selected and add all of them
	var that = $(this).parents("li.list-group-item")
	if (!data.selectedQuestions[that.attr("id")]){
	    selectItem(that,data.selectedQuestions)
	}

	// check if questions are solutioned
	for (var id in data.selectedQuestions) {
	    var q = data.questionsCache[id]
	    if (q.SolutionId<0){
		showErrorMessage("Han sido seleccionadas preguntas sin solucionar. Estas no pueden ser añadidas a un test")
		return false
	    }
	}

	for (var id in data.selectedQuestions) {
	    data.testsQuestions[id]=data.questionsCache[id]
	}

	// dump questions selected
	data.selectedQuestions={}

	listTestQuestions()
	$("#testAddedQuestionPanel").show()
	$("#testSelectQuestionPanel").hide()
    }



    // Event handler to remove questions selected from the test collection
    var removeQuestionsHandler = function(event){
	event.preventDefault()

	// mark this question as selected and remove all of them
	var that = $(this).parents("li.list-group-item")
	if (!data.selectedQuestions[that.attr("id")]){
	    selectItem(that,data.selectedQuestions)
	}
	for (var id in data.selectedQuestions) {
	    delete data.testsQuestions[id]
	}

	// dump questions selected
	data.selectedQuestions={}
	
	listTestQuestions()
    }


    // List every questions selected
    var listTestQuestions = function(){
	$("#testAddedQuestionPanel ul").empty()

	for (var id in data.testsQuestions) {
	    q = data.questionsCache[id]
	    if (!q){
		return
	    }
	    
	    $("#testAddedQuestionPanel .results")
	    var li =$('<li id='+q.Id+' class="list-group-item col-md-12">')
		.append('<div class="icons col-md-2 text-center">\
<input type="text" class="form-control item-input-value good-points"/>\
<input type="text" class="form-control item-input-value bad-points"/>\
<div class="icons row">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
</div>\
</div>')
	    
	    li.on("click",".item-select",selectQuestionHandler)

	    if (settings && settings.removeItemIcon){
		li.find(".icons.row").append('<a href="#" class="item-remove glyphicon glyphicon-remove"></a>')
		li.on("click",".item-remove",removeQuestionsHandler)
	    }

	    li.append(
		$('<div class="col-md-10">')
		    .append('<a href="/questions/get?id='+q.Id+'" class="item-link">'+resume(q.Text)+'</a>')
	    )

	    
	    $("#testAddedQuestionPanel .results").append(li)

	}
    } 
    


    // Event handler to select all questions listed
    var selectAllQuestionsHandler = function(event){
	//TODO
    }



    var bindFunctions = function(){

	// Show questions for select them
	$("#addMoreQuests").click(function(){
	    $("#testAddedQuestionPanel").hide()
	    $("#testSelectQuestionPanel").show()
	    
	})
	
	// Cancel the select questions action
	$("#cancelSelectedQuests").click(function(){
	    $("#testAddedQuestionPanel").show()
	    $("#testSelectQuestionPanel").hide()
	})
	

	// List Questions Tags
	$(".panel-select-questions .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(".panel-select-questions .tags").on("click",function(e){
	    e.preventDefault()
	    var tags=[]
	    $(".panel-select-questions .results").empty()
	    $(".panel-select-questions .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    questions.list(tags,listQuestionsResponse)
		}
	})

	listQuestionsTags(listQuestionsTagsResponse)
    }


    var init = function(options) {
	settings=options
	$("#testSelectQuestionPanel").hide()
	$("#testAddedQuestionPanel ul").empty()

	bindFunctions()
    }

    return{
	init: init,
	getSelected:function(){return data.selectedQuestions},
	getAdded:function(){return data.testsQuestions}
    }

})()




















var usersList = (function(){

    var data={
	selectedUsers:{},
	testsUsers:{},
	usersCache:{}
    }

    var settings={}

    /*

      Private and Dom functions 

    */


    // Mark question as selected 
    var selectItem = function(element,list){
	if (element.hasClass("list-group-item-info")) {
	    element.removeClass("list-group-item-info");
	    delete list[element.attr("id")]
	}else{
	    element.addClass("list-group-item-info");
	    list[element.attr("id")]=1
	}
    }




    // Listed the users tags for search questions
    var listUserTags = function(cb){
	$(".panel-select-users .results").empty()
	users.tags(cb)
    }


    // Callback after list questions by tag request
    var listUserResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(".panel-select-users .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(u){
		var li=$('<li id='+u.Id+' class="list-group-item col-md-12">\
<div class="icons row col-md-2">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
</div>\
<div class="text col-md-10">\
<a class="item-text" href="/users/get?id='+u.Id+'" >'+u.Name+'</a>\
</div>\
</li>')
		    .on("dblclick",".item-select",selectAllUsersHandler)
		    .on("click",".item-select",selectUserHandler)

		if (settings && settings.addItemIcon){
		    li.find(".icons").append('<a href="#" class="item-add glyphicon glyphicon-plus"></a>')
		    li.on("click",".item-add",addUserHandler)
		}

		$(".panel-select-users .results").append(li)
		
		data.usersCache[u.Id]=u
		1    })
	}
    }

    // Callback after lists user tags request
    var listUserTagsResponse = function(response){
	if (response){
	    $(".panel-select-users .tags").empty()
	    $.each(response,function(i,e){
		$(".panel-select-users .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    // Event handler to select a user
    var selectUserHandler = function(event){
	event.preventDefault()
	var that = $(this).parents("li.list-group-item")
	selectItem(that,data.selectedUsers)
    }


    // Event handler to add users to the tests collection
    var addUserHandler = function(event){
	event.preventDefault()
	
	// mark this user as selected and add all of them
	var that = $(this).parents("li.list-group-item")
	if (!data.selectedUsers[that.attr("id")]){
	    selectItem(that,data.selectedUsers)
	}
	for (var id in data.selectedUsers) {
	    data.testsUsers[id]=data.usersCache[id]
	}

	// dump users selected
	data.selectedUsers={}

	listTestUsers()
	$("#testAddedUserPanel").show()
	$("#testSelectUserPanel").hide()
    }


    // List every users added
    var listTestUsers = function(){
	$("#testAddedUserPanel ul").empty()

	for (var id in data.testsUsers) {
	    u = data.usersCache[id]
	    if (!u){
		return
	    }

	    var li = $('<li id='+u.Id+' class="list-group-item col-md-12">')
		.append('<div class="row icons col-md-2">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
\
</div>')
	    
	    li.on("dblclick",".item-select",selectAllUsersHandler)
	    li.on("click",".item-select",selectUserHandler)
	    if (settings && settings.removeItemIcon){
		li.find(".icons").append('<a href="#" class="item-remove glyphicon glyphicon-remove"></a>')
		li.on("click",".item-remove",removeUserHandler)
	    }
	    li.append(
		$('<div class="col-md-10">')
		    .append('<a href="/users/get?id='+u.Id+'" class="item-link">'+u.Name+'</a>')
	    )

	    $("#testAddedUserPanel .results").append(li)
	    
	}	    
    } 


    // Event handler to remove users selected from the test collection
    var removeUserHandler = function(event){
	event.preventDefault()

	// mark this user as selected and remove all of them
	var that = $(this).parents("li.list-group-item")
	if (!data.selectedUsers[that.attr("id")]){
	    selectItem(that,data.selectedUsers)
	}
	for (var id in data.selectedUsers) {
	    delete data.testsUsers[id]
	}

	// dump users selected
	data.selectedUsers={}
	
	listTestUsers()
    }


    // Event handler to select all users listed
    var selectAllUsersHandler = function(event){
	var panel = $(this).parents(".panel-selection").first()

	$(panel).find(" .results li").each(function(i,item){
	    var that=$(this)
	    if (!data.selectedUsers[that.attr("id")]){
		selectItem(that,data.selectedUsers)
	    }
	})
	    }


    var bindFunctions = function(){

	// Show users for select them
	$("#addMoreUsers").click(function(){
	    $("#testAddedUserPanel").hide()
	    $("#testSelectUserPanel").show()
	    //listUserTags(listUserTagsResponse)
	})

	// Cancel the select users action
	$("#cancelSelectedUser").click(function(){
	    $("#testAddedUserPanel").show()
	    $("#testSelectUserPanel").hide()
	})


	// List Users Tags
	$(".panel-select-users .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(".panel-select-users .tags").on("click",function(e){
	    e.preventDefault()
	    var tags=[]
	    $(".panel-select-users .results").empty()
	    $(".panel-select-users .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    users.list(tags,listUserResponse)
		}
	})

	// Select All Users action
	$("#selectAllUsers").click(selectAllUsersHandler)

	listUserTags(listUserTagsResponse)
    }


    var init = function(options) {
	settings=options
	
	// load users added before
	if (settings.test){
	    tests.listUsers(settings.test, function(users){
		users.forEach(function(u){
		    data.usersCache[u.Id]=u
		    data.testsUsers[u.Id]=u
		})
		listTestUsers()
	    })
	}

	$("#testSelectUserPanel").hide()
	$("#testAddedUserPanel ul").empty()
	bindFunctions()
    }

    return{
	init: init,
	getSelected:function(){return data.selectedUsers},
	getAdded:function(){return data.testsUsers}
    }

})()








var testsList = (function(){

    var data={
	selectedTests:{},
	testsAdded:{},
	testsCache:{}
    }

    var settings={}

    /*

      Private and Dom functions 

    */


    // Mark question as selected 
    var selectItem = function(element,list){
	if (element.hasClass("list-group-item-info")) {
	    element.removeClass("list-group-item-info");
	    delete list[element.attr("id")]
	}else{
	    element.addClass("list-group-item-info");
	    list[element.attr("id")]=1
	}
    }




    // Listed the users tags for search questions
    var listTestsTags = function(cb){
	$(".panel-select-tests .results").empty()
	tests.tags(cb)
    }


    // Callback after list questions by tag request
    var listTestResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(".panel-select-tests .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(t){
		var li=$('<li id='+t.Id+' class="list-group-item col-md-12">\
<div class="icons row col-md-2">\
<a href="#" class="item-select glyphicon glyphicon-ok"></a>\
</div>\
<div class="text col-md-10">\
<a class="item-text" href="/tests/get?id='+t.Id+'" >'+t.Title+'</a>\
</div>\
</li>')
		    .on("dblclick",".item-select",selectAllTestsHandler)
		    .on("click",".item-select",selectTestHandler)

		if (settings && settings.addItemIcon){
		    li.find(".icons").append('<a href="#" class="item-add glyphicon glyphicon-plus"></a>')
		    li.on("click",".item-add",addUserHandler)
		}

		$(".panel-select-tests .results").append(li)
		
		data.testsCache[t.Id]=t
		1    })
	}
    }

    // Callback after lists user tags request
    var listTestsTagsResponse = function(response){
	if (response){
	    $(".panel-select-tests .tags").empty()
	    $.each(response,function(i,e){
		$(".panel-select-tests .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    // Event handler to select a user
    var selectTestHandler = function(event){
	event.preventDefault()
	var that = $(this).parents("li.list-group-item")
	selectItem(that,data.selectedTests)
    }

    /*
    // Event handler to add users to the tests collection
    var addUserHandler = function(event){
    event.preventDefault()
    
    // mark this user as selected and add all of them
    var that = $(this).parents("li.list-group-item")
    if (!data.selectedTests[that.attr("id")]){
    selectItem(that,data.selectedTests)
    }
    for (var id in data.selectedTests) {
    data.testsAdded[id]=data.testsCache[id]
    }

    // dump users selected
    data.selectedTests={}

    listTestUsers()
    $("#testAddedUserPanel").show()
    $("#testSelectUserPanel").hide()
    }
    */
    /*
    // List every users added
    var listTestUsers = function(){
    $("#testAddedUserPanel ul").empty()

    for (var id in data.testsAdded) {
    u = data.testsCache[id]
    if (!u){
    return
    }

    var li = $('<li id='+u.Id+' class="list-group-item col-md-12">')
    .append('<div class="row icons col-md-2">\
    <a href="#" class="item-select glyphicon glyphicon-ok"></a>\
    \
    </div>')
    
    li.on("dblclick",".item-select",selectAllUsersHandler)
    li.on("click",".item-select",selectTestHandler)
    if (settings && settings.removeItemIcon){
    li.find(".icons").append('<a href="#" class="item-remove glyphicon glyphicon-remove"></a>')
    li.on("click",".item-remove",removeUserHandler)
    }
    li.append(
    $('<div class="col-md-10">')
    .append('<a href="/users/get?id='+u.Id+'" class="item-link">'+u.Name+'</a>')
    )

    $("#testAddedUserPanel .results").append(li)
    
    }	    
    } 
    */

    // Event handler to remove users selected from the test collection
    var removeTestHandler = function(event){
	event.preventDefault()

	// mark this user as selected and remove all of them
	var that = $(this).parents("li.list-group-item")
	if (!data.selectedTests[that.attr("id")]){
	    selectItem(that,data.selectedTests)
	}
	for (var id in data.selectedTests) {
	    delete data.testsAdded[id]
	}

	// dump users selected
	data.selectedTests={}
	
	listTestUsers()
    }


    // Event handler to select all users listed
    var selectAllTestsHandler = function(event){
	var panel = $(this).parents(".panel-selection").first()

	$(panel).find(" .results li").each(function(i,item){
	    var that=$(this)
	    if (!data.selectedTests[that.attr("id")]){
		selectItem(that,data.selectedTests)
	    }
	})
	    }


    var bindFunctions = function(){

	/*	// Show users for select them
		$("#addMoreUsers").click(function(){
		$("#testAddedUserPanel").hide()
		$("#testSelectUserPanel").show()
		})

		// Cancel the select users action
		$("#cancelSelectedUser").click(function(){
		$("#testAddedUserPanel").show()
		$("#testSelectUserPanel").hide()
		})
	*/

	// List Users Tags
	$(".panel-select-tests .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(".panel-select-tests .tags").on("click",function(e){
	    e.preventDefault()
	    var tags=[]
	    $(".panel-select-tests .results").empty()
	    $(".panel-select-tests .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    tests.list(tags,listTestResponse)
		}
	})

	// Select All Users action
	$("#selectAllUsers").click(selectAllTestsHandler)

	listTestsTags(listTestsTagsResponse)
    }


    var init = function(options) {
	settings=options

	//$("#testSelectUserPanel").hide()
	//$("#testAddedUserPanel ul").empty()
	bindFunctions()
    }

    return{
	init: init,
	getSelected:function(){return data.selectedTests},
	getAdded:function(){return data.testsAdded}
    }

})()















