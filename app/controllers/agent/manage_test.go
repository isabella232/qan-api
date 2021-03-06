/*
   Copyright (c) 2016, Percona LLC and/or its affiliates. All rights reserved.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package agent

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddVisualExplain(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		explain, err := addVisualExplain([]byte{})
		assert.Nil(t, explain)
		assert.Equal(t, err, errors.New("cannot unmarshal classic expain to do visual explain: unexpected end of JSON input"))
	})

	t.Run("ok", func(t *testing.T) {
		const agentResp = `{"Classic":[{"Id":1,"SelectType":"SIMPLE","Table":"sbtest1","Partitions":null,"CreateTable":null,"Type":"const","PossibleKeys":"PRIMARY","Key":"PRIMARY","KeyLen":"4","Ref":"const","Rows":1,"Filtered":null,"Extra":""}],"JSON":"{\n  \"query_block\": {\n    \"select_id\": 1,\n    \"table\": {\n      \"table_name\": \"sbtest1\",\n      \"access_type\": \"const\",\n      \"possible_keys\": [\"PRIMARY\"],\n      \"key\": \"PRIMARY\",\n      \"key_length\": \"4\",\n      \"used_key_parts\": [\"id\"],\n      \"ref\": [\"const\"],\n      \"rows\": 1,\n      \"filtered\": 100\n    }\n  }\n}"}`
		got, err := addVisualExplain([]byte(agentResp))
		assert.NoError(t, err)
		expected, err := loadExpected("testdata/visual-explain.txt", got)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expected), string(got))
	})
}

func loadExpected(file string, got []byte) ([]byte, error) {
	updateTestData := os.Getenv("UPDATE_TEST_DATA")
	if updateTestData != "" {
		ioutil.WriteFile(file, got, 0666)
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}
