

function chexInit(){

    function readUserEdited(){
	var u = $("#userEditForm").serializeObject()
	u.Tags = u.Tags.split(",").map(function(e){
	    return e.trim()
	})
	u.Tags.clean("")

	if (/\s+/.test(u.Tags.join())){
	    showErrorMessage("Las etiquetas no pueden contener espacios")
	    return
	}
	
	if ((u.Name=="") || (u.Maill=="")){
	    showErrorMessage("Existen campos sin información")
	    return
	}
	return u
    }



    // New Users

    $("#userEditForm #userNewSubmit").click(function(){
	var u = readUserEdited()
	if (!u) {
	    return
	}
	addUser(u,function(){
	    showInfoMessage("Usuario creado con éxito")
	    $("#userEditForm").each(function(){this.reset()})
		}
		,function(){
		    showErrorMessage("Error al crear usuario")
		})
    })




    // Edit Users

    $("#userEditForm #userUpdateSubmit").click(function(){
	var u = readUserEdited()
	if (!u){
	    return
	}
	editUser(u,function(){
	    showInfoMessage("Usuario editado con éxito")
	},function(){
	    showErrorMessage("Error al editar usuario")
	})
    })


    // List Users
    function loadListLabels(selector){
	getTags(function(data){
	    $.each(data,function(i,e){
		$(selector).append("<a href=\"#\" class=\"label label-default\">"+e+"</a>")
	    })
		})
    }
/*
    loadListLabels("#usersTags")
    $("#usersTags").on("click","*",function(e){
	$(this).toggleClass("label-primary")
    })
    $("#usersTags").on("click",function(e){
	tags=[]
	$("#usersTags").find(".label-primary").each(function(){
	    tags.push($(this).html())
	})
	    
	    if (tags.length>0){
		console.log("Hago el ajax con : "+tags)
		console.log(JSON.stringify(tags))
		getUsers(tags,function(data){
		    if (data){
			console.log("encontrados "+data.length)
		    }else{
			console.log("no encontro usuarios")
		    }
		})
	    }
    })
*/
    


    $(".alert").css("visibility", "hidden");
}









function showInfoMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-success").text(text)
    window.setTimeout(function() { $("#infoPanel").removeClass("alert-success").css("visibility", "hidden") }, 1500)
}

function showErrorMessage(text) {
    var alert = $("#infoPanel").css("visibility", "visible").addClass("alert-danger").text(text)
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
