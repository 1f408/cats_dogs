document.addEventListener('DOMContentLoaded', function(){
  function open_extlink(ev){
    ev.preventDefault();

    let url = this.getAttributeNS('http://www.w3.org/1999/xlink', 'href');
    if (url == null) {
      url = this.href;
    }

    if (url != null) {
      let w = window.open("");
      w.document.open();
      w.document.write('<meta http-equiv="refresh" content="0;url='+url+'">');
      w.document.close();
    }

    ev.stopPropagation();
  }
  document.body.querySelectorAll("a").forEach(function(el) {
    if(el.href == ""){ return; } 
    let l = new URL(el.href);
    if(window.location.origin == l.origin){ return; }

    el.setAttribute("target", "_blank");
    el.addEventListener("click", open_extlink);
  });
});
