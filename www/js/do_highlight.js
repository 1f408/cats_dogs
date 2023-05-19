document.addEventListener('DOMContentLoaded', function(){
  document.body.querySelectorAll(":not(pre) pre > code")
  .forEach(function(elem){
    if(elem.classList.lentgh == 0) { return; }
    let lang = null;
    for (let cn of elem.classList){
      if (cn.startsWith('language-')) {
        lang = cn.replace(/^language-/, '');
        break;
      }
    }

    if (lang != null && hljs.getLanguage(lang) != null){
      hljs.highlightElement(elem);
      elem.classList.add("plaintext");
    }
  });
});
