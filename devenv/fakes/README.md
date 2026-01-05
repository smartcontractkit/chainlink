# DataFeeds Fakes

This server contains DataFeeds specific fakes

To run it locally and develop (you need a fine-grained GH approved token here)
```
task build -- $tag
task run
```

Test it
```
curl -v "http://localhost:9111/static-fake"
curl -v "http://localhost:9111/dynamic-fake"
```
Publish it
```
task publish -- $tag
```