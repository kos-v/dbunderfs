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

package log

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type StdoutFormatter struct {
}

func (f *StdoutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatString := entry.Message

	switch entry.Level.String() {
	case "warning":
		formatString = color.YellowString(formatString)
	case "error", "fatal", "panic":
		formatString = color.RedString(formatString)
	}

	formatString += "\n"

	return []byte(formatString), nil
}
