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

import "flag"
import "github.com/gocql/gocql"
import "fmt"
import "text/tabwriter"
import "bufio"
import "os"

var netaddr = flag.String("addr","127.0.0.1","network address of cassandra")

var cql_session *gocql.Session
var tw *tabwriter.Writer
var inpt *bufio.Reader

func main() {
	tw = tabwriter.NewWriter(os.Stdout,0,0,2,' ',tabwriter.Debug)
	inpt = bufio.NewReader(os.Stdin)
	flag.Parse()
	cluster := gocql.NewCluster(*netaddr)
	session, err := cluster.CreateSession()
	if err!=nil {
		fmt.Println("no connection",err)
		return
	}
	cql_session = session
	repl()
}

