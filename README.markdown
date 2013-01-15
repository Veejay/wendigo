Wendigo
=======

Some form of grep with web UI.

What I want this program to do *eventually*:

User has a web UI with a search box, inputs a search term, hits enter and gets a list of code blocks of all the files
in a given directory that contain that search term with additional lines around the given line.

The code blocks should be syntax-higlighted à la gist (yet not as ugly as gists if possible) and the whole thing should leverage 
Go's concurrency capabilities if possible (i.e. trying to search in multiples files at a time and returning results to the UI as soon
as available) because this s**t really gets me going. I have next to NO idea how realistic or even 
sensical that is though. Let's say it makes sense for me right now.

Nice to have: 
-------------

+ Defining files to ignore
+ Defining a template for the resulting code blocks
+ Defining the number of lines to display 
+ Defining filetypes to filter
