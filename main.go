package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cccfs/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/fatih/color"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/exec"
	"path/filepath"
)

type Info struct {
	pvClaim string
	pvName string
	pvVolumeID string
	pvVolumePath string
	pvStatus string
	pvSCName string
	pvSize string
	pvSource string
}
func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		if ok, _ := util.Exists(filepath.Join(home, ".kube", "config")); ok {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			outfile, err := os.Create("/tmp/config")
			if err != nil {
				panic(err)
			}
			defer outfile.Close()
			cmd := exec.Command("bash", "-c", "kubectl config view --merge --flatten")
			cmd.Stdout = outfile
			err = cmd.Start(); if err != nil {
				panic(err)
			}
			cmd.Wait()
			kubeconfig = flag.String("kubeconfig", filepath.Join("/tmp/config"), "using $KUBECONFIG variables to the kubeconfig file")
		}
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	pvList, _ := clientset.CoreV1().PersistentVolumes().List(context.TODO(), v1.ListOptions{})
	info := Info{}
	fmt.Println("### Print pv info in console")
	for _, pv := range pvList.Items {
		info.pvClaim = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
		info.pvName = pv.Name
		info.pvVolumeID = pv.Annotations["gluster.kubernetes.io/heketi-volume-id"]
		info.pvSCName = pv.Spec.StorageClassName
		info.pvSize = pv.Spec.Capacity.Storage().String()
		info.pvStatus = string(pv.Status.Phase)
		if pv.Spec.Glusterfs != nil {
			info.pvVolumePath = pv.Spec.Glusterfs.Path
		}
		commaind := fmt.Sprintf(`
Claim: %s
PV: %s
VolumeID: %s
Size: %s
StorageClass: %s
Status: %s
VolumePath: %s
`, color.RedString(info.pvClaim), color.BlueString(info.pvName), color.HiWhiteString(info.pvVolumeID), color.CyanString(info.pvSize), color.CyanString(info.pvSCName), color.CyanString(info.pvStatus), color.CyanString(info.pvVolumePath))
		fmt.Print("*****")
		fmt.Println(commaind)
		//con := func (field string, val reflect.Value) color.Color {
		//	switch field {
		//	case "Claim":
		//		return color.BLUEBG
		//	case "PV":
		//		return color.BLUE
		//	case "VolumeID":
		//		return color.GREENBG
		//	case "VolumePath":
		//		return color.BLUE
		//	case "Size":
		//		return color.BLUE
		//	case "StorageClass":
		//		return color.BLUE
		//	case "Status":
		//		return color.BLUE
		//	default:
		//		return ""
		//	}
		//}
		//headers := []string{"PV", "VolumeID", "Size", "StorageClass", "Status", "VolumePath", "Claim"}
		//tb, err := gotable.CreateTable(headers, gotable.WithColorController(con))
		//if err != nil {
		//	fmt.Println("Create table failed: ", err.Error())
		//	return
		//}
		//value := gotable.CreateEmptyValueMap()
		//value["Claim"] = gotable.CreateValue(info.pvClaim)
		//value["PV"] = gotable.CreateValue(info.pvName)
		//value["VolumeID"] = gotable.CreateValue(info.pvVolumeID)
		//value["VolumePath"] = gotable.CreateValue(info.pvVolumePath)
		//value["Size"] = gotable.CreateValue(info.pvSize)
		//value["StorageClass"] = gotable.CreateValue(info.pvSCName)
		//value["Status"] = gotable.CreateValue(info.pvStatus)
		//tb.AddValue(value)
		//tb.PrintTable()
	}
}