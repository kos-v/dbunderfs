/*
   Copyright 2021 The DbunderFS Contributors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package mysql

import "github.com/kos-v/dsnparser"

type DSN struct {
	ParsedDSN dsnparser.DSN
}

func (d *DSN) GetDatabase() string {
	return d.ParsedDSN.GetPath()
}

func (d *DSN) GetTablePrefix() string {
	if !d.ParsedDSN.HasParam("tblprefix") {
		return ""
	}
	return d.ParsedDSN.GetParam("tblprefix")
}

func (d *DSN) ToString() string {
	protocol := "tcp"
	if d.ParsedDSN.HasParam("protocol") {
		protocol = d.ParsedDSN.GetParam("protocol")
	}

	return d.ParsedDSN.GetUser() + ":" + d.ParsedDSN.GetPassword() + "@" +
		protocol + "(" + d.ParsedDSN.GetHost() + ")/" +
		d.ParsedDSN.GetPath()
}
