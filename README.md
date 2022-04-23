# nhkpod

A podcast server providing Radiru programs written in golang

**All data processed by this program should be handled in accordance with the terms of use for Radiru https://www.nhk.or.jp/radio/info/kiyaku.html**

## Usage

1. Visit [Radiru ProgramList](https://www.nhk.or.jp/radioondemand/json/index_v3/index.json), copy `site_id` for the
   program you want to follow.

2. Edit `conf.yml` and declare podcast. Specify the copied `site_id` as the `id` of the podcast. (Optional) If you
   specify `corner_id`, the corner will be served as an independent podcast. Otherwise, all the corners in the site_id are
   served as a single podcast.

```
podcasts:
  - id: 1633
  - id: F295
    corner_id: 29
```

This will make podcasts `http://localhost:8080/audio/f295_29/feed.rss` and `http://localhost:8080/audio/f295_29/feed.rss`

3. Run cmd/nhkpod/main.go by `make build && ./nhkpod` or your preferred method to start the server

4. Initially, it will download all the available audio files associated with the `site_id` and `corner_id` you've specified.

5. Register url `http://<host>:8080/audio/<site_id>/feed.rss`
   or `http://<host>:8080/audio/<site_id>_<corner_id>/feed.rss` on your podcast client.

## Environment variables

| Key               | Description                                                         | Default    |
|-------------------|---------------------------------------------------------------------|------------|
| NHKPOD_SCHEDULE   | schedule for audio file download & podcast feed update (cron style) | 35 * * * * | 
| NHKPOD_LOG_FILE   | log file                                                            | log.log    |
| NHKPOD_AUDIO_DIR  | audio & podcast feed directory                                      | audio      |
| NHKPOD_CONF_PATH  | configuration file path                                             | conf.yml   |
| NHKPOD_PORT       | port for podcast server                                             | 8080       |
