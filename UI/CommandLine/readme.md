# Command Line Interface (CLI)

## Tutorial: Add command to CLI
1. Enter the following in terminal in the /UI/CommandLine directory:
    > cobra add \<command\>

    A new command is added in the /cmd folder in a file named \<command\>.go.

2. Insert any functionality and/or logic into the newly created \<command\>.go file, in the <i>run()</i> function

3. In order to take arguments, conversion from []string is needed, and should be applied on the "[]string args" argument of the <i>run</i> command

