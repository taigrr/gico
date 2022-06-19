# GiCo

Goals:
1. Provide a git post-commit hook binary that allows you to track commit frequency locally
  - requires a lock file
1. Save a history locally to a JSON file
  - requires a lock file
1. Provide functionality to display the data:
  - bubbletea interactive, scrollable graph
  - one-shot printed graph for a date range (default 1 year)
  - print out commit count
    - today
    - this week
    - this month
    - this year
1. Provide an import method to add old repos to history
1. Expose parsing / counting functions to other libraries



Notes:
- use a client-server model to make the execution of the git hook much faster (no need to load, parse, save the file every execution)
 - An appropriate communication method for an application like this would be Unix Domain Sockets, to reduce network overhead.
   However, a goal of the program is to serve generated graph images over an API, so since we already need network communication for that, it makes sense to reuse it rather than listen on two interfaces at once.
- env var PWD is set to git repo base even when in a subfolder (subfolder is kept in GIT_PREFIX)
- Author Email is kept in  GIT_AUTHOR_EMAIL=<string>
- Commit message file is in argv[1] for hook invocation, read in the file (and ignore empty lines + lines starting with `#` to get the message
- use the current date for the 

