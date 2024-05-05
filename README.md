# Dictionary Web Server

A simple web gateway to a dict server. This server only implements DEFINE and
always searches all of the databases on the dict server.

## Usage

``` bash
go run main.go --listen-addr :8000 --dict-server dict.org:2628
```

See https://servers.freedict.org/ for a list of dict servers.

## Copyright

Dictionary Web Server
Copyright (C) 2024 Raymond Miller <ray@1729.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
