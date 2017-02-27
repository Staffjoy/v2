# Contributing to Staffjoy

## Community and contact information

- [Pull requests](https://github.com/staffjoy/v2/pulls)
- Support and bug report email: [help@staffjoy.com](mailto:help@staffjoy.com)

## Testing

We run static code analysis wherever possible, including linting and
formatting checks. We also extensively monitor all systems and 
staging/production errors. 

We believe in strategically testing important logic, but that achieving 100% 
coverage is not always worth the time.

## Commit message format

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```
scripts: add the test-cluster command

this uses tmux to setup a test cluster that you can easily kill and
start for debugging.

Fixes #38
```

The format can be described more formally as follows:

```
<subsystem>: <what changed>
<BLANK LINE>
<why this change was made>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various
git tools.