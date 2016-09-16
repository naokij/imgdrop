$(function()
{
	// change locale and reload page
	$(document).on('click', '.lang-changed', function(){
		var $e = $(this);
		var lang = $e.data('lang');
		$.cookie('lang', lang, {path: '/', expires: 365});
		window.location.reload();
	});

});	