---
date: 2025-12-30
draft: true
title: "My Custom Go Blog Stack"
description: "How I arrived at the stack for this blog"
---

I have had some markdown blogs burning a hole in my obsidian vault for a while now and recently actually got some time to  add them into this website. And how do you create a static blog site from markdown files? [Hugo! Hugo! Hugo!](https://gohugo.io/) I hear you all shouting... 

Well I tried it and it wasn't for me. I already had build a (fairly rudimentary) website and just wanted a "Blog" section. So I tried using Hugo to just generate this section and it was bloated and felt clunky. Using another theme when I'd already made my own out of html and css felt inconsistent and didn't feel *mine*. "I'm a developer", I thought to myself smugly, "I can build this" (I mean how many of you aren't guilty of this thought)...

So I built my own. Here. You're reading from it now. It's not perfect but it's definitely *mine*.

I had a few requirements going into this:

### Requirements:
- Ideally I can write plain html, css with no js framework
- As little js as possible
- Easy to write new blog posts (no writing plain html and wrapping evertyhing in `<p>` tags as I go)
- I'd like to write as much in Go as possible. 
- Needs to be low cost
- Did I mention no JS?

### Nice to Haves:
- Easily extendable to have comments and likes (perhaps using twitter maybe?)
- A system for counting views
- Potentially incorporate a little sprinkling of HTMX (although for a static site I expect this is a pipe dream)

### Nice Not to Haves:
- JavaScript... (OK I'll stop now)

## So how did I do it?

I know that Hugo is written in Go so my first thought was what is Hugo doing? After some initial digging I found that it's using [Goldmark](https://github.com/yuin/goldmark) to parse markdown files. Excellent... If it works for them it'll work for me. 

<figure>
  <img src="/Images/hugo-github-path.png" alt="Hugo Github Path">
  <figcaption>Found it!</figcaption>
</figure>

Turns out that using `goldmark` is really quite simple. It's just a matter of creating a `goldmark.Markdown` instance and then calling `Convert` on it.

```go
renderer := goldmark.New()
err := renderer.Convert(source, &buf)
if err != nil {
  panic(err)
}
```

Nice and simple. However what's really great about goldmark is that it's completely extensible via plugins. You want code highlighting, frontmatter metadata parsing or latex parsing? No problem

```go
renderer := goldmark.New(
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		),
		meta.Meta,
		latex.NewLatex(),
	),
)
err := renderer.Convert(source, &buf)
if err != nil {
	panic(err)
}
```

And to take it even further, if you have some experience walking AST's you can quite easily write a custom parser for your own implementations. However I will leave that up to you to investigate further.



> [!NOTES]
> - [example blog](https://bits.logic.inc/p/ai-is-forcing-us-to-write-good-code)
> - subscriptions?
> - I like to write code... So let me write the code
> - Would I do it like this again?
> - I don't really hate JS - It's just fun to bash
