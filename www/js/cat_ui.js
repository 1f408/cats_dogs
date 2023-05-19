document.addEventListener('DOMContentLoaded', function(){
  let expand_url = function(url_pat, nv){
    return url_pat.replace(/\{([a-zA-Z_]+)\}/g,
    (m, n) => (nv.hasOwnProperty(n) ? encodeURIComponent(nv[n]) : ""));
  }
  let url_keys = function(url_pat){
    return [...url_pat.matchAll(/\{([a-zA-Z_]+)\}/g)].map(m => m[1]);
  }

  let set_msg = function(dom, msg){
    dom.selectAll("*").remove();
    dom.append("span").text(msg);
  }

  let set_error = function(dom, txt){
    dom.selectAll("*").remove();
    dom.append("p").text(txt);
  }

  let set_text = function(dom, txt){
    dom.selectAll("*").remove();
    dom.append("textarea")
      .property("disabled", true)
      .text(txt);
  }

  let set_csv = function(dom, txt){
    let csv = d3.csvParse(txt);

    dom.selectAll("*").remove();
    tbl = dom.append('table');
    let thead = tbl.append('thead');
    let tbody = tbl.append('tbody');

    thead.append('tr').selectAll('th')
      .data(csv.columns).enter().append('th').text(k => k);

    let rows = tbody.selectAll('tr')
      .data(csv).enter().append('tr');

    let cells = rows.selectAll('td')
      .data(row => csv.columns.map(
      col => ({column: col, value: row[col] }) ))
      .enter()
      .append('td')
      .text(d => d.value);
  }

  d3.selectAll("div.cat_ui").each(function(){
    let ui_div = d3.select(this);
    let out = ui_div.select(".cat_ui_res");
    if (out.node() == null) { return; }
    let ok_chk = ui_div.select(".cat_ui_ok");

    ui_div.select(".cat_ui_run").on("click", function(){
      let btn = d3.select(this);
      let url = btn.node().value;
      if (url == "") { return; }

      if (ok_chk.node() != null && !ok_chk.property("checked")) {
        set_msg(out, "Not confirmed.");
        return;
      }

      let v = {};
      url_keys(url).forEach(function(n){
        let o = ui_div.select(".cat_ui_in_" + n);
        if (o.node() != null) { v[n] = o.node().value; }
      });

      fetch(expand_url(url, v))
        .then(function(res){
          if (!res.ok) throw new Error(res.status + " " + res.statusText);
          let tp = res.headers.get('Content-Type');
            if (tp == null) {
            } else if (tp.match(/^text\/csv(;.*)$/)) {
              return res.text().then(function(d){ set_csv(out, d); });
            } else if (tp.match(/^text\//)) {
              return res.text().then(function(d){
                if (d.length == 0) {
                  set_msg(out, "Success!");
                } else {
                  set_text(out, d);
                }
              });
            }
            set_msg(out, "Error: unsuppoted content type.");
        })
        .catch(function(e){
          set_error(out, e);
          console.log("API error: " + e);
        });
    });
  });
});
