# MdRest

Digitize your Markdown documents by Json and html for restful app

## Features

* md files to html and json
* static search index support
* auto generate md file summary
* auto converts a relative path to an absolute path 
* watch file changes


## Quick Preview (HTML output)

```
Before                                 After

sample_docs                            data
├── Simple Article.md                  ├── mdrest_index.json
├── YAML Article.md                    ├── mdrest_search_index.json
├── _DraftArticle.md                   ├── mdrest_sitemap.json
├── first dir                          ├── first dir
│   ├── Hello word.md                  │   └── hello word.html
│   └── img                            ├── second dir
│       └── logo.png                   │   └── hello word.html
└── second\ dir                        ├── simple article.html
    └── Hello\ word.md                 └── yaml article.html

```

## Quick Preview (json output)

in config.json, set OutputType "json"

```
Before                                 After

sample_docs                            data
├── Simple Article.md                  ├── mdrest_index.json
├── YAML Article.md                    ├── mdrest_search_index.json
├── _DraftArticle.md                   ├── mdrest_sitemap.json
├── first dir                          ├── first dir
│   ├── Hello word.md                  │   └── hello word.json
│   └── img                            ├── second dir
│       └── logo.png                   │   └── hello word.json
└── second\ dir                        ├── simple article.json
    └── Hello\ word.md                 └── yaml article.json

```

## JSONs


### simple article.json

Before

```
Simple Article.md

last modified: 2017-01-17T10:42:14+08:00

file content:

# This is a Simple Article

This is a Simple Article

```

After

```json
{
  "date": "2017-01-17T10:42:14+08:00",
  "html": "\u003ch1 id=\"this-is-a-simple-article\"\u003eThis is a Simple Article\u003c/h1\u003e\n\n\u003cp\u003eThis is a Simple Article\u003c/p\u003e\n",
  "location": "simple article",
  "title": "Simple Article"
}
```

### yaml article.json

Before

```
YAML Article.md

last modified: 2017-01-17T10:42:14+08:00

file content:

---
title: Hello world
author: leenanxi
type: Page
tags: [golang, hello]
date: 2016-12-29
draft: true
---

# This is a article include YAML

This is content
```

After

```json
{
  "author": "leenanxi",
  "date": "2016-12-29T00:00:00Z",
  "draft": true,
  "html": "\u003ch1 id=\"this-is-a-article-include-yaml\"\u003eThis is a article include YAML\u003c/h1\u003e\n\n\u003cp\u003eThis is content\u003c/p\u003e\n",
  "location": "yaml article",
  "tags": [
    "golang",
    "hello"
  ],
  "title": "Hello world",
  "type": "Page"
}
```

### mdrest_sitemap.json (config.SiteMapDeep is 2)

```json 
[
  {
    "title": "Simple Article",
    "location": "simple article"
  },
  {
    "title": "Hello world",
    "location": "yaml article"
  },
  {
    "title": "first dir",
    "children": [
      {
        "title": "Hello word",
        "location": "first dir/hello word"
      }
    ]
  },
  {
    "title": "second dir",
    "children": [
      {
        "title": "Hello word",
        "location": "second dir/hello word"
      }
    ]
  }
]
```
### mdrest_search_index.json

```json
[
  {
    "location": "first dir/hello word",
    "text": "Hello This is title This can be a blog engine that not limit you.\nSimple Image Simple Internal MD files you can use this Hello word, to links to other md documents\n",
    "title": "Hello word"
  },
  {
    "location": "second dir/hello word",
    "text": "Hello This is title This can be a blog engine that not limit you.\nThis is some image\n",
    "title": "Hello word"
  },
  {
    "location": "simple article",
    "text": "This is a Simple Article This is a Simple Article\n",
    "title": "Simple Article"
  },
  {
    "location": "yaml article",
    "tags": [
      "golang",
      "hello"
    ],
    "text": "This is a article include YAML This is content\n",
    "title": "Hello world"
  }
]
```

## How to Use

```bash
go get github.com/russross/blackfriday
go get github.com/ghodss/yaml
cd git.tiup.us/linx/mdjson
make
cd build
./mdrest
```
then check out  sample_docs/web/data/*
