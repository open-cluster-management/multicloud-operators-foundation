// Licensed Materials - Property of IBM
// 5737-E67
// (C) Copyright IBM Corporation 2016, 2019 All Rights Reserved
// US Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
// IBM Confidential
// OCO Source Materials
// 5737-E67
// (C) Copyright IBM Corporation 2016, 2019 All Rights Reserved
// The source code for this program is not published or otherwise divested of its trade secrets, irrespective of what has been
// deposited with the U.S. Copyright Office.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/open-cluster-management/multicloud-operators-foundation/cmd/klusterlet/app"
	"github.com/open-cluster-management/multicloud-operators-foundation/cmd/klusterlet/app/options"
	"github.com/open-cluster-management/multicloud-operators-foundation/pkg/signals"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"github.com/spf13/pflag"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	s := options.NewKlusterletRunOptions()
	s.AddFlags(pflag.CommandLine)

	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := signals.SetupSignalHandler()
	if err := app.Run(s, stopCh); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
