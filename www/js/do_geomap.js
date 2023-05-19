document.addEventListener('DOMContentLoaded', function(){
  let cnt = 0;
  function gen_id() {
    let id = "geomap-" + String(cnt);
    cnt += 1;
    return id;
  }

  document.body.querySelectorAll(":not(pre) pre > code.language-geojson")
  .forEach(function(src){
    let text = src.textContent;
    if(text == ""){ return; }
    let json = JSON.parse(src.textContent);
    if(json == null){ return; }
    let id = gen_id();
    let tgt = document.createElement('div');
    tgt.id = id;
    tgt.classList.add('geomap');
    src.parentElement.replaceWith(tgt);
    let map = L.map(id);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', 
  {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }
).addTo(map);
    let geoLayer = L.geoJSON().addTo(map);
    geoLayer.addData(json);
    map.fitBounds(geoLayer.getBounds());
    map.on('resize', function(){
      map.fitBounds(geoLayer.getBounds());
    });
  });

  document.body.querySelectorAll(":not(pre) pre > code.language-topojson")
  .forEach(function(src){
    let text = src.textContent;
    if(text == ""){ return; }
    let json = JSON.parse(src.textContent);
    if(json == null){ return; }

    let id = gen_id();
    let tgt = document.createElement('div');
    tgt.id = id;
    tgt.classList.add('geomap');
    src.parentElement.replaceWith(tgt);
    let map = L.map(id);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    L.TopoJSON = L.GeoJSON.extend({
      addData: function (data) {
        let geojson, key;
        if (data.type === "Topology") {
          for (key in data.objects) {
            if (data.objects.hasOwnProperty(key)) {
              geojson = topojson.feature(data, data.objects[key]);
              L.GeoJSON.prototype.addData.call(this, geojson);
            }
          }
          return this;
        }
        L.GeoJSON.prototype.addData.call(this, data);
        return this;
      }
    });
    L.topoJson = function (data, options) {
      return new L.TopoJSON(data, options);
    };

    let geoLayer = L.topoJson().addTo(map);
    geoLayer.addData(json);
    map.fitBounds(geoLayer.getBounds());
    map.on('resize', function(){
      map.fitBounds(geoLayer.getBounds());
    });
  });

});
