package cli

// failureContract is the stderr contract every failure follows, composed
// into the long description of every command so they cannot drift apart.
const failureContract = `A failure prints exactly one line on stderr, 'namo: <message>', and no
usage text; a line break or carriage return inside the message is shown
as \n or \r.`

// textCommandContract is the contract every command that only prints text
// shares: the help command, the completion command, and the help topics.
const textCommandContract = `This command never prompts, works without a terminal, and reads and
writes no files.

` + failureContract
