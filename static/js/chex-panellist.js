


var panelList = (function(){

  var data={
    selectedQuestions:{},
    selectedUsers:{},
    testsQuestions:{},
    testsUsers:{},
    questionsCache:{},
    usersCache:{}
  }


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
	$(".panel-select-questions .results")
	      .append(
	  $('<li id='+q.Id+' class="list-group-item col-md-12">\
	    <div class="icons row col-md-2">\
	    <a href="#" class="item-select glyphicon glyphicon-ok"></a>\
	    <a href="#" class="item-add glyphicon glyphicon-plus"></a>\
	    </div>\
	    <div class="text col-md-10">\
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
    var q = data.questionsCache[that.attr("id")]
    if (q.SolutionId<0){
      showErrorMessage("Esta pregunta está sin solucionar por parte del autor y no puede ser seleccionada")
      return false
    }
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
    for (var id in data.selectedQuestions) {
      data.testsQuestions[id]=1
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
	      .append(
	$('<li id='+q.Id+' class="list-group-item col-md-12">')
		     .append('<div class="icons col-md-2 text-center">\
			     <input type="text" class="form-control item-input-value good-points"/>\
			     <input type="text" class="form-control item-input-value bad-points"/>\
			     <div class="icons row">\
			     <a href="#" class="item-select glyphicon glyphicon-ok"></a>\
			     <a href="#" class="item-remove glyphicon glyphicon-remove"></a>\
			     </div>\
			     </div>')
		     .on("click",".item-select",selectQuestionHandler)
		     .on("click",".item-remove",removeQuestionsHandler)

		     .append(
	  $('<div class="col-md-10">')
		     .append('<a href="/questions/get?id='+q.Id+'" class="item-link">'+resume(q.Text)+'</a>')
	)
      )
    }
  } 


  // Event handler to select all questions listed
  var selectAllQuestionsHandler = function(event){
    //TODO
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
	$(".panel-select-users .results")
	      .append(
	  $('<li id='+u.Id+' class="list-group-item col-md-12">\
	    <div class="icons row col-md-2">\
	    <a href="#" class="item-select glyphicon glyphicon-ok"></a>\
	    <a href="#" class="item-add glyphicon glyphicon-plus"></a>\
	    </div>\
	    <div class="text col-md-10">\
	    <a class="item-text" href="/users/get?id='+u.Id+'" >'+u.Name+'</a>\
	    </div>\
	    </li>')
								   .on("dblclick",".item-select",selectAllUsersHandler)
								   .on("click",".item-select",selectUserHandler)
								   .on("click",".item-add",addUserHandler)
	)
	
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
      data.testsUsers[id]=1
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
      
      $("#testAddedUserPanel .results")
	      .append(
	$('<li id='+u.Id+' class="list-group-item col-md-12">')
		     .append('<div class="row icons col-md-2">\
			     <a href="#" class="item-select glyphicon glyphicon-ok"></a>\
			     <a href="#" class="item-remove glyphicon glyphicon-remove"></a>\
			     </div>')
		     .on("dblclick",".item-select",selectAllUsersHandler)
		     .on("click",".item-select",selectUserHandler)
		     .on("click",".item-remove",removeUserHandler)

		     .append(
	  $('<div class="col-md-10">')
		     .append('<a href="/users/get?id='+u.Id+'" class="item-link">'+u.Name+'</a>')
	)
      )
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

    // Show questions for select them
    $("#addMoreQuests").click(function(){
      $("#testAddedQuestionPanel").hide()
      $("#testSelectQuestionPanel").show()
      listQuestionsTags(listQuestionsTagsResponse)
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


  var init = function() {
    $("#testSelectQuestionPanel").hide()
    $("#testAddedQuestionPanel ul").empty()

    $("#testSelectUserPanel").hide()
    $("#testAddedUserPanel ul").empty()
    bindFunctions()
  }

  return{
    init: init,
    getSelectedUsers:function(){return data.selectedUsers}
  }

})()

