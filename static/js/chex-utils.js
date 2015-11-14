
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








/*
  - Run the modal with a button, using the same id for the modal an for data-target attr form button:
  <button id="..." type="button" class="btn btn-primary btn-lg" data-toggle="modal" data-target="#myDialog">

  - Operate the modal manually:
  $('#myDialog').modal('toggle');
  $('#myDialog').modal('show');
  $('#myDialog').modal('hide');

  - Init the modal with an object as:
  var modalData={
     id:"myDialog",
     titleText:"Titulo del dialogo",
     bodyText:"Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
     actionButton:{
     text:"Ok",
     handler: actionButtonAnyHandler
  }
*/


var modal = {
    init:function(data){
	$("body").append(
	    $('<div class="modal fade" id="'+data.id+'" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">')
		.append(
		    $('<div class="modal-dialog" role="document">')
			.append(
			    $('<div class="modal-content">')
				.append(
				    $('<div class="modal-header">')
					.append(
					    $('<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>')
					)
					.append(
					    $('<h4 class="modal-title" id="myModalLabel">'+data.titleText+'</h4>')
					)
				)
			    
				.append(
				    $('<div class="modal-body">')
					.append(data.bodyText)
				)

				.append(
				    $('<div class="modal-footer">')
					.append(
					    $('<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>')
					)
				)
			)
		)
	)

	// insert action button
	if (data.actionButton){
	    $(".modal .modal-footer").append(
		$('<button type="button" class="btn btn-primary">'+data.actionButton.text+'</button>')
	    )
	}
    }
}
