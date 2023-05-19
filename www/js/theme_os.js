document.addEventListener('DOMContentLoaded', function(){
  var all_theme = d3.selectAll("link.theme");
  all_theme.each(function(){
    let elm = d3.select(this);
    let sname = this.id.replace(/_theme$/, '');
    let chk = (elm.attr("disabled") == null);
    let id = sname+"_theme_rdo";
    if(sname == ""){ return; }
  });

  var def_theme = d3.select("link.theme:not([disabled])");
  function set_theme(sname){
    let must = false;
    all_theme.property("disabled", function(){
      if(this.id == (sname + "_theme")){
        return false;
      }
      return true;
    });
  }
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    set_theme('dark');
  } else {
    set_theme('light');
  }
});
