# GiCo
A collection of tools for local processing of *Gi*t *Co*mmits

Included Tools:
- gico
- gitfetch
- mgfetch
- svg-server

## GiCo Core library

The core library for GiCo (separate from the `gico` executable of the same
name) provides utilities to read in a list of directories containing git
repositories and translating them into a list of day-aligned lists--each day
of the year points to a list structs containing metadata on each commit created
that day.

The repositories can be parsed synchronously or in parallel, using goroutines
and channels.
There is also a configurable caching interval to take advantage of memoization
as an exercise in dynamic programming, which makes the library suitable for
frequent calls such as those from a BubbleTea UI application.

It is up to more concrete implementations to decide how this data is visualized
or otherwise exposed to an end-user.
For this reason, several example tools have been provided.

## Executables
### gico

The gico binary is a tui tool that loads a list of git repos on a system
and turns it into an interactive Github-style heatmap (coloring support
included).

By default, gico uses [mg](https://github.com/taigrr/mg) to pull in a central list of
repos and parses the user's git config to extract the email and name of the
current user.
The GiCo library is used to load all the repos and translate them into a heatmap
and convert the values into a dynamically scaled, user-configurable color
pallette (see [simplecolorpalettes](https://github.com/taigrr/simplecolorpalettes)).

A searchable settings view is available to allow users to select and deselect
individual repos and authors to include in the graph.

Example GIF:


### gitfetch

Like your standard fetch program, gitfetch uses the GiCo library to parse the
git history of a repo in the current directory and print a gitgraph out to the
terminal.

### mgfetch

mgfetch uses [mg](https://github.com/taigrr/mg) to pull in a list of all git repos and
combines the heatmap lists into a single gitgraph, and prints it to the
terminal.`

### svg-server

svg-server uses GiCo in a similar way to mgfetch, by first pulling in a list of
all repos seen by [mg](https://github.com/taigrr/mg) and then generates svg files
on-the-fly depicting the resultant gitgraph.
svg-server is suitable for embedding an svg of your gitgraph onto your desktop
using conky, for example.
