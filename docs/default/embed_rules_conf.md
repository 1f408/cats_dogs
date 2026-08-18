# 埋め込み処理のデフォルト設定

```
video_ext = ["mp4", "m4v", "webm"]
audio_ext = ["mp3", "m4a", "wav", "wave", "flac"]

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="query"
query="v"
path="/watch"
player="https://www.youtube.com/embed/$0"
allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
referrerpolicy="strict-origin-when-cross-origin"

[[iframe]]
site_id="youtube"
host="youtu.be"
type="path"
path="/"
player="https://www.youtube.com/embed/$0"
allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
referrerpolicy="strict-origin-when-cross-origin"

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="path"
path="/embed"
player=""
allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
referrerpolicy="strict-origin-when-cross-origin"

[[iframe]]
site_id="vimeo"
host="vimeo.com"
type="path"
path="/"
player="https://player.vimeo.com/video/$0"

[[iframe]]
site_id="vimeo"
host="player.vimeo.com"
type="path"
path="/video"
player=""
```
