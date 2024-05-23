package main

import (
	"fmt"
	"os"

	"github.com/openshift-kni/lifecycle-agent/api/ibiconfig"
	"github.com/openshift-kni/lifecycle-agent/ib-cli/installationiso"
	"github.com/openshift-kni/lifecycle-agent/lca-cli/ops"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:   os.Stdout,
	Level: logrus.InfoLevel,
}

func main() {

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	testIBI := &ibiconfig.IBIPrepareConfig{
		RHCOSLiveISO:       "https://mirror.openshift.com/pub/openshift-v4/amd64/dependencies/rhcos/latest/rhcos-live.x86_64.iso",
		InstallationDisk:   "/dev/disk/by-id/scsi-0QEMU_QEMU_HARDDISK_drive-scsi0-0-0-0",
		SeedImage:          "quay.io/ocp-edge-qe/ib-seedimage-public:ci-4.15-trwest",
		SeedVersion:        "4.15.13",
		AuthFile:           "/home/trwest/pull-secret.txt",
		PullSecretFile:     "/home/trwest/pull-secret.txt",
		SSHPublicKeyFile:   "/home/trwest/.ssh/id_rsa.pub",
		Shutdown:           true,
		PrecacheBestEffort: true,
	}

	err := testIBI.Validate()
	if err != nil {
		log.Fatal(err)
	}
	testIBI.SetDefaultValues()

	log.Info("testing")

	hostCommandsExecutor := ops.NewRegularExecutor(log, true)
	op := ops.NewOps(log, hostCommandsExecutor)

	isoCreator := installationiso.NewInstallationIso(log, op, "/home/trwest/ibi_images")
	fmt.Println(isoCreator)
	if err = isoCreator.Create(testIBI); err != nil {
		err = fmt.Errorf("failed to create installation ISO: %w", err)
		log.Errorf(err.Error())
	}

	log.Info("Installation ISO created successfully!")
}
