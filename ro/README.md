ro/ro
=====

Example Reittiopas client.

Installation:
-------------

    go get github.com/errnoh/ro/ro
    ro --from "Homestreet 3" --to "Workpath 13"

Parameters:
-----------

    --from       Starting location (i.e. "Kauppakuja 1, helsinki")
    --to         Destination
    --date       Date (YYYYMMDD)
    --time       Time (HHMM)
    --limit      How many routes are displayed (default 3, max 5)
    --optimize   fastest|least_transfers|least_walking

