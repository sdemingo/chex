

function chexInit(){

    $("#userNewForm #userNewSubmit").click(function(){
	var u = $("#userNewForm").serializeObject()
	if ((u.username=="") || (u.email=="")){
	    showErrorMessage("#userNewAlert","Existen campos sin información")
	    return
	}

	addUser(u,function(){
	    showInfoMessage("#userNewAlert","Usuario creado con éxito")
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

