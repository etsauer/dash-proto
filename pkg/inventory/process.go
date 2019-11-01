package inventory

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	cp "github.com/redhat-cop/dash/pkg/copy"
)

// Process Inventory of ResourceGroups
func (i *Inventory) Process(ns *string) error {

	// create a temp directory for resources
	output_dir := filepath.Clean(i.Output)
	if _, err := os.Stat(output_dir); os.IsNotExist(err) {
		err = os.Mkdir(output_dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	defer func() {
		os.Remove("output_dir")
	}()

	if i.Namespace != "" {
		ns = &i.Namespace
	}

	if i.ResourceGroups != nil {
		for _, rg := range i.ResourceGroups {
			err := rg.Process(ns)
			if err != nil {
				return err
			}

			err = rg.Reconcile(ns, i.Args)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// Process ResourceGroup of Resource
func (rg *ResourceGroup) Process(ns *string) error {

	if rg.Namespace != "" {
		ns = &rg.Namespace
	}

	if rg.Resources != nil {
		for _, r := range rg.Resources {
			err := r.Process(ns)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Process a single resource
func (r *Resource) Process(ns *string) error {

	// create a temp directory for resources
	output_dir := filepath.Clean(r.Output + "/" + r.Action)
	if _, err := os.Stat(output_dir); os.IsNotExist(err) {
		err = os.Mkdir(output_dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if r.Namespace != "" {
		ns = &r.Namespace
	}

	log.Println("Resource: " + r.Name + ", Namespace: " + *ns)

	// TODO: Determine the type of resource
	if !reflect.DeepEqual(FileTemplate{}, r.File) {
		err := r.File.Process(ns, r)
		if err != nil {
			return err
		}
	} else if !reflect.DeepEqual(HelmChart{}, r.Helm) {
		err := r.Helm.Process(ns, r)
		if err != nil {
			return err
		}
	} else if !reflect.DeepEqual(OpenShiftTemplate{}, r.OpenShiftTemplate) {
		err := r.OpenShiftTemplate.Process(ns, r)
		if err != nil {
			return err
		}
	}

	return nil
}

// Reconcile ensures all generated resources are applied via the appropriate verb
func (rg *ResourceGroup) Reconcile(ns *string, args []string) error {

	applyP := rg.Output + "/apply"
	applyAbs, err := filepath.Abs(applyP)
	if err != nil {
		return err
	}
	cmdArgs := append(
		[]string{"apply", "-f", filepath.Clean(applyAbs), "--recursive"},
		args...,
	)
	if *ns != "" {
		cmdArgs = append(cmdArgs, "-n", *ns)
	}

	cmd := exec.Command("kubectl", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}
	fmt.Printf("%s\n", stdoutStderr)

	return nil

}

func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	log.Printf("Source: %s, Dest: %s\n", src, dst)
	if sourceFileStat.IsDir() {
		err = cp.Dir(src, dst)
		if err != nil {
			return err
		}
		return nil
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	err = cp.File(src, dst)
	if err != nil {
		return err
	}
	return nil
}
