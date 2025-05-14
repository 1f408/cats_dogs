document.addEventListener('DOMContentLoaded', function(){
  let all_theme = d3.selectAll("link.theme");

  let def_theme = d3.select("link.theme:not([disabled])");
  function set_theme(sname){
    let must = false;
    all_theme.property("disabled", function(){
      if(this.id == (sname + "_theme")){
        must = true;
        return false;
      }
      return true;
    });
    if(must){
      localStorage.ColorStyle = sname;
    } else {
      localStorage.removeItem("ColorStyle");
      def_theme.property("disabled", false);
    }
  }
  function change_theme(sname){
    let theme_form = d3.select("input.theme[value="+sname+"]")
      .property("checked", true);
    set_theme(sname);
  }
  d3.selectAll("input.theme")
    .on("change", function(){ set_theme(this.value); });
  change_theme(localStorage.getItem("ColorStyle"));
});
