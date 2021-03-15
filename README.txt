===============
   ethersim
===============

Simulates an old-style, half-duplex Ethernet link, including
features like exponential backoff and checksumming/voluntary
packet corruption.

Building:

On athena, clone this repository and run `setup ggo` to configure
a go toolchain. On other systems, follow the official guides from
the Golang team at Google to install Go for your operating system.

Once done, ensure you have make installed on your system (should be
there by default on athena), and run make to build ethersim.

Usage:

Make sure you have /tmp writeable, and execute ./ethersim to view
the command line options:

$ ./ethersim
Usage: ethersim [OPTIONS]
        -s: Start ethersim server / emulated media
        -c [ID]: Connect to session [ID]

First, start ethersim's server via the above. You will get an ID
back. On the same machine, in multiple terminals or as multiple users
(especially in the athena case), utilize this ID to connect to the
session. Each connected client (with -c) is an Ethernet node, attached
to the central node and ethernet media (-s).

