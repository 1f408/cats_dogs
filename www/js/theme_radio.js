document.addEventListener('DOMContentLoaded', function(){
  let all_theme = d3.selectAll("link.theme");
  let theme_form = d3.select("#style_config").append("form").classed("theme_frm", true);

  let mode = "";
  let mode_div;
  all_theme.each(function(){
    let elm = d3.select(this);
    let sname = this.id.replace(/_theme$/, '');
    let v = sname.match(/^(dark|light)([0-9A-Z])$/);
    if(v == null){ return; }

    let mname = v[1];
    let tname = v[2];
    let chk = (elm.attr("disabled") == null);
    let id = sname+"_theme_rdo";
    if(sname == ""){ return; }

    if(mode != mname){
      mode_div = theme_form.append("div");
      mode = mname;
      mode_div.append("span").text(mode+":");
    }
    mode_div.append("input").classed("theme_in", true).attr("id", id)
      .attr("type", "radio").attr("name", "theme").attr("value", sname)
      .property("checked", chk);
    mode_div.append("label").text(tname).attr("for", id);
  });

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
    d3.select(".theme_frm input.theme_in[name=\"theme\"][value="+sname+"]")
      .property("checked", true);
    set_theme(sname);
  }
  change_theme(localStorage.getItem("ColorStyle"));

  d3.selectAll(".theme_frm input.theme_in[name=\"theme\"")
    .on("change", function(){ set_theme(this.value); });
});
