# Vision


## Purpose

namo gives scripts, dotfiles, and coding agents one way to name the things they create: scratch files, logs, PR bodies, branches, and batches of identifiers. Before namo, each of those callers assembled names by hand from a `date` timestamp and a few random words. namo turns that habit into one command whose names sort by creation time, can be recognized at a glance in a directory listing, a log, or a branch list, and can go straight into a path. Coding agents make most of the calls.


## Responsibilities and boundaries

namo is responsible for producing good names on demand, as the first principle below defines them. Its defaults follow its callers' most common call, so that call stays short.

namo serves the naming uses above. Anyone may run it, but a request counts only when it serves those uses. That someone could use a feature is not a reason to add it, and a feature that only a different kind of caller would want is out of scope even if it is small and well built.

namo is a command-line program, and invoking the command is its interface. A single interface keeps namo free to restructure its generator under its clean-break policy. The accepted cost is that every consumer depends on the namo binary being installed.


## Principles


### What makes a good name

A good namo name is sortable, recognizable, grouped by purpose, safe, and distinct within its batch. With the default stamp, names sort by creation time. The words let a person pick a name out of a listing and remember it. A prefix keeps related names together and says what they are for. Names contain only lowercase ASCII letters, digits, and hyphens unless a caller opts into trusted input. No two names from one call share a slug. Safety wins over every other quality, which is why a normalized prefix gives up capitals and punctuation. Sorting and brevity trade off at the caller's choice: a caller who needs a short handle drops the stamp with a flag, and the default keeps the stamp.


### Every run stands alone

namo is stateless: each run builds its names from its own input, the clock, and fresh randomness, and nothing carries over to the next run. Any caller, including several agents running in parallel on one machine, therefore gets a usable name without coordinating with anyone, and no state can go stale or become corrupted. The accepted cost is that namo cannot promise uniqueness across runs; it relies on the stamp and the random words to make collisions unlikely. A guarantee that would need namo to remember what it issued conflicts with this principle, however useful it would be.


### The words come from hotdiva2000 as they are

namo's slugs come from hotdiva2000, and namo keeps that vocabulary as it is, irreverent words included. Those words make a name memorable enough to pick out of a listing, which is why namo uses words at all. The accepted cost is that some names are awkward where colleagues or the public can see them, such as branch names and shared directories.


### No configuration

namo is zero-config: every option is a flag on the invocation, and no file or environment variable of namo's own can change what it does. A given invocation therefore means the same thing on every machine and inside every agent environment. The accepted cost is that there are no personal defaults; a caller who wants a particular prefix or stamp layout spells it out every time.


### Safe to run unattended

namo runs to completion on its own: it never waits on a person, another program, or the network. Scripts and agents call it unattended, usually inside command substitution, so it must finish promptly and behave the same everywhere. The accepted cost is that namo gives up any improvement to its names that would depend on something outside its own process.


### Safe output by default

Without an explicit trusted-input flag, a name contains only lowercase ASCII letters, digits, and hyphens. Callers put namo's output directly into file paths and command substitution without quoting or validating it, so the default must never produce a character that could break a path or a shell. Input that could carry unsafe characters, such as a verbatim prefix or a custom stamp layout, is allowed only behind a flag that names it as trusted. namo keeps these escape hatches even when no caller uses them, because they let a caller handle a rare case without waiting for a change to namo. The accepted cost is that a normalized prefix loses capitals and punctuation.


### Better design over compatibility

namo makes clean breaks. When a better design for the name shape or the flags appears, namo changes to it without aliases, deprecation periods, or compatibility shims. namo's callers are scripts and agent instructions maintained alongside it, so keeping them working means updating them at the time of the change, not carrying old behavior inside namo. Clean breaks serve better designs, and low usage alone is not a reason to remove a capability. The accepted cost is that a script or agent instruction that is missed fails or receives a differently shaped name.
