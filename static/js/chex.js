

var DOMAIN=""


$(document).ready(function () {
    $(".dropdown-toggle").dropdown()
    users.init()
    questions.init()
    answers.init()
    tests.init()
})


/*
  
  Modulo de usuarios

*/

var users = (function(){

    var types={
	Name:"isNonEmpty",
	Mail:"isEmail",
	Tags:"isWordEnumeration"
    }

    var settings={
	form:"#userEditForm",
	panel:"#usersList"
    }

    var addUser =  function(u){
	$.ajax({
	    url:DOMAIN+'/users/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al crear usuario")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Usuario creado con éxito")
		    resetForm()
		}
	    },
	    error: error
	});
    }

    var editUser = function(u){
	$.ajax({
	    url:DOMAIN+'/users/update',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(u),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al editar usuario")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Usuario editado con éxito")
		}
	    },
	    error: error
	})
    }

    var listTags = function(){
	$.ajax({
	    url:DOMAIN+'/users/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: function(data){
		if (data){
		    $.each(data,function(i,e){
			$(settings.panel+" .tags")
			    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
		    })
			}
	    },
	    error: error
	})
    }

    var listUsers = function(tags){
	$.ajax({
	    url:DOMAIN+'/users/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: function(data){
		if ((!data) || (data.length==0)){
		    $(settings.panel+" .results")
			.append("<span class=\"list-group-item\">No hubo resultados</span>")
		}else{
		    data.forEach(function(e){
			$(settings.panel+" .results")
			    .append("<a href=\"/users/get?id="+e.Id+"\" class=\"list-group-item\">"+e.Name+"</a>")
		    })
		}
	    },
	    error: error
	})
    }

    var deleteUser = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
	    }
    
    var readForm = function(){
	var u = $(settings.form).serializeObject()
	u.Tags = u.Tags.split(",").map(function(e){
	    return e.trim()
	})
	u.Tags.clean("")

	validator.validate(u,types)
	if (validator.hasErrors()){
	    showErrorMessage("Existen campos mal formados o sin información")
	    return 
	}

	return u
    }
    
    var bindFunctions = function(){
	// Add User
	$(settings.form+" #userNewSubmit").click(function(){
	    var u = readForm()
	    if (!u) {
		return
	    }
	    addUser(u)
	})

	// Edit Users
	$(settings.form+" #userUpdateSubmit").click(function(){
	    var u = readForm()
	    if (!u){
		return
	    }
	    editUser(u)
	})

	// List Users
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(settings.panel+" .tags").on("click",function(e){
	    tags=[]
	    $(settings.panel+" .results").empty()
	    $(settings.panel+" .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    listUsers(tags)
		}
	})
    }


    var init = function() {
	listTags()
	bindFunctions()
	$(".alert").css("visibility", "hidden")
    }

    return{
	init: init,
    }
})()





/*
  
  Modulo de preguntas

*/

