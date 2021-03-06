// Copyright 2018 David Sansome
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"path"

	"github.com/golang/protobuf/proto"

	pb "github.com/davidsansome/tsurukame/proto"
)

type Mapper func(pb.Subject) pb.Subject

var (
	directory = flag.String("directory", "data", "Directory to read data files from")
	mapper    = flag.String("mapper", "", "One-off mapper to run")

	mappers = map[string]Mapper{
		// Put mappers here.
		"RemoveNone": RemoveNone,
	}
)

func main() {
	flag.Parse()

	mapper := mappers[*mapper]
	if mapper == nil {
		panic("Mapper not found")
	}
	if err := ListAll(mapper); err != nil {
		panic(err)
	}
}

func ListAll(mapper Mapper) error {
	files, err := ioutil.ReadDir(*directory)
	if err != nil {
		return err
	}

	for _, f := range files {
		filename := path.Join(*directory, f.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var oldSubject pb.Subject
		if err := proto.Unmarshal(data, &oldSubject); err != nil {
			return err
		}

		newSubject := mapper(oldSubject)

		data, err = proto.Marshal(&newSubject)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filename, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
