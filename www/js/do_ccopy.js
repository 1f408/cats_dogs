if(navigator.clipboard != null){
  document.addEventListener('DOMContentLoaded', function(){
    document.body.querySelectorAll("pre > code").forEach(function(tgt){
      let btn = document.createElement('button');
      btn.classList.add('ccopy');
      let icon = document.createElement('div');
      icon.classList.add('ccopy_icon');
      btn.appendChild(icon);

      tgt.parentNode.insertBefore(btn, tgt);
      btn.addEventListener("click", function(){
        navigator.clipboard.writeText(tgt.textContent);
      });
    });
  });
}
