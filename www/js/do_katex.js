document.addEventListener('DOMContentLoaded', function(){
  document.body.querySelectorAll(":not(pre) span.math.inline")
  .forEach(function(elem){
    katex.render(elem.textContent, elem, {
      displayMode: false,
      throwOnError: false,
    });
  });
  document.body.querySelectorAll(
      ":not(pre) span.math.display,:not(pre) pre > code.language-math")
  .forEach(function(elem){
    let text = elem.textContent;
    let p = document.createElement('p');
    katex.render(elem.textContent, p, {
      displayMode: true,
      throwOnError: false,
    });
    elem.parentElement.replaceWith(p);
  });
});
