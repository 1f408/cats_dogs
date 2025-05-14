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

  function is_anchor(l){
      return window.location.origin == l.origin
      && window.location.pathname == l.pathname
      && window.location.search == l.search;
  }

  document.body.querySelectorAll("body > .contents a").forEach(function(el) {
    if(el.href == ""){ return; } 
    let l = new URL(el.href);

    if(is_anchor(l)){
      el.classList.add("anchor_link");
    }

    if(window.location.origin == l.origin){ return; }

    el.setAttribute("target", "_blank");
    el.addEventListener("click", open_extlink);
  });
});
window.addEventListener('beforeprint', function(){
  document.body.querySelectorAll("body > .contents a:not(.anchor_link)").forEach(function(el) {
    let href = el.getAttribute('href');
    if(href == null){ return; }

    el.setAttribute('save_href', href);
    el.removeAttribute('href');
  });
});
window.addEventListener('afterprint', function(){
  document.body.querySelectorAll("body > .contents a:not(.anchor_link)").forEach(function(el) {
    let href = el.getAttribute('save_href');
    if(href == null){ return; }

    el.setAttribute('href', href);
    el.removeAttribute('save_href');
  });
});
