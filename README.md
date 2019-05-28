# gn2bazel

gn2bazel is an experimental tool to convert GN build files to Bazel.

It is not fully functional and doesn't implement all possible GN build rules,
nor does it necessarily produce a correct Bazel BUILD file.

It works by invoking `gn desc --format=json out.gn/<release> *` and converting
the JSON data to a BUILD file. To allow access to code in subdirectories, it
adds dummy BUILD files to each directory that publicly export all files. This
is far from ideal and there may be better approaches.

As this is experimental, any contributions or improvements are very welcome,
but they must be released publicly as per the GNU General Public License. Pull
requests will be gladly accepted.

There are several closed source tools that work similarly, but no known open
source tools, see:
https://groups.google.com/a/chromium.org/d/msg/chromium-dev/wl1T6XX2gg8/S77auAlOAAAJ,
https://crbug.com/webrtc/6412 and https://github.com/nodejs/TSC/issues/464.

## Installation:

`go get tmthrgd.dev/go/gn2bazel`

## Usage:

`gn2bazel [-dir <gn directory>] [-out <bazel directory>] [-exclude <exclude targets regexp>] <release>`

## License:

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
