$(function() {
  $('.send').click(function() {
    var text = $('.main-input').val();
    if (text != "") {
      text = text.replace(/[\.,-\/#!$%\^&\*;:{}=\-_`~()]/g,"");
      text = text.replace(/\s{2,}/g," ");
      text = text.toLowerCase();
      $.post("/input", {words: text});
    }
  });
});
