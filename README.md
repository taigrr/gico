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

Future:
 - possibly use a client-server model to make the execution of the git hook much faster (no need to load, parse, save the file every execution)
