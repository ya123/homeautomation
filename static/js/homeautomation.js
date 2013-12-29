function dim(name, value, callback) {
	$.ajax("/dim?name="+name + "&value=" + value, {
		async: false,
		error: function() {
			// console.log("error happened");
		},

		success: function(data) {
			// console.log(data.result);
			// console.log(data);
			callback(data.result);
		}
	} );
}

$(function() {
	
	$(".switch").click(function(){
		var sw = $(this);
		var name = sw.parent().parent().attr("data-name");

		var dimmer = sw.parent().parent().find(".dimmer");
		var input = sw.parent().parent().find(".dimmer-input");

		var onAjax;
		var offAjax;
		if (dimmer.length == 1 ) {
			onAjax = function() {
				$.ajax("/dim?name="+name+"&value=100", {
					async: false,
					success: function(data) {
						sw.addClass("on");
						sw.removeClass("btn-primary");
						sw.addClass("btn-danger");
						dimmer.slider( "value", 100 );
						input.val( 100 );
					}
				});
			};

			offAjax = function() {
				$.ajax("/off?name="+name, {
					async: false,
					success: function(data) {
						sw.removeClass("on");
						sw.removeClass("btn-danger");
						sw.addClass("btn-primary");
						dimmer.slider( "value", 0 );
						input.val( 0 );
					}
				});
			}

		} else {
			onAjax = function() {
				$.ajax("/on?name="+name, {
					async: false,
					success: function(data) {
						sw.addClass("on");
						sw.removeClass("btn-primary");
						sw.addClass("btn-danger");
					}
				});
			};

			offAjax = function() {
				$.ajax("/off?name="+name, {
					async: false,
					success: function(data) {
						sw.removeClass("on");
						sw.removeClass("btn-danger");
						sw.addClass("btn-primary");
					}
				});
			};
		}
		
		if (sw.hasClass("on")) {
			offAjax();
		} else {
			onAjax();
		}
		
	});

	$(".dimmer").each(function(){
		var dimmer = $(this);
		var name = dimmer.parent().parent().attr("data-name");
		var input = dimmer.parent().find(".dimmer-input");
		var sw = dimmer.parent().parent().find(".switch");

		var swOff = function() {
			sw.removeClass("on");
			sw.removeClass("btn-danger");
			sw.addClass("btn-primary");
		}

		var swOn = function() {
			sw.addClass("on");
			sw.removeClass("btn-primary");
			sw.addClass("btn-danger");
		}

		var callback = function(serverValue) {
			dimmer.slider( "value", serverValue );
			input.val( serverValue );
			//dimmer.slider( "option", "disabled", false);
		}
	
		dimmer.slider({
			min: 0,
			max: 100,
			step: 10,
			value: input.val(),
			stop: function( event, ui ) {
				input.val( ui.value);
				dim(name, ui.value, callback );

				if (ui.value == 0) {
					swOff();
				} else {
					swOn();
				}
			}
		});

		input.change(function () {
			dimmer.slider( "value", $(this).val() );
			dim(name, $(this).val(), callback );

			if ($(this).val() == 0) {
				swOff();
			} else {
				swOn();
			}
		});
	});
});