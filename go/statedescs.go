// Obdi - a REST interface and GUI for deploying software
// Copyright (C) 2014  Mark Clarkson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
)

type PostedData struct {
	Grain string
	Text  string
}

// ***************************************************************************
// GO RPC PLUGIN
// ***************************************************************************

func (t *Plugin) GetRequest(args *Args, response *[]byte) error {

	// Check for required query string entries

	if len(args.QueryString["env_id"]) == 0 {
		ReturnError("'env_id' must be set", response)
		return nil
	}

	//env_id, _ := strconv.ParseInt(args.QueryString["env_id"][0], 10, 64)

	if len(args.QueryString["version"]) == 0 {
		ReturnError("'version' must be set", response)
		return nil
	}

	version := args.QueryString["version"][0]

	// Get the DcSysName and [Env]SysName for this env_id for using
	// as arguments to the script that will be run later.
	// GET queries always return an array of items, even for 1 item.
	envs := []Env{}
	resp, err := GET("https://127.0.0.1/api/"+
		args.PathParams["login"]+"/"+args.PathParams["GUID"], "envs"+
		"?env_id="+args.QueryString["env_id"][0])
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		txt := fmt.Sprintf("Error reading Body ('%s').", err.Error())
		ReturnError(txt, response)
		return nil
	} else {
		json.Unmarshal(b, &envs)
	}
	// If envs is empty then we don't have permission to see it
	// or the env does not exist so bug out.
	if len(envs) == 0 {
		txt := "The requested environment id does not exist" +
			" or the permissions to access it are insufficient."
		ReturnError(txt, response)
		return nil
	}

	// If version is 0 then use the main branch, for example, 'test'
	// rather than 'test_0.1.0'

	version_text := ""
	if version != "0" {
		version_text = "_" + version
	}

	sa := ScriptArgs{
		ScriptName: "statelist.sh",
		CmdArgs:    envs[0].SysName + version_text,
		EnvVars:    "",
		EnvCapDesc: "SALT_WORKER",
		Type:       2,
	}

	// Send the Job ID as the RPC reply (back to the master)

	var jobid int64
	if jobid, err = t.RunScript(args, sa, response); err != nil {
		// RunScript wrote the error
		return nil
	}

	reply := Reply{jobid, "", SUCCESS, ""}
	jsondata, err := json.Marshal(reply)
	if err != nil {
		ReturnError("Marshal error: "+err.Error(), response)
		return nil
	}
	*response = jsondata

	return nil
}

func (t *Plugin) PostRequest(args *Args, response *[]byte) error {

	ReturnError("Internal error: Unimplemented HTTP POST", response)
	return nil
}

func (t *Plugin) HandleRequest(args *Args, response *[]byte) error {

	// All plugins must have this.

	if len(args.QueryType) > 0 {
		switch args.QueryType {
		case "GET":
			t.GetRequest(args, response)
			return nil
		case "POST":
			t.PostRequest(args, response)
			return nil
		}
		ReturnError("Internal error: Invalid HTTP request type for this plugin "+
			args.QueryType, response)
		return nil
	} else {
		ReturnError("Internal error: HTTP request type was not set", response)
		return nil
	}
}

func main() {

	//logit("Plugin starting")

	plugin := new(Plugin)
	rpc.Register(plugin)

	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		txt := fmt.Sprintf("Listen error. %s", err)
		logit(txt)
	}

	//logit("Plugin listening on port " + os.Args[1])

	//for {
	if conn, err := listener.Accept(); err != nil {
		txt := fmt.Sprintf("Accept error. %s", err)
		logit(txt)
	} else {
		//logit("New connection established")
		rpc.ServeConn(conn)
	}
	//}
}
