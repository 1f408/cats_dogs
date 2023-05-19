document.addEventListener('DOMContentLoaded', function(){
  document.body.querySelectorAll("pre > code.language-whitespace")
  .forEach(function(elem){
    elem.innerHTML = elem.textContent.replace(/ /g, "<span class=\"space\">\u2423</span>");
    elem.classList.add("plaintext");
    elem.classList.remove("language-whitespace");
  });
});
