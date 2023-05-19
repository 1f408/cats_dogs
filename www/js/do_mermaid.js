document.addEventListener('DOMContentLoaded', function(){
  function get_var(n) {
    return window.getComputedStyle(document.querySelector(":root"))
      .getPropertyValue(n).trim();
  }
  let is_dark = (get_var('--color-schame') == 'dark');
  const config = {
    theme: (is_dark ? 'dark' : 'default'),
    startOnLoad: false,
    darkMode: is_dark,
  };

  const mermaidAPI = mermaid.mermaidAPI;
  mermaidAPI.initialize(config);
  let cnt = 0;
  function gen_id() {
    let id = "mermaid-" + String(cnt);
    cnt += 1;
    return id;
  }
  let close_dt = Array.from(document.body.querySelectorAll("details"))
    .filter(e => { return e.open == false; }); 
  close_dt.forEach(function(e){ e.open = true; });
  document.body.querySelectorAll(":not(pre) pre > code.language-mermaid")
  .forEach(function(src){
    try {
      let id = gen_id();
      let tgt = document.createElement('p');
      tgt.classList.add('mermaid');
      let text = src.textContent;
      src.parentElement.replaceWith(tgt);
      tgt.innerHTML = mermaidAPI.render(id, text, undefined, tgt);
    } catch (error) {
      console.log("mermaid:", error);
    }
  });
  close_dt.forEach(function(e){ e.open = false; });
});
