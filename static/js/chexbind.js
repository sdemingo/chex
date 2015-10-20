

function chexInit(){

    function readUserEdited(){
	var u = $("#userEditForm").serializeObject()
	u.Tags = u.Tags.split(",").map(function(e){
	    return e.trim()
	})
	u.Tags.clean("")

	if (/\s+/.test(u.Tags.join())){
	    showErrorMessage("#userEditAlert","Las etiquetas no pueden contener espacios")
	    return
	}
	
	if ((u.username=="") || (u.email=="")){
	    showErrorMessage("#userEditAlert","Existen campos sin información")
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
	    showInfoMessage("#userEditAlert","Usuario creado con éxito")
	    $("#userEditForm").each(function(){
		this.reset();
	    })
		},function(){
		    showErrorMessage("#userEditAlert","Error al crear usuario")
		})
    })




    // Edit Users

    $("#userEditForm #userUpdateSubmit").click(function(){
	var u = readUserEdited()
	if (!u){
	    return
	}
	editUser(u,function(){
	    showInfoMessage("#userEditAlert","Usuario editado con éxito")
	},function(){
	    showErrorMessage("#userEditAlert","Error al editar usuario")
	})
    })




    $(".alert").css("visibility", "hidden");
}









function showInfoMessage(selector, text) {
    var alert = $(selector).css("visibility", "visible").addClass("alert-success").text(text)
    window.setTimeout(function() { $(selector).removeClass("alert-success").css("visibility", "hidden") }, 1500)
}

function showErrorMessage(selector, text) {
    var alert = $(selector).css("visibility", "visible").addClass("alert-danger").text(text)
    window.setTimeout(function() { $(selector).removeClass("alert-danger").css("visibility", "hidden") }, 1500)
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
