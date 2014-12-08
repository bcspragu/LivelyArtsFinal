$(function() {
  $('.send').click(function() {
    var text = $('.main-input').val();
    if (text != "") {
      text = text.replace(/[\.,-\/#!$%\^&\*;:{}=\-_`~()]/g,"");
      text = text.replace(/\s{2,}/g," ");
      $.post("/input", {words: text});
    }
  });
});
