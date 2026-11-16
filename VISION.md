# Vision


## Purpose

namo exists so that Javier's own tooling can name things the same way everywhere. Before namo, the author's scripts, dotfiles, and coding-agent instructions each assembled names by hand from a `date` timestamp and a few random words. namo turns that habit into one command: a name that sorts by creation time, that a person can recognize at a glance in a directory listing, a log, or a branch list, and that a script or agent can drop straight into a path without thinking about it. Its users are Javier, the scripts and dotfiles Javier maintains, and the coding agents Javier directs, which make most of the calls.


## Responsibilities and boundaries

namo is responsible for producing good names on demand: choosing the words, rendering the stamp, normalizing the prefix, and keeping names in one batch distinct. The command is namo's only interface. Its naming code stays internal, and other tools, including Go tools, get names by running the binary. Keeping a single interface keeps namo free to restructure the generator under its clean-break policy. The accepted cost is that every consumer depends on the namo binary being installed. Its defaults follow what its author does most, so the common call stays short.

namo serves its author's workflow first. Other people are welcome to use it, but a request counts only when it matches how Javier and Javier's tooling name things. That someone could use a feature is not a reason to add it, and a feature that only outside users would want is out of scope even if it is small and well built.

namo does not remember anything between runs, so it does not guarantee that a name is unique across runs, keep a registry of issued names, or coordinate between callers. It relies on the stamp and the random words to make collisions unlikely. It has no configuration, does not fetch anything, and does not ask questions. The principles below explain why each of these boundaries exists.


## Principles


### Every run stands alone

namo is stateless: it keeps no history, counter, cache, or lock, and one run never depends on another. This lets any caller, including several agents running in parallel on one machine, get a usable name without coordinating with anyone, and it means there is no state to go stale or become corrupted. The accepted cost is that namo cannot promise uniqueness across runs; a change that would need namo to remember what it issued conflicts with this principle, however useful the guarantee would be.


### The words come from hotdiva2000 as they are

namo's slugs come from hotdiva2000, and namo keeps its vocabulary as it is, irreverent words included. Those words make a name memorable enough to pick out of a listing, which is why namo uses words at all. The accepted cost is that some names are awkward where colleagues or the public can see them, such as branch names and shared directories. namo does not filter or replace the words, and it does not offer a tamer vocabulary as an option.


### No configuration

namo is zero-config: it reads no configuration file, defines no environment variables of its own, and takes every option as a flag. A given invocation therefore means the same thing on every machine and inside every agent environment, because no hidden setting can change it. The accepted cost is that there are no personal defaults; a caller who wants a particular prefix or stamp layout spells it out every time.


### Safe to run unattended

namo never prompts, runs no external program, and never accesses the network. Scripts and agents call it unattended, usually inside command substitution, so it must finish promptly, behave the same everywhere, and never wait on a person or a remote service. This rules out features such as remote word lists, online name registries, or interactive pickers, even when they would make names better.


### Safe output by default

Without an explicit trusted-input flag, a name contains only lowercase ASCII letters, digits, and hyphens. Callers put namo's output directly into file paths and command substitution without quoting or validating it, so the default must never produce a character that could break a path or a shell. Input that could carry unsafe characters, such as a verbatim prefix or a custom stamp layout, is allowed only behind a flag that names it as trusted. namo keeps these escape hatches even when no caller uses them, because they let a caller handle a rare case without waiting for a change to namo. The accepted cost is that a normalized prefix loses capitals and punctuation. A change that lets an unsafe character reach the output without such a flag conflicts with this principle.


### Better design over compatibility

namo makes clean breaks. When a better design for the name shape or the flags appears, namo changes to it without aliases, deprecation periods, or compatibility shims. Every caller is its author's own tooling, so keeping those callers working is done by updating them at the time of the change, not by carrying old behavior inside namo. The accepted cost is that a script or agent instruction that is missed fails or receives a differently shaped name.


## Applying this vision

A request clearly aligns when it makes namo better at the author's own naming, such as improving names for a use the author's scripts or agents already have or making the author's most common call shorter, and when it keeps every run independent, adds no configuration, stays offline and non-interactive, and keeps default output safe for paths and shells. A breaking change does not conflict merely because it breaks callers. A capability is not removed merely because the author's tooling rarely uses it.

A request clearly conflicts when it needs namo to remember earlier runs, read a configuration file or its own environment variables, prompt, run another program, or reach the network. It also conflicts when it lets unsafe characters into a name without an explicit trusted-input flag, filters, replaces, or offers an alternative to hotdiva2000's vocabulary, or gives namo an interface besides the command, such as a public Go package or a server mode. A request whose only beneficiaries are people outside the author's workflow conflicts as well.

Some requests need a human. These principles do not settle whether namo should do more than print names, for example by creating the file or directory it names, offer an output format other than plain text, read data such as prefixes from stdin, or accept positional arguments, and they do not settle whether namo should encode the convention for a particular use, such as a preset for branch or scratch names, rather than leaving that convention to its callers. Those are open choices to be judged on whether they serve the author's tooling. Whether a proposed default actually reflects the author's most common use is a fact only the author can confirm. An alignment judgment says only that a request fits this vision; it does not approve an implementation or replace technical review.
