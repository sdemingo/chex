

function chexInit(){

    // New Users

    $("#userNewForm #userNewSubmit").click(function(){
	var u = $("#userNewForm").serializeObject()
	u.Tags = u.Tags.split(",").map(function(e){
	    return e.trim()
	})
	u.Tags.clean("")

	if (/\s+/.test(u.Tags.join())){
	    showErrorMessage("#userNewAlert","Las etiquetas no pueden contener espacios")
	    return
	}
	

	if ((u.username=="") || (u.email=="")){
	    showErrorMessage("#userNewAlert","Existen campos sin información")
	    return
	}

	addUser(u,function(){
	    showInfoMessage("#userNewAlert","Usuario creado con éxito")
	    $("#userNewForm").each(function(){
		this.reset();
	    })
	},function(){
	    showErrorMessage("#userNewAlert","Error al crear usuario")
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