var questions = (function(){
    var settings={
	form:"#questEditForm",
	panel:"#questList"
    }


    /*

      Ajax Api

    */

    var addQuest =  function(q,cb){
	$.ajax({
	    url:DOMAIN+'/questions/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(q),
	    success: cb,
	    error: error
	});
    }

    var editQuest = function(q,cb){

    }

    var listTags = function(cb){
	$.ajax({
	    url:DOMAIN+'/questions/tags/list',
	    type: 'get',
	    dataType: 'json',
	    success: cb,
	    error: error
	})
    }

    var listQuests = function(tags,cb){
	$.ajax({
	    url:DOMAIN+'/questions/list',
	    type: 'get',
	    dataType: 'json',
	    data: {tags:tags.join(",")},
	    success: cb,
	    error: error
	})
    }

    var deleteQuest = function(q,cb){

    }


    /*

      Dom functions 

    */


    var addQuestResponse = function(response){
	if (response.Error){
	    showErrorMessage("Error al crear pregunta")
	    console.log(data.Error)
	}else{
	    showInfoMessage("Pregunta creada con éxito")
	    resetForm(settings.form)
	}
    }

    var listQuestResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(settings.panel+" .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(e){
		$(settings.panel+" .results")
		    .append("<li class=\"list-group-item\"><a href=\"/questions/get?id="+e.Id+"\" >"+resume(e.Text)+"</a></li>")
	    })
	}
    }

    var listTagsResponse = function(response){
	if (response){
	    $.each(response,function(i,e){
		$(settings.panel+" .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }

    
    var readForm = function(){
	var q = $(settings.form).serializeObject()
	q.Tags = q.Tags.split(",").map(function(e){
	    return e.trim()
	})
	q.Tags.clean("")
	
	return q
    }
    

    var bindFunctions = function(){
	// Edit Quest form
	$(settings.form+" .btn-add").on("click",function(){
	    $(settings.form+" .question-options")
		.append("<div class=\"input-group\">\
<input type=\"text\" class=\"form-control\" name=\"Options\">\
<span class=\"input-group-btn\">\
<button type=\"button\" class=\"btn btn-default btn-del\">-</button></div>\
</span>")
	})

	$(settings.form).on("click",".btn-del",function(){
	    $(this).closest("div.input-group").remove()
	})


	// Add Quest
	$(settings.form+" #questNewSubmit").click(function(){
	    var q = readForm()
	    if (!q) {
		return
	    }
	    addQuest(q,addQuestResponse)
	})

	// List Quests
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(settings.panel+" .tags").on("click",function(e){
	    e.preventDefault()
	    tags=[]
	    $(settings.panel+" .results").empty()
	    $(settings.panel+" .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    listQuests(tags,listQuestResponse)
		}
	})
    }





    /*
      
      Public interface of the module

    */


    var init = function() {
	listTags(listTagsResponse)
	bindFunctions()
    }


    return{
	init: init,
	list: listQuests,
	tags: listTags,
	add: addQuest,
	del: deleteQuest
    }

})()





/*
  
  Modulo de respuestas

*/

var answers = (function(){
    var settings={
	form:"#answerEditForm",
	panel:"#answerPanel"
    }


    var addAnswer =  function(a){
	$.ajax({
	    url:DOMAIN+'/answers/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(a),
	    success: function(data){
		if (data.Error){
		    showErrorMessage("Error al crear respuesta")
		    console.log(data.Error)
		}else{
		    showInfoMessage("Respuesta creada con éxito")
		    //resetForm()
		}
	    },
	    error: error
	});
    }

    var editAnswer = function(u){

    }


    var listAnswers = function(){

    }

    var deleteAnswer = function(){

    }

    var resetForm = function(){
	$(settings.form).each(function(){this.reset()})
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
	    addAnswer(a)
	    location.reload()
	})
    }


    var init = function() {
	bindFunctions()
    }

    return{
	init: init,	
    }

})()




/*
  
  Modulo de tests

*/

var tests = (function(){
    var settings={
	form:"",
	panel:"#testSelectQuestionPanel"
    }


    /*

      Ajax Api

    */

    var addTest =  function(test,cb){
	$.ajax({
	    url:DOMAIN+'/tests/add',
	    type: 'post',
	    dataType: 'json',
	    data: JSON.stringify(q),
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

    // Listed the questions tags for search questions
    var listQuestionsTags = function(cb){
	questions.tags(cb)
    }

    // Callback after list questions request
    var listQuestionsResponse = function(response){
	if ((!response) || (response.length==0) || !Array.isArray(response)){
	    $(settings.panel+" .results")
		.append("<span class=\"list-group-item\">No hubo resultados</span>")
	}else{
	    response.forEach(function(e){
		$(settings.panel+" .results")
		    .append("<li class=\"list-group-item\"><a href=\"/questions/get?id="+e.Id+"\" >"+resume(e.Text)+"</a></li>")
	    })
	}
    }

    // Callback after lists tags request
    var listQuestionsTagsResponse = function(response){
	if (response){
	    $(settings.panel+" .tags").empty()
	    $.each(response,function(i,e){
		$(settings.panel+" .tags")
		    .append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		}
    }


    var readForm = function(){

    }
    
    var bindFunctions = function(){

	// Add questions button
	$(settings.form+" #testNewSubmit").click(function(){
	    
	})

	// Show questions for select them
	$("#addMoreQuests").click(function(){
	    $("#testSelectedQuestionPanel").hide()
	    $("#testSelectQuestionPanel").show()
	    listQuestionsTags(listQuestionsTagsResponse)
	})

	// Add selected quests and show all
	$("#addSelectedQuests").click(function(){
	    $("#testSelectedQuestionPanel").show()
	    $("#testSelectQuestionPanel").hide()
	})


	// List Tests
	$(settings.panel+" .tags").on("click","*",function(e){
	    $(this).toggleClass("label-primary")
	})

	$(settings.panel+" .tags").on("click",function(e){
	    e.preventDefault()
	    tags=[]
	    $(settings.panel+" .results").empty()
	    $(settings.panel+" .tags").find(".label-primary").each(function(){
		tags.push($(this).html())
	    })
		if (tags.length>0){
		    questions.list(tags,listQuestionsResponse)
		}
	})

    }


    var init = function() {
	$("#testSelectQuestionPanel").hide()
	$("#testSelectedQuestionPanel ul").empty()
	bindFunctions()
    }

    return{
	init: init,
    }

})()



var validator = {
    
    types:{
	isNonEmpty :{
	    validate:function(value){
		return value!= ""
	    },
	    instructions: "value cannot be empty"
	},
	isNumber : {
	    validate:function(value){
		return !isNaN(value)
	    },
	    instructions: "value must be a number"
	},
	isWordEnumeration : {
	    validate:function(value){
		return (/^(\s*\w+\s*,)*\s*\w+\s*$/m.test(value))
	    },
	    instructions: "value must a word without spaces sequence"
	},
	isEmail : {
	    validate:function(value){
		var re = /^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$/i;
		return (re.test(value))
	    },
	    instructions: "value must a valid email"
	}		
    },
    
    config:{},
    
    messages:[],
    
    validate:function(data,types){
	var i,msg,type,checker,result
	
	this.messages=[]
	for (i in data){
	    if (data.hasOwnProperty(i)){
		type = types[i]
		checker = this.types[type]
		if (!type){
		    continue
		}
		if (!checker){
		    console.log("Error: no checker for this type")
		}
		result = checker.validate(data[i])
		if (!result){
		    msg = "Invalid value for "+i+":, "+checker.instructions
		    this.messages.push(msg)
		}
	    }
	}
	return this.hasErrors()
    },

    hasErrors: function(){
	return this.messages.length !=0
    },

    getErrors: function(){
	m=this.messages.join("\n")
	this.messages=[]
	return m
    }
}




function resetForm(form){
    $(form).each(function(){
	this.reset()
    })
}


function error (data){
    console.log("Internal server error: "+data)
}


function resume(text,max){
    if (!max){
	max=150
    }
    if (text.length > max){
	return text.substring(0,max)+" ..."
    }else{
	return text
    }
}


function showInfoMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-success").text(text)
    window.scrollTo(0,0);
    window.setTimeout(function() { $("#infoPanel").removeClass("alert-success").css("visibility", "hidden") }, 1500)
}

function showErrorMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-danger").text(text)
    window.scrollTo(0,0);
    window.setTimeout(function() { $("#infoPanel").removeClass("alert-danger").css("visibility", "hidden") }, 1500)
}


$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
	    if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
	    }
	    o[this.name].push(this.value || '');
        } else {
	    o[this.name] = this.value || '';
        }
    });
    return o;
};

Array.prototype.clean = function(deleteValue) {
    for (var i = 0; i < this.length; i++) {
	if (this[i] == deleteValue) {         
	    this.splice(i, 1);
	    i--;
	}
    }
    return this;
};
