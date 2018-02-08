# LogView

See [/docs/design/CIServer.md](CIServer.md) for design notes.

## Port Forward

This could be improved, but it's not critical for now because only YourBase devs need to test the logviewer locally.

```
export POD_NAME=$(kubectl get pods --namespace logging -l "app=elasticsearch,component=client" -o jsonpath="{.items[0].metadata.name}")
kubectl -n logging port-forward $POD_NAME 9200:9200
```

## ElasticSearch

*UPDATE* The stuff below isn't relevant right now because we are simply searching for a build ID which is a uuid without hyphens and doesn't get affected by this. Keeping here for later reference if the same problem does appear again.

We can't search for hyphenized words by default because they can get split out. This affects the search for branch=yb-test for example.

Possible solutions:

- change the dynamic index template on ES so these fields don't get analyzed, or are a different token is used for the analysis
- change dynamic index template on ES to add multi-fields, and add ".raw" for these fields so we can search them without tokenization
- change the query itself and search for "branch:yb" + "branch:test" instead. Works OK but not a perfect match.
- sidestep the issue completely and use a random simple string for each build (we're doing this now)ba

#### Index Settings

We have to change the `branch` index otherwise we can't search for hyphenized words.

```
---
settings:
  analysis:
    analyzer:
      german_analyzer:
        tokenizer: standard
        filter: [standard, stop, lowercase, asciifolding, german_stemmer]
    filter:
      german_stemmer:
        type: stemmer
        name: light_german
mappings:
  my_type:
    properties:
      text:
        type: "string"
        analyzer: "german_analyzer"
```

curl -XPOST localhost:9200/my_index --data-binary @settings.yaml

Full process:
gre
```
# delete any previous versions of that index
curl -XDELETE $ES_HOST/$INDEX
# recreate using our settings
curl -XPOST $ES_HOST/$INDEX --data-binary @settings.yaml
# put our data in
curl -XPOST $ES_HOST/$INDEX/_bulk --data-binary @docs.json
# any additional setup you have to do
#....
```

Reference: http://asquera.de/blog/2013-07-10/an-elasticsearch-workflow/

Logstash indexes are created every day, so we need `dynamic_templates`:

```
With dynamic_templates, you can take complete control over the mapping that is generated for newly detected fields. You can even apply a different mapping depending on the field name or datatype.
```

https://www.elastic.co/guide/en/elasticsearch/guide/current/custom-dynamic-mapping.html#dynamic-templates

At some point this
procedure will be scripted and deployed automatically (possible via a helm chart).

### Kibana Dev Tools inputs tips.

Check the type of an index:

```
GET logstash-2018.01.26/_analyze
{
  "field": "branch",
  "text": "log-view"
}
```

Example problematic output (note how the branch name got split)

```
{
  "tokens": [
    {
      "token": "log",
      "start_offset": 0,
      "end_offset": 3,
      "type": "<ALPHANUM>",
      "position": 0
    },
    {
      "token": "view",
      "start_offset": 4,
      "end_offset": 8,
      "type": "<ALPHANUM>",
      "position": 1
    }
  ]
}
```

