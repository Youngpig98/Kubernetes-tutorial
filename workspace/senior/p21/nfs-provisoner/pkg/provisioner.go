package pkg

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v8/controller"
)

type NfsProvisioner struct {
	Server        string
	RootMountPath string
	MountPath     string
}

var _ controller.Provisioner = &NfsProvisioner{}

func (nfs *NfsProvisioner) Provision(_ context.Context, opt controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
	// 创建卷
	mountPath := filepath.Join(nfs.MountPath, opt.PVName)
	if err := os.Mkdir(mountPath, 0755); err != nil {
		log.Printf("mkdir %s error: %s", mountPath, err.Error())
		return nil, controller.ProvisioningFinished, errors.Wrap(err, "mkdir error")
	}

	// 返回pv对象
	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: opt.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *opt.StorageClass.ReclaimPolicy,
			AccessModes:                   opt.PVC.Spec.AccessModes,
			MountOptions:                  opt.StorageClass.MountOptions,
			Capacity: v1.ResourceList{
				v1.ResourceStorage: opt.PVC.Spec.Resources.Requests[v1.ResourceStorage],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server:   nfs.Server,
					Path:     filepath.Join(nfs.RootMountPath, opt.PVName),
					ReadOnly: false,
				},
			},
		},
	}, controller.ProvisioningFinished, nil
}

func (nfs *NfsProvisioner) Delete(_ context.Context, pv *v1.PersistentVolume) error {
	return os.Remove(filepath.Join(nfs.MountPath, pv.Name))
}
