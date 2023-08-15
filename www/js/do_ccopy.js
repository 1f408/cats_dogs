if(navigator.clipboard != null){
  document.addEventListener('DOMContentLoaded', function(){
    document.body.querySelectorAll("pre > code").forEach(function(tgt){
      let btn = document.createElement('button');
      btn.classList.add('ccopy');
      tgt.parentNode.insertBefore(btn, tgt);
      btn.addEventListener("click", function(){
        navigator.clipboard.writeText(tgt.textContent);
      });
    });
  });
}
