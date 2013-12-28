function dim(name, value, callback) {
	$.ajax("/dim?name="+name + "&value=" + value, {
		error: function() {
			console.log("error happened");
		},

		success: function(data) {
			console.log(data.result);
			console.log(data);
			callback(data.result);
		}
	} );
}

$(function() {
	$(".dimmer").each(function(){
		var dimmer = $(this);
		var input = dimmer.parent().find(".dimmer-input");

		var callback = function(serverValue) {
			dimmer.slider( "value", serverValue );
			input.val( serverValue );
		}
	
		dimmer.slider({
			min: 0,
			max: 100,
			value: input.val(),
			slide: function( event, ui ) {
				input.val( ui.value);
				dim("", ui.value, callback );
				// ?name=gerade gemacht

			}
		});

		input.change(function () {
			dimmer.slider( "value", $(this).val() );
			dim("", $(this).val(), callback );
		});
	});
});