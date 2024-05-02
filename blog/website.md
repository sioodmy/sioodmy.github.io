1714668335

# My adventures overengineering simple website

So I was recently working on my personal website and encountered a problem. Up until now my site was a fork of some random's github blog with a few things changed. Of course, this was only supposed to be a temporary solution, so I started making mine from scratch.

## Solution 1 - static HTML

At first I thought a simple static site would be more than enough for me. I liked playing with the UI of the site, I liked the minimalist style without unnecessary JS. However, I quickly got to the point that if I ever needed to update some text on the site I would have to write html by hand. This wouldn't actually be a big problem considering how often I do it, but my ego forced me to improve. And here we come to my adventure with blog generators.

## Solution 2 - Zola/Hugo

I decided that writing my blog using an existing generator would not be a great challenge and would be "correct enough" to use it.
I went back to my fork to analyze again how it all works. I wasted 2h of my life converting my site to Zola, only to come to the conclusion that it is "too much". Seriously, I felt overwhelmed, too many features I would never use.

Then I remembered the existence of Hugo, Luke Smith made a video about it, so it should be pretty minimalist. Oh boy, was I wrong. After reading a piece of documentation, I came to a similar conclusion as with Zola. I knew I had to do it my way

## Solution 3 - Nix/Shell

I decided that as a Nix guy, I should write the whole thing in pure Nix (with the small help of shell scripts). Thus, I spent a couple of hours sifting through comically badly written documentation and writing my scripts. Eventually I arrived at something that worked. It seemingly worked, but it gave the impression of poorly written scripts glued together with a duck tape. That's why I quickly gave up on this solution

## Solution 4 - Nix/Go

And here we come to the "best" solution I have come to. I wrote my own generator in go, specifically for this site. I designed it to do exactly three things:

-   generates a subpage with a list and short descriptions of my projects
-   generates all the blog posts
-   sorts and arranges the posts in an index

...meaning everything I would need to do manually, nothing more, nothing less.

From now on all my blogposts look like this:

```
1714668335
# Title

```

First two lines are reserved for post meta data, first line is a date in the standard UNIX format and the second is the title. Pretty simple.

I also wrote a Nix flake to automate the build process. It provides two packages: generator and the website.
Website package is just a wrapper around the generator, it takes it as a build input, calls it, and copies generated files over to $out.

This way I can just call the website package in `services.nginx.virtualHosts.<name>.root`. Super convinient.

## Conclusion

Was it all worth it?

no.

I could have just used Zola. But doing it all was fun, so no regrets.
