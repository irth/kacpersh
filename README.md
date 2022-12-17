# kacpersh

My friend Kacper told me about an inconvenience that sometimes happened to him - you execute a long running command, then want to do something with the output, but you forgot to redirect it to a file or something, so you have to re-run it all.

The idea is to make the shell/terminal/whatever capture all commands, and make the last command's outputs available to the user easily.

Also, Kacper sounds a bit like "capture". So, it fits.

# Proof of concept

https://gist.github.com/irth/ffba9da0a9a4f6df54f02fe06605f19c

[![asciicast](https://asciinema.org/a/vVu3yRs9bGB8Zrjs9ia301ALw.svg)](https://asciinema.org/a/vVu3yRs9bGB8Zrjs9ia301ALw)

# TODO:

## v1

- [x] launch users shell (according to SHELL/profile by default) instead of a hardcoded one
- [x] communicate over a unix socket, path to which is passed through an env variable (so that u can have more than one shell)
- [ ] generate config for zsh, so that it can be eval'ed in .zshrc easily

## v2

- [ ] in-band signalling - this will allow as to implement buffering for performance, as we don't have to keep perfect synchronisation anymore
- [ ] see if we can support bash and fish

## v3

- [ ] capture multiple commands instead of just the last one, decreasing the odds of losing data
