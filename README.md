# yaml-diffs
TODO:
- mention line number also while showing diff
- handle cases where configs can be like:
```yaml
  plugins:
      - jekyll-avatar
      - jekyll-feed
      - jekyll-mentions
      - jekyll-redirect-from
      - jekyll-seo-tag
      - jekyll-sitemap
      - jemoji
```

- handle cases where config values can be array:
```yaml
defaults:
  -
    scope:
      path: "_docs"
      type: "docs"
    values:
      layout: "docs"
  -
    scope:
      path: "_posts"
      type: "posts"
    values:
      layout: "news_item"
      image: /img/twitter-card.png
 
```
