/*
Copyright (c) 2018 Simon Schmidt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import "fmt"

import "regexp"
import "reflect"
import "strings"

var list_keyspaces = regexp.MustCompile(`(?i:list\s+keyspaces)`)
var list_tables = regexp.MustCompile(`(?i:list\s+tables\s+([a-z0-9_]+))`)
var list_columns = regexp.MustCompile(`(?i:list\s+columns\s+([a-z0-9_]+)\.([a-z0-9_]+)(\s+full)?)`)
var list_do_short = regexp.MustCompile(`(?i:do\s+)(.*)`)
var list_do_long = regexp.MustCompile(`(?i:do\:)`)

func unwrap(i interface{}) func()string {
	switch v := i.(type) {
	case *string: return func()string { return fmt.Sprintf("%q",*v) }
	case *[]byte: return func()string { return fmt.Sprintf("%q",*v) }
	}
	return func()string { return strings.Replace(fmt.Sprint(reflect.ValueOf(i).Elem().Interface()),"\t"," ",0) }
}

func repl() {
	for ;; tw.Flush() {
		fmt.Print("-->");
		line,err := inpt.ReadString('\n')
		if err!=nil { break }
		
		var query string
		if sub := list_keyspaces.FindStringSubmatch(line); len(sub)!=0 {
			query = `SELECT * FROM system_schema.keyspaces;`
		} else if sub := list_tables.FindStringSubmatch(line); len(sub)!=0 {
			query = `SELECT keyspace_name,table_name FROM system_schema.columns WHERE keyspace_name = '`+sub[1]+`' GROUP BY table_name;`
		} else if sub := list_columns.FindStringSubmatch(line); len(sub)!=0 {
			//fmt.Fprintf(tw,"%q\n",sub)
			if len(sub[3])!=0 {
				query = `SELECT * FROM system_schema.columns WHERE keyspace_name = '`+sub[1]+`' AND table_name = '`+sub[2]+`';`
			} else {
				query = `SELECT keyspace_name as "keyspace",table_name as "table",column_name as "column","type"
				FROM system_schema.columns WHERE keyspace_name = '`+sub[1]+`' AND table_name = '`+sub[2]+`';`
			}
		} else if sub := list_do_short.FindStringSubmatch(line); len(sub)!=0 {
			query = sub[1]
		} else if sub := list_do_long.FindStringSubmatch(line); len(sub)!=0 {
			query,err = inpt.ReadString(';')
			if err!=nil {
				fmt.Fprintf(tw,"stdin error: %v\n",err)
				continue
			}
			fmt.Fprintf(tw,"query is %q\n",query)
		} else {
			fmt.Fprintf(tw,"unknown command %q\n",line)
		}
		iter := cql_session.Query(query).Iter()
		cols := iter.Columns()
		if len(cols)==0 {
			fmt.Fprintf(tw,"no result: error=%v\n",iter.Close())
			continue
		}
		fetch := make([]interface{},len(cols))
		retrs := make([]func()string,len(cols))
		for i,col := range cols {
			fetch[i] = col.TypeInfo.New()
			retrs[i] = unwrap(fetch[i])
			fmt.Fprintf(tw,"%s\t",col.Name)
		}
		fmt.Fprintln(tw)
		for i := range cols {
			fmt.Fprintf(tw,"%d ---\t",i)
		}
		fmt.Fprintln(tw)
		for iter.Scan(fetch...) {
			for _,r := range retrs {
				fmt.Fprintf(tw,"%s\t",r())
			}
			fmt.Fprintln(tw)
		}
	}
}

