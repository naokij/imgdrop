function showImage(url, thumb){
    $("#upload-drop").addClass("uk-hidden");
    $("#upload-result").removeClass("uk-hidden");
    $("#uploaded-img").attr("src", thumb);
    $("#upload-url").attr("value",url);
}

$(function(){

    var progressbar = $("#progressbar"),
        bar         = progressbar.find('.uk-progress-bar'),
        settings    = {

        action: '/upload', // upload url

        allow : '*.(jpg|jpeg|gif|png)', // allow only images

        filelimit: 1,

        loadstart: function() {
            $("#upload-error").addClass("uk-hidden");
            bar.css("width", "0%").text("0%");
            progressbar.removeClass("uk-hidden");
        },

        progress: function(percent) {
            percent = Math.ceil(percent);
            bar.css("width", percent+"%").text(percent+"%");
        },

        allcomplete: function(response) {

            bar.css("width", "100%").text("100%");

            setTimeout(function(){
                progressbar.addClass("uk-hidden");
            }, 250);
            try{
                result = JSON.parse(response);
                if (!result.Success) {
                    $("#upload-error").removeClass("uk-hidden");
                    $("#upload-error").text(result.Error);
                }else{
                    showImage(result.URL, result.Thumb)
                }
            }catch(e){
                $("#upload-error").removeClass("uk-hidden");
                $("#upload-error").text("Invalid response from server:"+response);
            }
        }
    };

    var select = UIkit.uploadSelect($("#upload-select"), settings),
        drop   = UIkit.uploadDrop($("#upload-drop"), settings);
});