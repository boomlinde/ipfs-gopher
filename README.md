ipfs-gopher is a gopher proxy that allows serving gopher over IPFS. For
plain files, it just forwards the file content to the client. For menus,
it uses a simple format with a header to distinguish menus from other
files, and to use the local proxy when resolving IPFS links.

Example session
---------------

    $ ipfs-gopher &
    $ curl gopher://localhost:7070/1/ipfs/QmRgTQvW7ab1bwjmeL4EQh1AVcWrUJn5ou3g1jv9aAqgRo/menu1

## Usage
    $ ipfs-gopher -h
    Usage of ipfs-gopher:
      -daemon string
        	The address of the IPFS daemon (default "localhost:5001")
      -host string
        	The host to use in IPFS selectors (default "localhost")
      -listen string
        	The address of the proxy (default "localhost:7070")
      -port string
        	The port to use in IPFS selectors (default "7070")

Menu format
-----------

Because the menu content presented to the gopher client depends on
ipfs-gopher settings and thus need to be modified in transit, we need to
separate regular files from menu files. To do this, the menu starts with
a single line, containing:

    <<<ipfs-gopher-menu>>>

ipfs-gopher menus are otherwise much like plain gopher menus:

    <type><name><TAB><selector><TAB><host><TAB><port><CRLF>

The most important difference is that the port and host may be omitted.
When they are, the selector should be understood as an IPFS link and
ipfs-gopher will fill in the rest of the columns. Relative IPFS links
may be used as selectors.

You may also omit the selector (useful for info type entries). Finally,
ipfs-gopher menu files may use Unix style line endings (LF) as well as
network style line endings (CRLF).

Linking to ipfs-gopher holes
----------------------------

Because the settings are client dependent, I suggest that any non-IPFS
gopher hole that wants to link to an ipfs-gopher hole should do so
through the hostname "ipfs-gopher" on port 7070. This way, users can
create a host alias for their preferred host. If not, a Sufficiently
Smart™ Gopher client could also replace the "ipfs-gopher" hostname with
the preferred hostname and the port with the preferred port.

TODO
----

-   Report IPFS errors to client
