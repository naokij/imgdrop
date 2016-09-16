$(function()
{	
    $.validate({
	  form: '#login-form',
	  modules: 'html5',
	  errorElementClass: 'uk-form-danger',
	  errorMessageClass: 'uk-text-danger'
	});
});