# pos [![Build Status](https://travis-ci.org/Urethramancer/pos.svg?branch=master)](https://travis-ci.org/Urethramancer/pos)
Invoice tool for the command line.

## What
This tool works either in interactive (shell) mode, or as a set of tool commands (like Git or Go) to manage clients, jobs, tasks and invoices for a tiny sole proprietorship.

## How
- Prepare a Postgres database server and a user with privileges to create database.
- Run `pos setup` to set up the connection.
- Run `pos` and figure out the tool commands from help
- Run `pos sh` to open a session to your database and figure it out in interactive mode.

## Status: Non-functional
It's still being built. Clients and contacts are all there is right now, and tool commands are not implemented at all apart from setup. It's mostly available this early to allow Travis CI. Hence the sparse documentation. Many things are still in flux, but the end-goal is to get me a tool to produce printable invoices and keep track of jobs. If others can use it, great for them.
