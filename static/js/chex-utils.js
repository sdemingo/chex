
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
